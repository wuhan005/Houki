// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/wuhan005/Houki/internal/modules"
)

var Forward = &ForwardProxy{}

// ForwardProxy is the MitM proxy server which supports start and shutdown.
type ForwardProxy struct {
	http.Server

	enable  bool
	address string
	proxy   *goproxy.ProxyHttpServer
}

// IsEnabled returns the proxy status.
func (p *ForwardProxy) IsEnabled() bool {
	return p.enable
}

// Address returns the proxy address.
func (p *ForwardProxy) Address() string {
	return p.address
}

// Start starts the proxy.
// If the proxy has already been started, it will do nothing.
func (p *ForwardProxy) Start(address string) error {
	if p.enable {
		return nil
	}

	p.address = address
	return p.start()
}

// Shutdown shuts down the proxy server.
// If the proxy has not been started, it will return an error.
func (p *ForwardProxy) Shutdown() error {
	if !p.enable {
		return nil
	}

	return p.shutdown()
}

func (p *ForwardProxy) start() error {
	p.proxy = goproxy.NewProxyHttpServer()
	p.registerDispatcher()

	p.Server = http.Server{
		Addr:    p.address,
		Handler: p.proxy,
	}
	p.enable = true

	errChan := make(chan error)
	go func() {
		if err := p.Server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			logrus.Trace("Server closed.")
		} else if err != nil {
			errChan <- err
			p.enable = false
			logrus.WithError(err).Error("Failed to start proxy server")
		}
	}()

	select {
	case err := <-errChan:
		return err

	case <-time.After(2 * time.Second): // We trust the server has been started successfully if no error received after 2 seconds.
		logrus.WithField("address", p.address).Info("Proxy server started")
		return nil
	}
}

func (p *ForwardProxy) shutdown() error {
	err := p.Server.Shutdown(context.TODO())
	if err != nil {
		return errors.Wrap(err, "shut down")
	}

	p.enable = false
	return nil
}

// SetCA sets the goproxy server certificate globally.
func (p *ForwardProxy) SetCA(caCertPath, caKeyPath string) error {
	caCert, err := os.ReadFile(caCertPath)
	if err != nil {
		return errors.Wrap(err, "read CA certificate")
	}
	caKey, err := os.ReadFile(caKeyPath)
	if err != nil {
		return errors.Wrap(err, "read CA key")
	}

	proxyCA, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return errors.Wrap(err, "parse X509 key pair")
	}
	if proxyCA.Leaf, err = x509.ParseCertificate(proxyCA.Certificate[0]); err != nil {
		return errors.Wrap(err, "parse certificate")
	}

	goproxy.GoproxyCa = proxyCA
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&proxyCA)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&proxyCA)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&proxyCA)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&proxyCA)}

	return nil
}

// registerDispatcher registers the MitM server's dispatcher, which used to modify the request and response message.
func (p *ForwardProxy) registerDispatcher() {
	p.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	p.proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			logrus.WithError(err).Error("Failed to read request body")
		}
		// Set the body back to the request.
		req.Body = io.NopCloser(bytes.NewReader(body))

		modules.DoRequest(req, body)

		return req, nil
	})

	p.proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if ctx.Resp == nil || ctx.Resp.Header.Get("Content-Type") == "text/event-stream" {
			return ctx.Resp
		}

		body, err := io.ReadAll(ctx.Resp.Body)
		if err != nil {
			logrus.WithError(err).Error("Failed to read response body")
		}
		// Set the body back to the response.
		ctx.Resp.Body = io.NopCloser(bytes.NewReader(body))

		modules.DoResponse(ctx.Resp, body)

		return ctx.Resp
	})
}
