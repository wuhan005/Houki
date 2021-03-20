package proxy

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/pkg/errors"
	"github.com/wuhan005/Houki/internal/ca"
	log "unknwon.dev/clog/v2"
)

type Proxy struct {
	server *http.Server
	proxy  *goproxy.ProxyHttpServer
}

func New() (*Proxy, error) {
	p := &Proxy{
		server: &http.Server{},
		proxy:  goproxy.NewProxyHttpServer(),
	}

	caCrt, caKey, err := ca.Get()
	if err != nil {
		return nil, err
	}

	if err := p.SetCA(caCrt, caKey); err != nil {
		return nil, errors.Wrap(err, "set CA")
	}

	return p, nil
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

func (p *Proxy) Run(addr string) {
	p.server.Addr = addr
	p.server.Handler = p.proxy

	go func() {
		if err := p.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("Failed to start proxy server: %v", err)
		}
	}()
	log.Info("Proxy server listening on %s", addr)
}

func (p *Proxy) Stop() error {
	return p.server.Shutdown(nil)
}
