package module

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/wuhan005/Houki/internal/conf"

	log "unknwon.dev/clog/v2"
)

func Initialize() error {
	_, err := os.Stat("modules")
	if os.IsNotExist(err) {
		err := os.Mkdir("modules", 0755)
		if err != nil {
			return errors.Wrap(err, "create modules folder")
		}
	}
	return nil
}

var modules []*module

// Load loads the module files from modules folder.
func Load() []*module {
	moduleList := conf.Get().Modules
	for _, modName := range moduleList {
		mod, err := NewModule(filepath.Join("./modules/", modName+".yml"))
		if err != nil {
			log.Error("Failed to load module %q: %v", modName, err)
			continue
		}

		modules = append(modules, mod)
	}

	log.Trace("Load %d modules.", len(modules))
	return modules
}

func List() []*module {
	return modules
}

func DoRequest(req *http.Request, body []byte) {
	for _, mod := range modules {
		mod.DoRequest(req, body)
	}
}

func DoResponse(resp *http.Response, body []byte) {
	for _, mod := range modules {
		mod.DoResponse(resp, body)
	}
}
