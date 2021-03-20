package route

import (
	"io/fs"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/frontend"
	"github.com/wuhan005/Houki/internal/route/module"
	"github.com/wuhan005/Houki/internal/route/proxy"
)

type web struct {
	server *http.Server
}

func New() *web {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Content-type", "User-Agent"},
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:8081"},
	}))

	store := cookie.NewStore([]byte(randstr.String(50)))
	r.Use(sessions.Sessions("Houki", store))

	api := r.Group("/api")
	// Proxy
	pxy := api.Group("/proxy")
	pxy.GET("/status", __(proxy.GetStatus))
	pxy.POST("/start", __(proxy.Start))
	pxy.POST("/stop", __(proxy.Stop))

	// Modules
	api.GET("/modules", __(module.GetModules))

	fe, err := fs.Sub(frontend.FS, "dist")
	if err != nil {
		log.Fatal("Failed to sub path `dist`: %v", err)
	}
	r.StaticFS("/m", http.FS(fe))

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

func __(handler func(*gin.Context) (int, interface{})) func(*gin.Context) {
	return func(c *gin.Context) {
		c.JSON(handler(c))
	}
}
