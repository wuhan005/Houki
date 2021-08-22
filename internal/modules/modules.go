// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/Houki/internal/db"
)

type ModuleList map[string]*Module

func (modules ModuleList) Get(id string) (*Module, error) {
	module, ok := modules[id]
	if ok {
		return module, nil
	}
	return nil, errors.New("module not found")
}

var modulesMap = make(ModuleList)

// Reload loads all the modules.
func Reload(ctx context.Context) (map[string]*Module, error) {
	modules, err := db.Modules.List(ctx, db.GetModuleOptions{EnabledOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "load from database")
	}

	modulesMap = make(map[string]*Module)

	for _, mod := range modules {
		module, err := NewModule(mod.FilePath, mod.Enabled)
		if err != nil {
			log.Error("Failed to load module %q: %v", mod.ID, err)
			continue
		}

		modulesMap[mod.ID] = module
	}
	return modulesMap, nil
}

func List() ModuleList {
	return modulesMap
}

func Get(id string) (*Module, error) {
	return modulesMap.Get(id)
}

func DoRequest(req *http.Request, body []byte) {
	for _, module := range modulesMap {
		if module.Enabled {
			module.DoRequest(req, body)
		}
	}
}

func DoResponse(resp *http.Response, body []byte) {
	for _, module := range modulesMap {
		if module.Enabled {
			module.DoResponse(resp, body)
		}
	}
}
