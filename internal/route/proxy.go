// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	"github.com/flamego/template"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/context"
	"github.com/wuhan005/Houki/internal/db"
	"github.com/wuhan005/Houki/internal/form"
	"github.com/wuhan005/Houki/internal/proxy"
)

type ProxyHandler struct{}

func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{}
}

func (*ProxyHandler) Dashboard(ctx context.Context, t template.Template, data template.Data) {
	modules, err := db.Modules.List(ctx.Request().Context(), db.GetModuleOptions{})
	if err != nil {
		log.Error("Failed to list modules: %v", err)
		ctx.ServerError()
		return
	}
	data["Modules"] = modules
	data["Enabled"] = proxy.IsEnabled()
	t.HTML(http.StatusOK, "proxy")
}

func (*ProxyHandler) Start(ctx context.Context, f form.StartProxy) {
	if f.Address == "" {
		ctx.Error(40000, "Proxy address is required")
		return
	}

	if err := proxy.Start(f.Address); err != nil {
		log.Error("Failed to start proxy: %v", err)
		return
	}

	if err := proxy.ReloadAllModules(ctx.Request().Context()); err != nil {
		log.Error("Failed to reload modules: %v", err)
		return
	}

	ctx.Success("success")
}

func (*ProxyHandler) ShutDown(ctx context.Context) {
	defer func() { ctx.Redirect("/proxy/") }()

	err := proxy.Shutdown()
	if err != nil {
		log.Error("Failed to shutdown proxy: %v", err)
	}
}
