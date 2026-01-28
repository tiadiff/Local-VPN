package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"vpn_proto/config"
	"vpn_proto/crypto"
	"vpn_proto/tunnel"
	"vpn_proto/utils"
)

func main() {
	printBanner()

	cfg, err := config.Load()
	if err != nil {
		utils.Error("Error loading config: %v", err)
		flag.Usage()
		os.Exit(1)
	}

	// Check if running as root/admin if required
	if os.Geteuid() != 0 && cfg.Mode == config.ModeServer {
		utils.Warn("Running server without sudo/admin. TUN mode will fail if requested.")
	}

	utils.Info("Starting VPN in %s mode", cfg.Mode)

	if cfg.GenCerts {
		utils.Info("Generating mTLS certificates...")
		if err := crypto.GenerateCerts(); err != nil {
			utils.Error("Failed to generate certs: %v", err)
			os.Exit(1)
		}
		utils.Success("Certificates generated: ca.crt, server.crt, client.crt")
		return
	}

	ctx := context.Background()

	switch cfg.Mode {
	case config.ModeServer:
		err = tunnel.StartServer(ctx, cfg)
	case config.ModeSocks:
		err = tunnel.StartSocksClient(ctx, cfg)
	default:
		err = tunnel.StartClient(ctx, cfg)
	}

	if err != nil {
		utils.Error("Critical error: %v", err)
		os.Exit(1)
	}
}

func printBanner() {
	banner := `
%s██╗   ██╗██████╗ ███╗   ██╗%s
%s██║   ██║██╔══██╗████╗  ██║%s
%s██║   ██║██████╔╝██╔██╗ ██║%s
%s╚██╗ ██╔╝██╔═══╝ ██║╚██╗██║%s
%s ╚████╔╝ ██║     ██║ ╚████║%s
%s  ╚═══╝  ╚═╝     ╚═╝  ╚═══╝%s
    `
	fmt.Printf(banner,
		utils.Blue, utils.Reset,
		utils.Blue, utils.Reset,
		utils.Cyan, utils.Reset,
		utils.Cyan, utils.Reset,
		utils.Purple, utils.Reset,
		utils.Purple, utils.Reset,
	)
	fmt.Println("\n" + utils.Bold + "   SECURE TUNNEL PROTOTYPE" + utils.Reset)
}
