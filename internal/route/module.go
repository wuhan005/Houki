// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	"github.com/flamego/flamego"
	"github.com/flamego/template"

	"github.com/wuhan005/Houki/internal/context"
	"github.com/wuhan005/Houki/internal/form"
	"github.com/wuhan005/Houki/internal/modules"
)

type ModulesHandler struct{}

func NewModulesHandler() *ModulesHandler {
	return &ModulesHandler{}
}

func (*ModulesHandler) List(ctx context.Context, t template.Template, data template.Data) {
	modules := modules.List()
	data["modules"] = modules
	t.HTML(http.StatusOK, "modules")
}

func (*ModulesHandler) Get(ctx context.Context) error {
	//moduleID := ctx.Params("id")
	//module, err := modules.Get(moduleID)
	//if err != nil {
	//	return ctx.Error(40400, "Module not found.")
	//}
	//
	//return ctx.Success(module)
	return nil
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

func (*ModulesHandler) New(ctx context.Context, f form.NewModule) {
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
