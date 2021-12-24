// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	"github.com/flamego/flamego"
	"github.com/flamego/template"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/context"
	"github.com/wuhan005/Houki/internal/db"
	"github.com/wuhan005/Houki/internal/form"
)

type ModulesHandler struct{}

func NewModulesHandler() *ModulesHandler {
	return &ModulesHandler{}
}

func (*ModulesHandler) List(ctx context.Context, t template.Template, data template.Data) {
	modules, err := db.Modules.List(ctx.Request().Context(), db.GetModuleOptions{})
	if err != nil {
		log.Error("Failed to list modules: %v", err)
		ctx.ServerError()
		return
	}
	data["modules"] = modules
	t.HTML(http.StatusOK, "modules")
}

const (
	Enable  = "enable"
	Disable = "disable"
)

func (*ModulesHandler) SetStatus(status string) flamego.Handler {
	return func(ctx context.Context) error {
		//moduleID := ctx.Params("id")
		//module, err := db.Modules.Get(ctx.Request().Context(), moduleID)
		//if err != nil {
		//	return ctx.Error(50000, fmt.Sprintf("Failed to get module: %v", err))
		//}
		//
		//err = db.Modules.Update(ctx.Request().Context(), module.ID, db.UpdateModuleOptions{
		//	FilePath: module.FilePath,
		//	Enabled:  status == Enable,
		//})
		//if err != nil {
		//	return ctx.Error(50000, fmt.Sprintf("Failed to enable module: %v", err))
		//}
		//
		//return ctx.Success()
		return nil
	}
}

func (*ModulesHandler) New(ctx context.Context, t template.Template, data template.Data) {
	t.HTML(http.StatusOK, "new_module")
}

func (*ModulesHandler) NewAction(ctx context.Context, f form.NewModule) {
	//err := db.Modules.Create(ctx.Request().Context(), db.CreateModuleOptions{
	//	ID:       f.ID,
	//	FilePath: f.FilePath,
	//})
	//if err != nil {
	//	return ctx.Error(50000, fmt.Sprintf("Failed to create new module: %v", err))
	//}
	//return ctx.Success()
}

func (*ModulesHandler) Upload(ctx context.Context) error {
	panic("implement me")
}

func (*ModulesHandler) Download(ctx context.Context) error {
	panic("implement me")
}

func (*ModulesHandler) Edit(ctx context.Context) error {
	panic("implement me")
}
