package tunnel

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"vpn_proto/config"
	"vpn_proto/crypto"
	"vpn_proto/utils"

	"github.com/songgao/water"
)

func StartClient(cfg *config.Config) error {
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

	log.Printf("Connected using %s", crypto.GetCipherSuiteName(conn.ConnectionState().CipherSuite))

	// 2. Authenticate
	err = utils.WritePacket(conn, []byte(cfg.Secret))
	if err != nil {
		return fmt.Errorf("failed to send auth: %v", err)
	}
	log.Printf("Authentication sent")

	// 3. Create TUN connection
	ifce, err := water.New(water.Config{
		DeviceType: water.TUN,
	})
	if err != nil {
		return fmt.Errorf("failed to create TUN interface: %v", err)
	}
	defer ifce.Close()

	log.Printf("Interface %s created", ifce.Name())
	log.Println("Note: You must configure IP and routes manually for this prototype to work fully as a VPN.")
	log.Printf("Example: sudo ifconfig %s 10.0.0.2 10.0.0.1 up", ifce.Name())
	log.Printf("Example: sudo route add default 10.0.0.1")

	// 4. Start Pump
	// TUN -> TCP
	go func() {
		buf := make([]byte, 2048)
		for {
			n, err := ifce.Read(buf)
			if err != nil {
				log.Printf("TUN read error: %v", err)
				return
			}
			err = utils.WritePacket(conn, buf[:n])
			if err != nil {
				log.Printf("TCP write error: %v", err)
				conn.Close()
				return
			}
		}
	}()

	// TCP -> TUN
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
	return nil
}
