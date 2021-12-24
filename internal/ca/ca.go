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

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"
)

const (
	CommonName       = "Houki CA"
	OrganizationName = "Houki"
)

func init() {
	if err := createFolder(); err != nil {
		log.Fatal("Failed to create certificate folder: %v", err)
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

// Get trys to read certificate from file.
// It will create the certificate if fails.
func Get() ([]byte, []byte, error) {
	crt, key, err := readFromFile()
	if err == nil {
		return crt, key, nil
	}

	return GenerateCertificate(true)
}

//func GetPin() ([]byte, error) {
//	certBytes, _, err := Get()
//	if err != nil {
//		return nil, errors.Wrap(err, "get certificate")
//	}
//
//	cert, err := x509.ParseCertificate(certBytes)
//	if err != nil {
//		return nil, errors.Wrap(err, "x509.ParseCertificate")
//	}
//	cert.PublicKey.(*ecdsa.PublicKey)
//
//	publicDer, err := x509.MarshalPKCS1PublicKey()
//	if err != nil {
//		return nil, errors.Wrap(err, "x509.MarshalPKIXPublicKey")
//	}
//	sum := sha256.Sum256(publicDer)
//	pin := make([]byte, base64.StdEncoding.EncodedLen(len(sum)))
//	base64.StdEncoding.Encode(pin, sum[:])
//
//	return pin, nil
//}

func GenerateCertificate(save bool) ([]byte, []byte, error) {
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
		if err != nil {
			log.Error("Failed to save CA to file: %v", err)
		}
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
