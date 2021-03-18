package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

const (
	CommonName       = "Houki CA"
	OrganizationName = "Houki"
)

func Init() {
	_, err := os.Stat(".certificate")
	if os.IsNotExist(err) {
		err := os.Mkdir(".certificate", 0644)
		if err != nil {
			log.Fatal("Failed to create .certificate folder: %v", err)
		}
		log.Info("Create .certificate folder.")
	}
}

func Get() ([]byte, []byte, error) {
	crt, key, err := readFromFile()
	if err == nil {
		log.Trace("Read CA from file.")
		return crt, key, nil
	}

	return Generate(true)
}

func Generate(save bool) ([]byte, []byte, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, errors.Wrap(err, "generate serial number")
	}

	rootKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	b, err := x509.MarshalECPrivateKey(rootKey)
	if err != nil {
		return nil, nil, errors.Wrap(err, "marshal ECDSA private key")
	}

	rootKeyBytes := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b})
	rootTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   CommonName,
			Organization: []string{OrganizationName},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(5 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &rootKey.PublicKey, rootKey)
	rootCrtBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return nil, nil, err
	}

	if save {
		err := saveToFile(rootCrtBytes, rootKeyBytes)
		log.Error("Failed to save CA to file: %v", err)
	}

	return rootCrtBytes, rootKeyBytes, nil
}

func readFromFile() ([]byte, []byte, error) {
	crt, err := os.ReadFile(".certificate/ca.crt")
	if err != nil {
		return nil, nil, err
	}

	key, err := os.ReadFile(".certificate/ca.key")
	if err != nil {
		return nil, nil, err
	}

	return crt, key, nil
}

func saveToFile(crt, key []byte) error {
	err := os.WriteFile(".certificate/ca.crt", crt, 0644)
	if err != nil {
		return err
	}

	return os.WriteFile(".certificate/ca.key", key, 0644)
}
