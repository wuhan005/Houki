package web

import (
	"net/http"

	"gopkg.in/macaron.v1"
	log "unknwon.dev/clog/v2"
)

type web struct {
	server *http.Server
}

func New() *web {
	m := macaron.Classic()
	
	return &web{
		server: &http.Server{
			Handler: m,
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
