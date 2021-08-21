// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"bytes"
	"io"
	"net/http"

	"github.com/elazarl/goproxy"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/modules"
)

// registerDispatcher registers the MitM server's dispatcher, which used to modify the request and response message.
func (p *Proxy) registerDispatcher() {
	p.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)

	p.proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Error("Failed to read request body: %v", err)
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
			log.Error("Failed to read response body: %v", err)
		}
		// Set the body back to the response.
		ctx.Resp.Body = io.NopCloser(bytes.NewReader(body))

		modules.DoResponse(ctx.Resp, body)

		return ctx.Resp
	})
}
