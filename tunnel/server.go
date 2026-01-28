package tunnel

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"vpn_proto/config"
	"vpn_proto/crypto"
	"vpn_proto/utils"

	"strings"
	"time"

	"github.com/songgao/water"
)

// StartServer starts the VPN server
func StartServer(ctx context.Context, cfg *config.Config) error {
	// 1. Create TUN interface (Only needed if we do L3 routing, but we keep it for backward compat)
	// For pure proxy mode, we don't strict need it, but the code currently mixes them.
	// Let's keep it but handle errors gracefully if not sudo.

	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	var ifceName string
	if err != nil {
		utils.Warn("Failed to create TUN (running without sudo?): %v. L3 mode disabled.", err)
	} else {
		defer ifce.Close()
		ifceName = ifce.Name()
		utils.Success("Interface %s created", ifceName)
	}

	// 2. Setup TLS listener with mTLS
	tlsConfig, err := crypto.LoadServerTLS(cfg.CertFile, cfg.KeyFile, cfg.CAFile)
	if err != nil {
		return fmt.Errorf("failed to load mTLS config: %v (Did you run -gen-certs?)", err)
	}

	lc := net.ListenConfig{}
	ln, err := lc.Listen(ctx, "tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		return err
	}
	defer ln.Close()

	utils.Success("Server listening on 0.0.0.0:%d", cfg.Port)

	// Accept loop
	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil // Normal shutdown
			default:
				utils.Error("Accept error: %v", err)
				continue
			}
		}

		// Handle connection in a goroutine
		go handleClient(conn, ifce, tlsConfig)
	}
}

func handleClient(rawConn net.Conn, ifce *water.Interface, tlsConfig *tls.Config) {
	// Upgrade to TLS
	conn := tls.Server(rawConn, tlsConfig)
	// Don't defer close immediately, we might need it open long time
	defer conn.Close()

	if err := conn.Handshake(); err != nil {
		utils.Warn("TLS handshake failed: %v", err)
		return
	}

	conn.SetDeadline(time.Now().Add(10 * time.Second))
	authData, err := utils.ReadPacket(conn)
	conn.SetDeadline(time.Time{}) // Clear deadline
	if err != nil {
		utils.Warn("Failed to read auth: %v", err)
		return
	}

	authStr := string(authData)
	// Check format: "SECRET" (Legacy/TUN) or "SECRET|CONNECT|target" (Proxy)
	// With mTLS, "SECRET" might be just a placeholder or empty.

	parts := strings.SplitN(authStr, "|", 3)
	// secret := parts[0]
	// We trust the connection because of mTLS.
	// However, we still parse the pipe format for Mode selection.

	if len(parts) == 3 && parts[1] == "CONNECT" {
		handleProxyRequest(conn, parts[2])
	} else {
		if ifce == nil {
			utils.Warn("Client requested TUN mode but TAP/TUN is not available")
			return
		}
		handleTunSession(conn, ifce)
	}
}

func handleProxyRequest(clientConn net.Conn, target string) {
	if IsBlocked(target) {
		utils.Block("Connection denied: %s", target)
		utils.WritePacket(clientConn, []byte("FAIL")) // Or silent drop
		return
	}

	// Connect to target
	// Use DoH to resolve hostname to IP
	host, port, _ := net.SplitHostPort(target)

	// Check if host is an IP, if not, resolve via DoH
	if net.ParseIP(host) == nil && host != "localhost" {
		utils.Secure("Resolving %s via DoH...", host)
		ip, err := ResolveDoH(host)
		if err != nil {
			utils.Warn("DoH failed for %s: %v. Fallback to system DNS.", host, err)
			// Fallback or fail?
			// Let's fallback for robustness but log warning
		} else {
			target = net.JoinHostPort(ip, port)
		}
	}

	targetConn, err := net.DialTimeout("tcp", target, 10*time.Second)
	if err != nil {
		utils.Warn("Failed to dial target %s: %v", target, err)
		utils.WritePacket(clientConn, []byte("FAIL"))
		return
	}
	defer targetConn.Close()

	// utils.Info("Proxy: %s", target)
	// reduced log spam
	utils.WritePacket(clientConn, []byte("OK"))

	// Pipe
	go io.Copy(clientConn, targetConn)
	io.Copy(targetConn, clientConn)
}

func handleTunSession(conn net.Conn, ifce *water.Interface) {
	utils.Info("Starting TUN session for %s", conn.RemoteAddr())

	// Start Pump
	// Tun -> TCP
	go func() {
		buf := make([]byte, 2048) // MTU is usually 1500
		for {
			n, err := ifce.Read(buf)
			if err != nil {
				utils.Error("TUN read error: %v", err)
				return
			}
			err = utils.WritePacket(conn, buf[:n])
			if err != nil {
				// log.Printf("TCP write error: %v", err)
				conn.Close()
				return
			}
		}
	}()

	// TCP -> Tun
	for {
		packet, err := utils.ReadPacket(conn)
		if err != nil {
			if err != io.EOF {
				utils.Error("TCP read error: %v", err)
			}
			break
		}
		_, err = ifce.Write(packet)
		if err != nil {
			utils.Error("TUN write error: %v", err)
			break
		}
	}
}
