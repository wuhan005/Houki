// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"net/http"

	csmap "github.com/mhmtszr/concurrent-swiss-map"
	"github.com/pkg/errors"
)

var enableModules = csmap.Create[uint, *Body]()

// LoadModule loads module with the given ID.
func LoadModule(id uint, body *Body) error {
	if err := body.Init(); err != nil {
		return errors.Wrap(err, "init module")
	}

	enableModules.Store(id, body)
	return nil
}

// UnloadModule unloads module with the given ID.
func UnloadModule(id uint) error {
	enableModules.Delete(id)
	return nil
}

// ReloadAllModules loads all modules.
func ReloadAllModules(modules map[uint]*Body) error {
	enableModules.Clear()
	for id, body := range modules {
		body := body
		if err := body.Init(); err != nil {
			return errors.Wrap(err, "init module")
		}

		enableModules.Store(id, body)
	}

	return nil
}

func DoRequest(req *http.Request, body []byte) {
	enableModules.Range(func(key uint, mod *Body) (stop bool) {
		mod.DoRequest(req, body)
		return false
	})
}

func DoResponse(resp *http.Response, body []byte) {
	enableModules.Range(func(key uint, mod *Body) (stop bool) {
		mod.DoResponse(resp, body)
		return false
	})
}
