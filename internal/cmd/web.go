// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/flamego/flamego"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/wuhan005/Houki/internal/context"
	"github.com/wuhan005/Houki/internal/db"
	"github.com/wuhan005/Houki/internal/route"
)

var Web = &cli.Command{
	Name:        "web",
	Usage:       "Start web server",
	Description: "",
	Action:      runWeb,
	Flags: []cli.Flag{
		stringFlag("port, p", "8000", "Temporary port number to prevent conflict"),
	},
}

func runWeb(c *cli.Context) error {
	if _, err := db.New(); err != nil {
		return errors.Wrap(err, "new db")
	}

	f := flamego.Classic()
	f.Use(
		context.Contexter(),
	)

	f.Get("/", func(ctx context.Context) {
		ctx.Redirect("/proxy/")
	})

	proxy := route.NewProxyHandler()
	f.Group("/proxy", func() {
		f.Get("/status", proxy.Dashboard)

		f.Group("/forward", func() {
			f.Post("/start")
			f.Post("/shutdown")
		})
		f.Group("/reverse", func() {
			f.Post("/start")
			f.Post("/shutdown")
		})
		//f.Post("/start", binding.JSON(form.StartProxy{}), proxy.Start)
		//f.Post("/shut-down", proxy.ShutDown)
	})

	modules := route.NewModulesHandler()
	f.Group("/modules", func() {
		f.Combo("").Get().Post()
		f.Post("/reload", modules.RefreshModule)
		//f.Combo("/new").Get(modules.New).Post(binding.JSON(form.NewModule{}), modules.NewAction)
		f.Group("/{id}", func() {
			f.Combo("").
				Get(modules.Get).
				Put().
				Delete()
			f.Post("/enable", modules.SetStatus(route.Enable))
			f.Post("/disable", modules.SetStatus(route.Disable))
		})
		//f.Post("/{id}/enable", modules.SetStatus(route.Enable))
		//f.Post("/{id}/disable", modules.SetStatus(route.Disable))
		//f.Get("/{id}", modules.Get)
		//f.Put("/{id}", binding.JSON(form.UpdateModule{}), modules.Update)
		//f.Delete("/{id}", modules.Delete)
	})

	httpPort := c.Int("port")
	f.Run("0.0.0.0", httpPort)

	return nil
}
