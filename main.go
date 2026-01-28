package main

import (
	"flag"
	"log"
	"os"

	"vpn_proto/config"
	"vpn_proto/tunnel"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Error loading config: %v", err)
		flag.Usage()
		os.Exit(1)
	}

	log.Printf("Starting VPN in %s mode", cfg.Mode)

	if cfg.Mode == config.ModeServer {
		err = tunnel.StartServer(cfg)
	} else if cfg.Mode == config.ModeSocks {
		err = tunnel.StartSocksClient(cfg)
	} else {
		err = tunnel.StartClient(cfg)
	}

	if err != nil {
		log.Fatalf("Critical error: %v", err)
	}
}
