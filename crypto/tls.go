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
	"os"
	"time"
)

// GenerateCerts generates CA, Server, and Client certs and saves them to disk.
func GenerateCerts() error {
	// 1. CA
	caCert, caKey, caPEM, _ := createCert(true, nil, nil, "VPN Proto CA")
	if err := saveFile("ca.crt", caPEM); err != nil {
		return err
	}
	if err := saveFile("ca.key", encodeKey(caKey)); err != nil {
		return err
	}

	// 2. Server (Signed by CA)
	_, _, servPEM, servKeyPEM := createCert(false, caCert, caKey, "VPN Server")
	if err := saveFile("server.crt", servPEM); err != nil {
		return err
	}
	if err := saveFile("server.key", servKeyPEM); err != nil {
		return err
	}

	// 3. Client (Signed by CA)
	_, _, clientPEM, clientKeyPEM := createCert(false, caCert, caKey, "VPN Client")
	if err := saveFile("client.crt", clientPEM); err != nil {
		return err
	}
	if err := saveFile("client.key", clientKeyPEM); err != nil {
		return err
	}

	return nil
}

func createCert(isCA bool, parent *x509.Certificate, parentKey any, cn string) (*x509.Certificate, *rsa.PrivateKey, []byte, []byte) {
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)

	template := x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().UnixNano()),
		Subject:               pkix.Name{Organization: []string{"VPN Proto"}, CommonName: cn},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  isCA,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("0.0.0.0")},
	}

	if isCA {
		template.KeyUsage |= x509.KeyUsageCertSign
		parent = &template
		parentKey = priv
	}

	der, _ := x509.CreateCertificate(rand.Reader, &template, parent, &priv.PublicKey, parentKey)
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	return &template, priv, certPEM, keyPEM
}

func saveFile(name string, data []byte) error {
	return os.WriteFile(name, data, 0600)
}

func encodeKey(key *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
}

// LoadServerTLS loads keys and enforces Client Auth
func LoadServerTLS(certFile, keyFile, caFile string) (*tls.Config, error) {
	// Load Server Cert
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	// Load CA to verify clients
	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

// LoadClientTLS loads client cert and trusts CA
func LoadClientTLS(certFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		MinVersion:   tls.VersionTLS13,
	}, nil
}

func GetCipherSuiteName(id uint16) string {
	return tls.CipherSuiteName(id)
}
