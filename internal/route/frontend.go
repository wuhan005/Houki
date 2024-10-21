package route

import (
	"net/http"

	"github.com/wuhan005/Houki/internal/context"
	"github.com/wuhan005/Houki/web"
)

func Frontend(c context.Context) {
	if c.Request().Method != http.MethodGet && c.Request().Method != http.MethodHead {
		return
	}

	name := "index.html"
	f, err := http.FS(web.Embed).Open(name)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	fi, err := f.Stat()
	if err != nil {
		return // File exists but failed to open.
	}

	http.ServeContent(c.ResponseWriter(), c.Request().Request, name, fi.ModTime(), f)
}
