package module

import (
	"net/http"
	"os"

	log "unknwon.dev/clog/v2"
)

func Init() {
	_, err := os.Stat("modules")
	if os.IsNotExist(err) {
		err := os.Mkdir("modules", 0755)
		if err != nil {
			log.Fatal("Failed to create modules folder: %v", err)
		}
		log.Info("Create modules folder.")
	}
}

var modules []*module

// Load loads the module files from modules folder.
func Load() []*module {
	module, err := NewModule("./modules/demo.yml")
	if err != nil {
		log.Fatal("Load modules: %v", err)
	}
	modules = append(modules, module)
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
