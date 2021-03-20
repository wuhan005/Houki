package proxy

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/elazarl/goproxy"
	"github.com/pkg/errors"
	"github.com/wuhan005/Houki/internal/ca"
)

type Proxy struct {
	server *goproxy.ProxyHttpServer
}

func New() (*Proxy, error) {
	p := &Proxy{
		server: goproxy.NewProxyHttpServer(),
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
