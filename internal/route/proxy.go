// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/wuhan005/Houki/internal/context"
	"github.com/wuhan005/Houki/internal/db"
	"github.com/wuhan005/Houki/internal/form"
	"github.com/wuhan005/Houki/internal/modules"
	"github.com/wuhan005/Houki/internal/proxy"
)

type ProxyHandler struct{}

func NewProxyHandler() *ProxyHandler {
	return &ProxyHandler{}
}

func (*ProxyHandler) Status(ctx context.Context) error {
	return ctx.Success(map[string]interface{}{
		"forward": map[string]interface{}{
			"enabled": proxy.Forward.IsEnabled(),
			"address": proxy.Forward.Address(),
		},
		"reverse": map[string]interface{}{
			"enabled": proxy.Reverse.IsEnabled(),
			"address": proxy.Reverse.Address(),
		},
	})
}

func (*ProxyHandler) reloadAllModules(ctx context.Context) error {
	enabledModules, err := db.Modules.All(ctx.Request().Context(), db.AllModuleOptions{
		EnabledOnly: true,
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to list modules")
		return ctx.ServerError()
	}

	moduleSets := make(map[uint]*modules.Body, len(enabledModules))
	for _, module := range enabledModules {
		moduleSets[module.ID] = module.Body
	}

	if err := modules.ReloadAllModules(moduleSets); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to reload modules")
		return ctx.ServerError()
	}
	return nil
}

func (h *ProxyHandler) StartForward(ctx context.Context, f form.StartProxy) error {
	_ = h.reloadAllModules(ctx)
	if ctx.ResponseWriter().Written() {
		return nil
	}

	if err := proxy.Forward.Start(f.Address); err != nil {
		return ctx.Error(http.StatusInternalServerError, "Failed to start proxy: %v", err)
	}
	return ctx.Success("Forward proxy started successfully")
}

func (*ProxyHandler) ShutdownForward(ctx context.Context) error {
	if err := proxy.Forward.Shutdown(); err != nil {
		return ctx.Error(http.StatusInternalServerError, "Failed to shutdown proxy: %v", err)
	}
	return ctx.Success("Forward proxy shutdown successfully")
}

func (h *ProxyHandler) StartReverse(ctx context.Context, f form.StartProxy) error {
	_ = h.reloadAllModules(ctx)
	if ctx.ResponseWriter().Written() {
		return nil
	}

	//if err := proxy.Reverse.Start(f.Address); err != nil {
	//	return ctx.Error(http.StatusInternalServerError, "Failed to start proxy: %v", err)
	//}

	return ctx.Success("Reverse proxy started successfully")
}

func (*ProxyHandler) ShutdownReverse(ctx context.Context) error {
	return ctx.Success("Reverse proxy shutdown successfully")
}
