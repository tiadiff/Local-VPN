package tunnel

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"vpn_proto/config"
	"vpn_proto/crypto"
	"vpn_proto/utils"

	"github.com/songgao/water"
)

func StartClient(ctx context.Context, cfg *config.Config) error {
	serverEndpoint := fmt.Sprintf("%s:%d", cfg.ServerAddr, cfg.Port)
	log.Printf("Connecting to server at %s...", serverEndpoint)

	// 1. Connect to Server
	tlsConfig, err := crypto.LoadClientTLS(cfg.CertFile, cfg.KeyFile, cfg.CAFile)
	if err != nil {
		return fmt.Errorf("failed to load mTLS: %v", err)
	}
	conn, err := tls.Dial("tcp", serverEndpoint, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	// 2. Create TUN connection
	ifce, err := water.New(water.Config{DeviceType: water.TUN})
	if err != nil {
		// Fallback or error? For client TUN mode, this is critical.
		return fmt.Errorf("failed to create TUN interface (sudo required): %v", err)
	}
	defer ifce.Close()

	utils.Success("Connected to server! Interface: %s", ifce.Name())

	// 3. Auth
	authMsg := []byte(cfg.Secret)
	if err := utils.WritePacket(conn, authMsg); err != nil {
		return err
	}

	// 4. Pump
	// TUN -> TCP
	go func() {
		buf := make([]byte, 2048)
		for {
			n, err := ifce.Read(buf)
			if err != nil {
				utils.Error("TUN read error: %v", err)
				return
			}
			err = utils.WritePacket(conn, buf[:n])
			if err != nil {
				// log.Printf("TCP write error: %v", err)
				return
			}
		}
	}()

	// TCP -> TUN
	done := make(chan struct{})
	go func() {
		select {
		case <-ctx.Done():
			conn.Close()
		case <-done:
		}
	}()
	defer close(done)

	for {
		packet, err := utils.ReadPacket(conn)
		if err != nil {
			if err != io.EOF {
				select {
				case <-ctx.Done():
					return nil
				default:
					utils.Error("TCP read error: %v", err)
				}
			}
			break
		}
		_, err = ifce.Write(packet)
		if err != nil {
			utils.Error("TUN write error: %v", err)
			break
		}
	}
	return nil
}
