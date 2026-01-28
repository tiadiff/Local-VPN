package config

import (
	"flag"
	"fmt"
)

type Mode string

const (
	ModeServer Mode = "server"
	ModeClient Mode = "client"
	ModeSocks  Mode = "socks"
)

type Config struct {
	Mode       Mode
	ServerAddr string
	Port       int
	SocksPort  int
	Secret     string // Shared secret for simple authentication
}

func Load() (*Config, error) {
	modeStr := flag.String("mode", "", "Mode of operation: 'server', 'client', or 'socks'")
	serverAddr := flag.String("server", "127.0.0.1", "Server address (for client mode)")
	port := flag.Int("port", 3000, "Port to listen on (server) or connect to (client)")
	socksPort := flag.Int("socks", 1080, "Local SOCKS5 port (for socks mode)")
	secret := flag.String("secret", "default-secret", "Shared secret for authentication")

	flag.Parse()

	var mode Mode
	switch *modeStr {
	case "server":
		mode = ModeServer
	case "client":
		mode = ModeClient
	case "socks":
		mode = ModeSocks
	default:	
        // Default to client if not specified but server arg is present, else usage
		return nil, fmt.Errorf("invalid mode: %s", *modeStr)
	}

	return &Config{
		Mode:       mode,
		ServerAddr: *serverAddr,
		Port:       *port,
		SocksPort:  *socksPort,
		Secret:     *secret,
	}, nil
}
