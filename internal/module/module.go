package module

import (
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/db"
)

// Initialize creates the `modules` folder if not existed.
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

var enabledModules []*module

// Reload loads the module files from `modules` folder into proxy.
func Reload() ([]*module, error) {
	modules, err := Scan()
	if err != nil {
		return nil, err
	}

	enabledModIDs, err := GetEnabledModules()
	if err != nil {
		return nil, errors.Wrap(err, "get enabled modules")
	}

	// Clean the enabled modules.
	enabledModules = make([]*module, 0)
	for _, enabledModID := range enabledModIDs {
		ok := false
		for _, mod := range modules {
			if enabledModID == mod.ID {
				modInstance, err := NewModule(filepath.Join("./modules/", mod.FileName))
				if err != nil {
					_ = db.Modules.Disable(enabledModID)
					log.Error("Failed to load module %q: %v", enabledModID, err)
					break
				}
				ok = true
				enabledModules = append(enabledModules, modInstance)
			}
		}
		if !ok {
			_ = db.Modules.Disable(enabledModID)
			log.Error("Failed to load module, not found: %v", enabledModID)
		}
	}
	return enabledModules, nil
}

// Enable enables a module by given module ID and reload the enabled modules.
func Enable(moduleID string) error {
	if err := db.Modules.Enable(moduleID); err != nil {
		return err
	}
	_, err := Reload()
	return err
}

// Disable disables a module by given module ID and reload the enabled modules.
func Disable(moduleID string) error {
	if err := db.Modules.Disable(moduleID); err != nil {
		return err
	}
	_, err := Reload()
	return err
}

// Scan parse the YAML files as modules from `modules` folder.
func Scan() ([]*module, error) {
	files, err := os.ReadDir("./modules/")
	if err != nil {
		return nil, errors.Wrap(err, "read dir")
	}

	var modules []*module
	moduleIDs := make(map[string]struct{})
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if path.Ext(file.Name()) == ".yml" {
			mod, err := ParseFile(path.Join("./modules/", file.Name()))
			if err != nil {
				continue
			}

			if _, ok := moduleIDs[mod.ID]; ok {
				log.Error("Module ID existed: %q", mod.ID)
				continue
			}
			moduleIDs[mod.ID] = struct{}{}
			modules = append(modules, mod)
		}
	}

	return modules, nil
}

// GetEnabledModules returns the enabled modules uid.
func GetEnabledModules() ([]string, error) {
	return db.Modules.GetEnabled()
}

func DoRequest(req *http.Request, body []byte) {
	for _, mod := range enabledModules {
		mod.DoRequest(req, body)
	}
}

func DoResponse(resp *http.Response, body []byte) {
	for _, mod := range enabledModules {
		mod.DoResponse(resp, body)
	}
}
