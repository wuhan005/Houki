// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/wuhan005/Houki/internal/ca"
	"github.com/wuhan005/Houki/internal/context"
	"github.com/wuhan005/Houki/internal/form"
	"github.com/wuhan005/Houki/internal/proxy"
)

type CertificateHandler struct{}

func NewCertificateHandler() *CertificateHandler {
	return &CertificateHandler{}
}

func (*CertificateHandler) Get(ctx context.Context) error {
	cert, key, err := ca.Get()
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get certificate")
		return ctx.ServerError()
	}

	metadata, err := ca.Parse(cert)
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to parse certificate")
		return ctx.ServerError()
	}

	return ctx.Success(map[string]interface{}{
		"certificate": string(cert),
		"privateKey":  string(key),
		"metadata":    metadata,
	})
}

func (*CertificateHandler) Update(ctx context.Context, f form.UpdateCertificate) error {
	// Check if the certificate is valid.
	if _, err := ca.Parse([]byte(f.Certificate)); err != nil {
		return ctx.Error(http.StatusBadRequest, "Invalid certificate: %v", err)
	}

	if err := ca.Save([]byte(f.Certificate), []byte(f.PrivateKey)); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to save certificate")
		return ctx.ServerError()
	}
	ca.CleanCache()

	if err := proxy.Forward.SetCA(ca.CertificatePath, ca.KeyPath); err != nil {
		return ctx.Error(http.StatusInternalServerError, "Failed to set CA: %v", err)
	}

	return ctx.Success("Certificate updated successfully")
}
