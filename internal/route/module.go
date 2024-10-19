// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"encoding/json"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/context"
	"github.com/wuhan005/Houki/internal/db"
	"github.com/wuhan005/Houki/internal/dbutil"
	"github.com/wuhan005/Houki/internal/form"
	"github.com/wuhan005/Houki/internal/modules"
	"github.com/wuhan005/Houki/internal/proxy"
)

type ModulesHandler struct{}

func NewModulesHandler() *ModulesHandler {
	return &ModulesHandler{}
}

func (*ModulesHandler) List(ctx context.Context) error {
	modules, total, err := db.Modules.List(ctx.Request().Context(), db.ListModuleOptions{
		EnabledOnly: ctx.QueryBool("enabled"),
		Pagination: dbutil.Pagination{
			Page:     ctx.QueryInt("page"),
			PageSize: ctx.QueryInt("pageSize"),
		},
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to list modules")
		return ctx.ServerError()
	}

	return ctx.Success(map[string]interface{}{
		"modules": modules,
		"total":   total,
	})
}

func (*ModulesHandler) Create(ctx context.Context, f form.CreateModule) error {
	var body modules.Body
	if err := json.Unmarshal(f.Body, &body); err != nil {
		return ctx.Error(http.StatusBadRequest, "Failed to parse module body: %v", err)
	}

	module, err := db.Modules.Create(ctx.Request().Context(), db.CreateModuleOptions{
		Name: f.Name,
		Body: &body,
	})
	if err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to create module")
		return ctx.ServerError()
	}

	return ctx.Success(module)
}

func (*ModulesHandler) Moduler(ctx context.Context) error {
	moduleID := uint(ctx.ParamInt("id"))

	module, err := db.Modules.Get(ctx.Request().Context(), moduleID)
	if err != nil {
		if errors.Is(err, db.ErrModuleNotFound) {
			return ctx.Error(http.StatusNotFound, "Module not found")
		}
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to get module")
		return ctx.ServerError()
	}

	ctx.Map(module)
	return nil
}

func (*ModulesHandler) RefreshModule(ctx context.Context) {
	if err := proxy.ReloadAllModules(ctx.Request().Context()); err != nil {
		log.Error("Failed to reload modules: %v", err)
		ctx.ServerError()
		return
	}
	ctx.Success("success")
}

const (
	Enable  = "enable"
	Disable = "disable"
)

func (*ModulesHandler) SetStatus(status string) flamego.Handler {
	return func(ctx context.Context, module *db.Module) {
		err := db.Modules.SetStatus(ctx.Request().Context(), module.ID, status == Enable)
		if err != nil {
			log.Error("Failed to set module status: %v", err)
			ctx.ServerError()
			return
		}

		if err := proxy.ReloadAllModules(ctx.Request().Context()); err != nil {
			log.Error("Failed to reload modules: %v", err)
			ctx.ServerError()
			return
		}

		ctx.Success("success")
	}
}

func (*ModulesHandler) Get(ctx context.Context, module *db.Module) error {
	return ctx.Success(module)
}

func (*ModulesHandler) Update(ctx context.Context, module *db.Module, f form.UpdateModule) error {
	var body modules.Body
	if err := json.Unmarshal(f.Body, &body); err != nil {
		return ctx.Error(http.StatusBadRequest, "Failed to parse module body: %v", err)
	}

	if err := db.Modules.Update(ctx.Request().Context(), module.ID, db.UpdateModuleOptions{
		Name: f.Name,
		Body: &body,
	}); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to update module")
		return ctx.ServerError()
	}
	if err := proxy.ReloadAllModules(ctx.Request().Context()); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to reload all modules")
		return ctx.ServerError()
	}

	return ctx.Success("Update module successfully")
}

func (*ModulesHandler) Delete(ctx context.Context, module *db.Module) error {
	if err := db.Modules.Delete(ctx.Request().Context(), module.ID); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to delete module")
		return ctx.ServerError()
	}

	if err := proxy.ReloadAllModules(ctx.Request().Context()); err != nil {
		logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to reload all modules")
		return ctx.ServerError()
	}

	return ctx.Success("Delete module successfully")
}
