package proxy

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/ca"
)

var proxy *Proxy

type Proxy struct {
	http.Server

	enable bool
	proxy  *goproxy.ProxyHttpServer
}

func Initialize() (*Proxy, error) {
	p := &Proxy{
		proxy: goproxy.NewProxyHttpServer(),
	}

	caCrt, caKey, err := ca.Get()
	if err != nil {
		return nil, err
	}

	if err := p.SetCA(caCrt, caKey); err != nil {
		return nil, errors.Wrap(err, "set CA")
	}

	proxy = p
	return p, nil
}

func IsEnable() bool {
	return proxy.isEnable()
}

func Start(addr string) {
	if proxy.enable {
		return
	}
	proxy.run(addr)
}

func Stop() error {
	if !proxy.enable {
		return errors.New("Proxy server has been started.")
	}
	return proxy.stop()
}

func (p *Proxy) SetCA(caCert, caKey []byte) error {
	proxyCA, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return err
	}
	if proxyCA.Leaf, err = x509.ParseCertificate(proxyCA.Certificate[0]); err != nil {
		return err
	}

	goproxy.GoproxyCa = proxyCA
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&proxyCA)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&proxyCA)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&proxyCA)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&proxyCA)}

	return nil
}

func (p *Proxy) run(addr string) {
	p.serve()

	p.Server = http.Server{
		Addr:    addr,
		Handler: p.proxy,
	}
	p.enable = true

	go func() {
		if err := p.Server.ListenAndServe(); err == http.ErrServerClosed {
			log.Trace("Server closed.")
		} else if err != nil {
			p.enable = false
			log.Error("Failed to start proxy server: %v", err)
		}
	}()

	log.Info("Proxy server listening on %s", addr)
}

func (p *Proxy) stop() error {
	err := p.Server.Shutdown(context.TODO())
	if err != nil {
		return errors.Wrap(err, "shut down")
	}

	p.enable = false
	return nil
}

func (p *Proxy) isEnable() bool {
	return p.enable
}
