// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/flamego/binding"
	"github.com/flamego/flamego"
	"github.com/flamego/session"
	"github.com/flamego/template"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

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
	f.Use(session.Sessioner())

	f.Get("/", func(ctx context.Context) {
		ctx.Redirect("/proxy/")
	})

	proxy := route.NewProxyHandler()
	f.Group("/proxy", func() {
		f.Get("/", proxy.Dashboard)
		f.Post("/start", binding.Form(form.StartProxy{}), proxy.Start)
		f.Post("/shut_down", proxy.ShutDown)
	})

	modules := route.NewModulesHandler()
	f.Group("/modules", func() {
		f.Get("/", modules.List)
		f.Post("/", modules.New)
		f.Get("/{id}", modules.Get)
		f.Post("/{id}/enable", modules.SetStatus(route.Enable))
		f.Post("/{id}/disable", modules.SetStatus(route.Disable))
		f.Put("/{id}")
		f.Delete("/{id}")
	})

	//r := gin.Default()
	//// TODO remove CORS headers
	//r.Use(cors.New(cors.Config{
	//	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
	//	AllowHeaders:     []string{"Content-type", "User-Agent"},
	//	AllowCredentials: true,
	//	AllowOrigins:     []string{"http://localhost:8080"},
	//}))
	//
	//store := cookie.NewStore([]byte(randstr.String(50)))
	//r.Use(sessions.Sessions("Houki", store))
	//
	//api := r.Group("/api")
	//
	//// Proxy
	//pxy := api.Group("/proxy")
	//pxy.GET("/status", __(proxy.GetStatus))
	//pxy.POST("/start", __(proxy.Start))
	//pxy.POST("/stop", __(proxy.Stop))
	//// Proxy CA
	//pxy.GET("/ca", __(proxy.FetchCA))
	//pxy.POST("/ca/generate", __(proxy.GenerateCA))
	//
	//// Modules
	//api.GET("/modules", __(module.ListModules))
	//api.POST("/module/enable/:id", __(module.EnableModule))
	//api.POST("/module/disable/:id", __(module.DisableModule))
	//
	//// Frontend static assets
	//fe, err := fs.Sub(frontend.FS, "dist")
	//if err != nil {
	//	log.Fatal("Failed to sub path `dist`: %v", err)
	//}
	//r.StaticFS("/m", http.FS(fe))
	//
	httpPort := c.Int("port")
	f.Run("0.0.0.0", httpPort)
	return nil
}
