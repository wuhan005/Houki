// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"bytes"
	"io"
	"net/http"
	"time"

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
			"host":   req.URL.Host,
			"path":   req.URL.Path,
			"time":   time.Now().Unix(),
		})
		if err != nil {
			log.Error("Failed to send log: %v", err)
		}

		module.DoRequest(req, body)

		return req, nil
	})

	p.proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if ctx.Resp == nil || ctx.Resp.Header.Get("Content-Type") == "text/event-stream" {
			return ctx.Resp
		}

		body, err := io.ReadAll(ctx.Resp.Body)
		if err != nil {
			log.Error("Failed to read response body: %v", err)
		}
		ctx.Resp.Body = io.NopCloser(bytes.NewReader(body))

		module.DoResponse(ctx.Resp, body)

		return ctx.Resp
	})
}
