package tunnel

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"vpn_proto/config"
	"vpn_proto/crypto"
	"vpn_proto/utils"

	"github.com/songgao/water"
	"strings"
	"time"
)

// StartServer starts the VPN server
func StartServer(cfg *config.Config) error {
	// 1. Create TUN interface (Only needed if we do L3 routing, but we keep it for backward compat)
    // For pure proxy mode, we don't strict need it, but the code currently mixes them.
    // Let's keep it but handle errors gracefully if not sudo.
    
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
    var ifceName string
	if err != nil {
		log.Printf("Warning: failed to create TUN (running without sudo?): %v. L3 mode will fail, Proxy mode will work.", err)
	} else {
        defer ifce.Close()
        ifceName = ifce.Name()
	    log.Printf("Interface %s created", ifceName)
    }

	// 2. Setup TLS listener
	tlsConfig, err := crypto.GenerateSelfSignedCert()
	if err != nil {
		return fmt.Errorf("failed to generate cert: %v", err)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		return err
	}

	log.Printf("Server listening on 0.0.0.0:%d", cfg.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		// Handle connection in a goroutine
		go handleClient(conn, ifce, cfg, tlsConfig)
	}
}

func handleClient(rawConn net.Conn, ifce *water.Interface, cfg *config.Config, tlsConfig *tls.Config) {
	// Upgrade to TLS
	conn := tls.Server(rawConn, tlsConfig)
	// Don't defer close immediately, we might need it open long time
    defer conn.Close()

	if err := conn.Handshake(); err != nil {
		log.Printf("TLS handshake failed: %v", err)
		return
	}

    conn.SetDeadline(time.Now().Add(10 * time.Second))
    authData, err := utils.ReadPacket(conn)
    conn.SetDeadline(time.Time{}) // Clear deadline
    if err != nil {
        log.Printf("Failed to read auth: %v", err)
        return
    }
    
    authStr := string(authData)
    // Check format: "SECRET" (Legacy/TUN) or "SECRET|CONNECT|target" (Proxy)
    
    parts := strings.SplitN(authStr, "|", 3)
    secret := parts[0]
    
    if secret != cfg.Secret {
        log.Printf("Invalid secret from %s", conn.RemoteAddr())
        return
    }
    
    if len(parts) == 3 && parts[1] == "CONNECT" {
        handleProxyRequest(conn, parts[2])
    } else {
        if ifce == nil {
            log.Printf("Client requested TUN mode but TAP/TUN is not available")
            return
        }
        handleTunSession(conn, ifce)
    }
}

func handleProxyRequest(clientConn net.Conn, target string) {
    if IsBlocked(target) {
        log.Printf("[BLOCKED] Connection to tracker denied: %s", target)
        utils.WritePacket(clientConn, []byte("FAIL")) // Or silent drop
        return
    }

    log.Printf("Proxying connection to %s", target)
    
    // Connect to target
    targetConn, err := net.DialTimeout("tcp", target, 10*time.Second)
    if err != nil {
        log.Printf("Failed to dial target %s: %v", target, err)
        utils.WritePacket(clientConn, []byte("FAIL"))
        return
    }
    defer targetConn.Close()
    
    utils.WritePacket(clientConn, []byte("OK"))
    
    // Pipe
    go io.Copy(clientConn, targetConn)
    io.Copy(targetConn, clientConn)
}

func handleTunSession(conn net.Conn, ifce *water.Interface) {
	log.Printf("Starting TUN session for %s", conn.RemoteAddr())

	// Start Pump
	// Tun -> TCP
	go func() {
		buf := make([]byte, 2048) // MTU is usually 1500
		for {
			n, err := ifce.Read(buf)
			if err != nil {
				log.Printf("TUN read error: %v", err)
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
				log.Printf("TCP read error: %v", err)
			}
			break
		}
		_, err = ifce.Write(packet)
		if err != nil {
			log.Printf("TUN write error: %v", err)
			break
		}
	}
}
