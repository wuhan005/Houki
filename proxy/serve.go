package proxy

import (
	"net/http"

	"github.com/elazarl/goproxy"
)

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.server.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	
	p.server.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		return req, nil
	})

	p.server.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		return resp
	})

	p.server.ServeHTTP(w, r)
}
