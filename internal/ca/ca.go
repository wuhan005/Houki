// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

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

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	CommonName       = "Houki CA"
	OrganizationName = "Houki"
	CertificatePath  = ".certificate/ca.crt"
	KeyPath          = ".certificate/ca.key"
)

func init() {
	if err := createFolder(); err != nil {
		logrus.WithError(err).Fatal("Failed to create certificate folder")
	}
}

func createFolder() error {
	_, err := os.Stat(".certificate")
	if os.IsNotExist(err) {
		err := os.Mkdir(".certificate", 0755)
		if err != nil {
			return errors.Wrap(err, "create .certificate folder")
		}
	}
	return nil
}

var cacheStorage = cache.New[string, []byte]()

// Get trys to read certificate from file.
// It will create the certificate if fails.
func Get() ([]byte, []byte, error) {
	crt, crtExists := cacheStorage.Get(CertificatePath)
	key, keyExists := cacheStorage.Get(KeyPath)
	if crtExists && keyExists {
		return crt, key, nil
	}

	crt, key, err := readFromFile()
	if err == nil {
		return crt, key, nil
	}
	return nil, nil, errors.New("certificate not found")
}

func CleanCache() {
	cacheStorage.Delete(CertificatePath)
	cacheStorage.Delete(KeyPath)
}

func Generate() ([]byte, []byte, error) {
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

	return rootCrtBytes, rootKeyBytes, nil
}

func Save(crt, key []byte) error {
	if err := os.WriteFile(CertificatePath, crt, 0644); err != nil {
		return errors.Wrap(err, "write certificate")
	}
	if err := os.WriteFile(KeyPath, key, 0644); err != nil {
		return errors.Wrap(err, "write key")
	}
	return nil
}

type Metadata struct {
	Issuer             string    `json:"issuer"`
	ValidFrom          time.Time `json:"validFrom"`
	ValidTo            time.Time `json:"validTo"`
	PublicKeyAlgorithm string    `json:"publicKeyAlgorithm"`
	SerialNumber       string    `json:"serialNumber"`
	SignatureAlgorithm string    `json:"signatureAlgorithm"`
}

func Parse(crt []byte) (*Metadata, error) {
	block, _ := pem.Decode(crt)
	if block == nil {
		return nil, errors.New("failed to parse certificate")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "parse certificate")
	}

	metadata := &Metadata{
		Issuer:             cert.Issuer.String(),
		ValidFrom:          cert.NotBefore,
		ValidTo:            cert.NotAfter,
		PublicKeyAlgorithm: cert.PublicKeyAlgorithm.String(),
		SerialNumber:       cert.SerialNumber.String(),
		SignatureAlgorithm: cert.SignatureAlgorithm.String(),
	}
	return metadata, nil
}

func readFromFile() ([]byte, []byte, error) {
	crt, err := os.ReadFile(CertificatePath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "read certificate")
	}

	key, err := os.ReadFile(KeyPath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "read key")
	}

	return crt, key, nil
}
