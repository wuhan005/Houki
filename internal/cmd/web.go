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
	"github.com/wuhan005/Houki/internal/form"
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

	f.Group("/api", func() {

		proxy := route.NewProxyHandler()
		f.Group("/proxy", func() {
			f.Get("/status", proxy.Status)

			f.Group("/forward", func() {
				f.Post("/start", form.Bind(form.StartProxy{}), proxy.StartForward)
				f.Post("/shutdown", proxy.ShutdownForward)
			})
			f.Group("/reverse", func() {
				f.Post("/start", form.Bind(form.StartProxy{}), proxy.StartReverse)
				f.Post("/shutdown", proxy.ShutdownReverse)
			})
		})

		modules := route.NewModulesHandler()
		f.Group("/modules", func() {
			f.Combo("").Get(modules.List).Post(form.Bind(form.CreateModule{}), modules.Create)
			f.Group("/{id}", func() {
				f.Combo("").
					Get(modules.Get).
					Put(form.Bind(form.UpdateModule{}), modules.Update).
					Delete(modules.Delete)
				f.Post("/enable", modules.SetStatus(route.Enable))
				f.Post("/disable", modules.SetStatus(route.Disable))
			}, modules.Moduler)

			f.Post("/reload", modules.ReloadModules)
		})

		certificate := route.NewCertificateHandler()
		f.Combo("/certificate").
			Get(certificate.Get).
			Put(form.Bind(form.UpdateCertificate{}), certificate.Update)
	})

	httpPort := c.Int("port")
	f.Run("0.0.0.0", httpPort)

	return nil
}
