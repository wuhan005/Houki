package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/elazarl/goproxy"
	"github.com/wuhan005/Houki/internal/module"
	log "unknwon.dev/clog/v2"
)

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.server.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	p.server.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Error("Failed to read request body: %v", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(body))

		module.DoRequest(req, body)

		return req, nil
	})

	p.server.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("Failed to read response body: %v", err)
		}
		resp.Body = io.NopCloser(bytes.NewReader(body))

		module.DoResponse(resp, body)

		return resp
	})

	p.server.ServeHTTP(w, r)
}
