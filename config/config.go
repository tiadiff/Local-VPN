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
	Secret     string // Shared secret (Secondary, can be deprecated)
	GenCerts   bool
	CertFile   string
	KeyFile    string
	CAFile     string
}

func Load() (*Config, error) {
	modeStr := flag.String("mode", "", "Mode of operation: 'server', 'client', or 'socks'")
	serverAddr := flag.String("server", "127.0.0.1", "Server address (for client mode)")
	port := flag.Int("port", 3000, "Port to listen on (server) or connect to (client)")
	socksPort := flag.Int("socks", 1080, "Local SOCKS5 port (for socks mode)")
	secret := flag.String("secret", "default-secret", "Shared secret for authentication")

	genCerts := flag.Bool("gen-certs", false, "Generate new mTLS certificates and exit")
	certFile := flag.String("cert", "server.crt", "Certificate file (server.crt or client.crt)")
	keyFile := flag.String("key", "server.key", "Private Key file")
	caFile := flag.String("ca", "ca.crt", "CA Certificate file")

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
		if *genCerts {
			// UI or GenCerts don't strictly require mode
			return &Config{
				GenCerts:   *genCerts,
				Port:       *port,
				SocksPort:  *socksPort,
				ServerAddr: *serverAddr,
				CertFile:   *certFile,
				KeyFile:    *keyFile,
				CAFile:     *caFile,
			}, nil
		}
		return nil, fmt.Errorf("invalid or missing mode: %s", *modeStr)
	}

	return &Config{
		Mode:       mode,
		ServerAddr: *serverAddr,
		Port:       *port,
		SocksPort:  *socksPort,
		Secret:     *secret,
		GenCerts:   *genCerts,
		CertFile:   *certFile,
		KeyFile:    *keyFile,
		CAFile:     *caFile,
	}, nil
}
