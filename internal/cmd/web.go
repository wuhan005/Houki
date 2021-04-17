// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/frontend"
	"github.com/wuhan005/Houki/internal/route/module"
	"github.com/wuhan005/Houki/internal/route/proxy"
	"github.com/wuhan005/Houki/internal/sse"
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
	r := gin.Default()
	// TODO remove CORS headers
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Content-type", "User-Agent"},
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:8080"},
	}))

	store := cookie.NewStore([]byte(randstr.String(50)))
	r.Use(sessions.Sessions("Houki", store))

	sse.Initialize()
	api := r.Group("/api")
	api.GET("/logs", proxy.LogHandler)

	// Proxy
	pxy := api.Group("/proxy")
	pxy.GET("/status", __(proxy.GetStatus))
	pxy.POST("/start", __(proxy.Start))
	pxy.POST("/stop", __(proxy.Stop))
	// Proxy CA
	pxy.POST("/ca/generate")
	pxy.GET("/ca/download")

	// Modules
	api.GET("/modules", __(module.ListModules))
	api.POST("/module/enable/:id", __(module.EnableModule))
	api.POST("/module/disable/:id", __(module.DisableModule))

	// Frontend static assets
	fe, err := fs.Sub(frontend.FS, "dist")
	if err != nil {
		log.Fatal("Failed to sub path `dist`: %v", err)
	}
	r.StaticFS("/m", http.FS(fe))

	httpPort := c.String("port")
	return r.Run(fmt.Sprintf("%s:%s", "0.0.0.0", httpPort))
}

func __(handler func(*gin.Context) (int, interface{})) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(handler(c))
	}
}
