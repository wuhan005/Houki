// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"context"
	"net/http"

	"github.com/wuhan005/Houki/internal/modules"
)

var enableModules = make(moduleList)

type moduleList map[string]*modules.Body

// LoadModule loads module with the given ID.
func LoadModule(ctx context.Context, id string) error {
	//mod, err := db.Modules.Get(ctx, id)
	//if err != nil {
	//	return errors.Wrap(err, "get module")
	//}
	//
	//enableModules[mod.ID] = mod.Body
	//return nil
	return nil
}

// ReloadAllModules loads all modules.
func ReloadAllModules(ctx context.Context) error {
	//modules, err := db.Modules.List(ctx, db.GetModuleOptions{EnabledOnly: true})
	//if err != nil {
	//	return errors.Wrap(err, "load from database")
	//}
	//
	//enableModules = make(moduleList)
	//
	//for _, mod := range modules {
	//	body := mod.Body
	//
	//	if err := body.Init(); err != nil {
	//		log.Error("Failed to load module %q: %v", mod.ID, err)
	//		continue
	//	}
	//
	//	enableModules[mod.ID] = body
	//}
	//return nil
	return nil
}

func DoRequest(req *http.Request, body []byte) {
	for _, mod := range enableModules {
		mod.DoRequest(req, body)
	}
}

func DoResponse(resp *http.Response, body []byte) {
	for _, mod := range enableModules {
		mod.DoResponse(resp, body)
	}
}
