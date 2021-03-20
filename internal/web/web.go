package web

import (
	"io/fs"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"github.com/wuhan005/Houki/frontend"
	log "unknwon.dev/clog/v2"
)

type web struct {
	server *http.Server
}

func New() *web {
	r := gin.Default()

	store := cookie.NewStore([]byte(randstr.String(50)))
	r.Use(sessions.Sessions("Houki", store))

	fe, err := fs.Sub(frontend.FS, "dist")
	if err != nil {
		log.Fatal("Failed to sub path `dist`: %v", err)
	}
	r.StaticFS("/", http.FS(fe))

	return &web{
		server: &http.Server{
			Handler: r,
		},
	}
}

func (w *web) Run(addr string) {
	w.server.Addr = addr

	go func() {
		if err := w.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("Failed to start web server: %v", err)
		}
	}()
	log.Info("Web server listening on %s", addr)
}

func (w *web) Stop() error {
	return w.server.Shutdown(nil)
}
