package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

// GenerateSelfSignedCert generates a self-signed certificate and private key
// returning the TLS config to use it.
func GenerateSelfSignedCert() (*tls.Config, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"VPN Proto"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

    // Add localhost and 127.0.0.1 for local testing
	template.IPAddresses = append(template.IPAddresses, net.ParseIP("127.0.0.1"))
    template.DNSNames = append(template.DNSNames, "localhost")

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}, nil
}

// ClientTLSConfig returns a TLS config for the client.
// In a real app we would load CA certs, here we skip verification for the prototype
// or we could pin the certificate.
func ClientTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true, // For prototype only
		MinVersion:         tls.VersionTLS13,
	}
}

func GetCipherSuiteName(id uint16) string {
	return tls.CipherSuiteName(id)
}
