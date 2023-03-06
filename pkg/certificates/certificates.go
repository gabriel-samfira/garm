package certificates

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func writeKeyFile(path string, key interface{}) error {
	keyOut, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", path, err)
	}
	defer keyOut.Close()

	if err := pem.Encode(keyOut, pemBlockForKey(key)); err != nil {
		return fmt.Errorf("failed to write data to %s: %s", path, err)
	}
	return nil
}

func writeCertFile(path string, certBytes []byte) error {
	certOut, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", path, err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certBytes}); err != nil {
		return fmt.Errorf("failed to write data to %s: %s", path, err)
	}
	return nil
}

func newPrivateKey(keyPath string) (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, err
	}

	if err := writeKeyFile(keyPath, key); err != nil {
		return nil, err
	}
	return key, nil
}

// GenerateNewInternalCA offers a simple way to generate a new self signed CA.
// Garm will not automatically distribute certificates, but will offer some subcommands
// that aid in generating certificates for cluster members.
// These certificates are meant to be used only to secure connections between peers, for
// the embeded etcd server.
// Garm is both client and server when etcd is embeded, and these certificates are meant
// to enable secure communication between peers. The embeded etcd will not be consumed
// by any external actor.
func GenerateNewInternalCA(keyPath, certPath string) error {
	if _, err := os.Stat(keyPath); err == nil {
		return fmt.Errorf("file %s already exists", keyPath)
	}

	if _, err := os.Stat(certPath); err == nil {
		return fmt.Errorf("file %s already exists", certPath)
	}

	keyData, err := newPrivateKey(keyPath)
	if err != nil {
		return err
	}

	tpl, err := NewCertificateTemplate("GARM internal CA")
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	tpl.KeyUsage = x509.KeyUsageCertSign
	tpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
	tpl.BasicConstraintsValid = true
	tpl.IsCA = true

	derBytes, err := x509.CreateCertificate(rand.Reader, &tpl, &tpl, publicKey(keyData), keyData)
	if err != nil {
		return fmt.Errorf("creating certificate: %w", err)
	}

	if err := writeCertFile(certPath, derBytes); err != nil {
		return fmt.Errorf("writing certificate: %w", err)
	}

	return nil
}

func getPemBlockFromFile(pth string) (*pem.Block, error) {
	certContents, err := os.ReadFile(pth)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate %s: %w", pth, err)
	}

	certDecoded, _ := pem.Decode(certContents)
	if certDecoded == nil {
		return nil, fmt.Errorf("failed to decode certificate")
	}
	return certDecoded, nil
}

// LoadCertFromFile loads an x509 certificate from file.
func LoadCertFromFile(pth string) (*x509.Certificate, error) {
	certPem, err := getPemBlockFromFile(pth)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cert: %w", err)
	}
	crt, err := x509.ParseCertificate(certPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}
	return crt, nil
}

// LoadECPrivateKeyFromFile loads and ECDSA private key from file.
func LoadECPrivateKeyFromFile(pth string) (*ecdsa.PrivateKey, error) {
	keyPem, err := getPemBlockFromFile(pth)
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}

	key, err := x509.ParseECPrivateKey(keyPem.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse key: %w", err)
	}
	return key, nil
}

func NewCertificateTemplate(name string) (x509.Certificate, error) {
	now := time.Now()
	notBefore := now.Add(time.Duration(-24) * time.Hour)
	// This is not valid on all systems, but this is for internal use only.
	// notAfter := notBefore.Add(13872 * time.Hour) // 578 days
	notAfter := notBefore.Add(87600 * time.Hour) // 10 years
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return x509.Certificate{}, fmt.Errorf("failed to generate serial number: %s", err)
	}

	tpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{name},
			CommonName:   name,
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
	}
	return tpl, nil
}

func NewServerCertificateTemplate(name string, validFor []string, clientAuth bool) (x509.Certificate, error) {
	ipAddresses := []net.IP{
		net.ParseIP("127.0.0.1"),
	}
	dnsNames := []string{
		"localhost",
		"localhost.localdomain",
	}

	for _, val := range validFor {
		if ip := net.ParseIP(val); ip == nil {
			dnsNames = append(dnsNames, val)
		} else {
			ipAddresses = append(ipAddresses, ip)
		}
	}

	tpl, err := NewCertificateTemplate(name)
	if err != nil {
		return x509.Certificate{}, fmt.Errorf("failed to get template: %w", err)
	}

	extKeyUsage := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	if clientAuth {
		extKeyUsage = append(extKeyUsage, x509.ExtKeyUsageClientAuth)
	}

	tpl.IPAddresses = ipAddresses
	tpl.DNSNames = dnsNames
	tpl.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	tpl.ExtKeyUsage = extKeyUsage
	tpl.BasicConstraintsValid = true
	return tpl, nil
}

func NewClientCertificateTemplate(name string) (x509.Certificate, error) {
	tpl, err := NewCertificateTemplate(name)
	if err != nil {
		return x509.Certificate{}, fmt.Errorf("failed to get template: %w", err)
	}

	tpl.KeyUsage = x509.KeyUsageDigitalSignature
	tpl.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
	tpl.BasicConstraintsValid = true
	return tpl, nil
}

func NewCertificatesFromTemplate(tpl x509.Certificate, keyPath, certPath, caCert, caKey string) error {
	if _, err := os.Stat(keyPath); err == nil {
		return fmt.Errorf("peer key already exists")
	}

	if _, err := os.Stat(certPath); err == nil {
		return fmt.Errorf("peer certificate already exists")
	}

	parentCert, err := LoadCertFromFile(caCert)
	if err != nil {
		return fmt.Errorf("failed to load CA cert: %w", err)
	}

	parentKey, err := LoadECPrivateKeyFromFile(caKey)
	if err != nil {
		return fmt.Errorf("failed to load parent key: %w", err)
	}

	keyData, err := newPrivateKey(keyPath)
	if err != nil {
		return err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &tpl, parentCert, publicKey(keyData), parentKey)
	if err != nil {
		return fmt.Errorf("failed to create certificate: %w", err)
	}

	if err := writeCertFile(certPath, derBytes); err != nil {
		return err
	}

	return nil
}
