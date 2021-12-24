// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"net/http"

	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"github.com/flamego/session"
	"github.com/flamego/template"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/assets"
	"github.com/wuhan005/Houki/internal/context"
	"github.com/wuhan005/Houki/internal/db"
	"github.com/wuhan005/Houki/internal/form"
	"github.com/wuhan005/Houki/internal/proxy"
	"github.com/wuhan005/Houki/internal/route"
	"github.com/wuhan005/Houki/templates"
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
	_, err := proxy.SetDefaultProxy()
	if err != nil {
		return errors.Wrap(err, "set default proxy")
	}

	_, err = db.New()
	if err != nil {
		return errors.Wrap(err, "new db")
	}

	f := flamego.Classic()
	f.Use(context.Contexter())

	fs, err := template.EmbedFS(templates.FS, ".", []string{".tmpl"})
	if err != nil {
		log.Fatal("Failed to embed template file system: %v", err)
	}
	f.Use(template.Templater(template.Options{
		FileSystem: fs,
	}))

	f.Use(flamego.Static(flamego.StaticOptions{
		FileSystem: http.FS(assets.StaticFS),
	}))
	f.Use(session.Sessioner())

	f.Get("/", func(ctx context.Context) {
		ctx.Redirect("/proxy/")
	})

	proxy := route.NewProxyHandler()
	f.Group("/proxy", func() {
		f.Get("/", proxy.Dashboard)
		f.Post("/start", binding.JSON(form.StartProxy{}), proxy.Start)
		f.Post("/shut-down", proxy.ShutDown)
	})

	modules := route.NewModulesHandler()
	f.Group("/modules", func() {
		f.Post("/reload", modules.RefreshModule)
		f.Combo("/new").Get(modules.New).Post(binding.JSON(form.NewModule{}), modules.NewAction)
		f.Post("/{id}/enable", modules.SetStatus(route.Enable))
		f.Post("/{id}/disable", modules.SetStatus(route.Disable))
		f.Get("/{id}", modules.Get)
		f.Put("/{id}", binding.JSON(form.UpdateModule{}), modules.Update)
		f.Delete("/{id}", modules.Delete)
	})

	chromeDp := route.NewChromeDpHandler()
	f.Group("/chromedp", func() {
		f.Post("/", chromeDp.New)
	})

	httpPort := c.Int("port")
	f.Run("0.0.0.0", httpPort)
	return nil
}
