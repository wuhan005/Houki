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

var enabledModules map[string]*Module

// Reload loads the enabled modules.
func Reload(ctx context.Context) (map[string]*Module, error) {
	modules, err := db.Modules.List(ctx, db.GetModuleOptions{EnabledOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "load from database")
	}

	enabledModules = make(map[string]*Module)

	for _, mod := range modules {
		module, err := NewModule(mod.FilePath)
		if err != nil {
			log.Error("Failed to load module %q: %v", mod.ID, err)
			continue
		}

		enabledModules[mod.ID] = module
	}
	return enabledModules, nil
}

func DoRequest(req *http.Request, body []byte) {
	for _, module := range enabledModules {
		module.DoRequest(req, body)
	}
}

func DoResponse(resp *http.Response, body []byte) {
	for _, module := range enabledModules {
		module.DoResponse(resp, body)
	}
}
