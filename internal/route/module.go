// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/flamego/flamego"
	"github.com/flamego/template"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/context"
	"github.com/wuhan005/Houki/internal/db"
	"github.com/wuhan005/Houki/internal/form"
	"github.com/wuhan005/Houki/internal/module"
	"github.com/wuhan005/Houki/internal/proxy"
)

type ModulesHandler struct{}

func NewModulesHandler() *ModulesHandler {
	return &ModulesHandler{}
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
	return func(ctx context.Context) {
		moduleID := ctx.Params("id")
		mod, err := db.Modules.Get(ctx.Request().Context(), moduleID)
		if err != nil {
			if errors.Is(err, db.ErrModuleNotFound) {
				ctx.Error(40400, "Module not found")
				return
			}
			log.Error("Failed to get module: %v", err)
			ctx.ServerError()
			return
		}

		err = db.Modules.SetStatus(ctx.Request().Context(), mod.ID, status == Enable)
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

func (*ModulesHandler) New(ctx context.Context, t template.Template) {
	t.HTML(http.StatusOK, "new_module")
}

func (*ModulesHandler) NewAction(ctx context.Context, f form.NewModule) {
	var body module.Body
	if err := json.Unmarshal(f.Body, &body); err != nil {
		ctx.Error(40000, "Failed to parse module body: %v", err)
		return
	}

	if f.ID == "" {
		ctx.Error(40000, "Module ID is required")
		return
	}

	if !regexp.MustCompile("^[a-zA-Z0-9_-]+$").MatchString(f.ID) {
		ctx.Error(40000, "Module ID can only contain letters, numbers, underscores and dashes")
		return
	}

	err := db.Modules.Create(ctx.Request().Context(), db.CreateModuleOptions{
		ID:   f.ID,
		Body: &body,
	})
	if err != nil {
		if errors.Is(err, db.ErrModuleExists) {
			ctx.Error(40000, "Module already exists")
			return
		}

		ctx.Error(50000, fmt.Sprintf("Failed to create new module: %v", err))
		return
	}
	ctx.Success("success")
}

func (*ModulesHandler) Get(ctx context.Context, t template.Template, data template.Data) {
	id := ctx.Params("id")

	mod, err := db.Modules.Get(ctx.Request().Context(), id)
	if err != nil {
		ctx.Redirect("/")
		return
	}

	data["Module"] = mod

	var moduleBody bytes.Buffer
	encoder := json.NewEncoder(&moduleBody)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(mod.Body)

	data["ModuleBody"] = moduleBody.String()

	t.HTML(http.StatusOK, "update_module")
}

func (*ModulesHandler) Update(ctx context.Context, f form.UpdateModule) {
	id := ctx.Params("id")

	var body module.Body
	if err := json.Unmarshal(f.Body, &body); err != nil {
		ctx.Error(40000, "Failed to parse module body: %v", err)
		return
	}

	if err := db.Modules.Update(ctx.Request().Context(), id, db.UpdateModuleOptions{
		Body: &body,
	}); err != nil {
		log.Error("Failed to update module: %v", err)
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

func (*ModulesHandler) Delete(ctx context.Context) {
	id := ctx.Params("id")
	if err := db.Modules.Delete(ctx.Request().Context(), id); err != nil {
		log.Error("Failed to delete module: %v", err)
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
