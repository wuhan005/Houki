// Copyright 2024 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package form

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/flamego/flamego"
	"github.com/wuhan005/govalid"

	"github.com/wuhan005/Houki/internal/context"
)

func Bind(model interface{}) flamego.Handler {
	if reflect.TypeOf(model).Kind() == reflect.Ptr {
		panic("form: pointer can not be accepted as binding model")
	}

	return func(ctx context.Context) error {
		obj := reflect.New(reflect.TypeOf(model))
		r := ctx.Request().Request
		if r.Body != nil {
			defer func() { _ = r.Body.Close() }()
			if err := json.NewDecoder(r.Body).Decode(obj.Interface()); err != nil {
				return ctx.Error(http.StatusBadRequest, "Failed to parse request body")
			}
		}

		errors, ok := govalid.Check(obj.Interface())
		if !ok {
			return ctx.Error(http.StatusBadRequest, errors[0].Error())
		}

		// Validation passed.
		ctx.Map(obj.Elem().Interface())
		return nil
	}
}
