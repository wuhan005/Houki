// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/wuhan005/Houki/internal/ca"
	"github.com/wuhan005/Houki/internal/modules"
)

var Reverse = &ReverseProxy{}

// ReverseProxy is the reverse proxy server which supports start and shutdown.
type ReverseProxy struct {
	http.Server

	enable  bool
	address string
	isHttps bool
	proxy   *httputil.ReverseProxy
}

// IsEnabled returns the proxy status.
func (p *ReverseProxy) IsEnabled() bool {
	return p.enable
}

// Address returns the proxy address.
func (p *ReverseProxy) Address() string {
	return p.address
}

// Start starts the proxy.
// If the proxy has already been started, it will do nothing.
func (p *ReverseProxy) Start(address string) error {
	if p.enable {
		return nil
	}

	p.address = address
	p.isHttps = strings.HasSuffix(p.address, ":443")
	return p.start()
}

// Shutdown shuts down the proxy server.
// If the proxy has not been started, it will return an error.
func (p *ReverseProxy) Shutdown() error {
	if !p.enable {
		return nil
	}

	return p.shutdown()
}

func (p *ReverseProxy) start() error {
	p.proxy = &httputil.ReverseProxy{}
	p.registerDispatcher()

	p.Server = http.Server{
		Addr:    p.address,
		Handler: p.proxy,
		TLSConfig: &tls.Config{
			GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
				caCertBytes, caKeyBytes, err := ca.Get()
				if err != nil {
					return nil, errors.Wrap(err, "get ca cert")
				}
				caCert, err := tls.X509KeyPair(caCertBytes, caKeyBytes)
				if err != nil {
					return nil, errors.Wrap(err, "parse x509 key pair")
				}

				host := chi.ServerName
				selfSignedCert, err := ca.SignHost(caCert, []string{host})
				if err != nil {
					return nil, errors.Wrap(err, "sign host")
				}
				return selfSignedCert, nil
			},
		},
	}
	p.enable = true

	errChan := make(chan error)
	go func() {
		if p.isHttps {
			if err := p.Server.ListenAndServeTLS("", ""); errors.Is(err, http.ErrServerClosed) {
				logrus.Trace("Server closed.")
			} else if err != nil {
				errChan <- err
				p.enable = false
				logrus.WithError(err).Error("Failed to start proxy server")
			}
		} else {
			if err := p.Server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
				logrus.Trace("Server closed.")
			} else if err != nil {
				errChan <- err
				p.enable = false
				logrus.WithError(err).Error("Failed to start proxy server")
			}
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

func (p *ReverseProxy) shutdown() error {
	if err := p.Server.Shutdown(context.TODO()); err != nil {
		return errors.Wrap(err, "shut down")
	}

	p.enable = false
	return nil
}

// registerDispatcher registers the reverse proxy server's dispatcher, which used to modify the request and response message.
func (p *ReverseProxy) registerDispatcher() {
	p.proxy.Director = func(req *http.Request) {
		req.URL.Host = req.Host

		if p.isHttps {
			req.URL.Scheme = "https"
		} else {
			req.URL.Scheme = "http"
		}

		var body []byte
		if req.Body != nil {
			var err error
			body, err = io.ReadAll(req.Body)
			if err != nil {
				logrus.WithError(err).Error("Failed to read request body")
			}
			// Set the body back to the request.
			req.Body = io.NopCloser(bytes.NewReader(body))
		}

		modules.DoRequest(req, body)
	}

	p.proxy.ModifyResponse = func(resp *http.Response) error {
		if resp == nil || resp.Header.Get("Content-Type") == "text/event-stream" {
			return nil
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.WithError(err).Error("Failed to read response body")
		}
		// Set the body back to the response.
		resp.Body = io.NopCloser(bytes.NewReader(body))

		modules.DoResponse(resp, body)

		return nil
	}

	p.proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}
