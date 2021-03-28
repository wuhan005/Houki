package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/elazarl/goproxy"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/module"
	"github.com/wuhan005/Houki/internal/sse"
)

func (p *Proxy) serve() {
	p.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	p.proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Error("Failed to read request body: %v", err)
		}
		req.Body = io.NopCloser(bytes.NewReader(body))

		// Send log
		err = sse.Write("request", map[string]interface{}{
			"method": req.Method,
			"url":    req.URL.String(),
		})
		if err != nil {
			log.Error("Failed to send log: %v", err)
		}

		module.DoRequest(req, body)

		return req, nil
	})

	p.proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error("Failed to read response body: %v", err)
		}
		resp.Body = io.NopCloser(bytes.NewReader(body))

		module.DoResponse(resp, body)

		return resp
	})
}
