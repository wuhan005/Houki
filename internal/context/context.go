// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"net/http"

	"github.com/flamego/flamego"
)

// Context represents context of a request.
type Context struct {
	flamego.Context
}

func (c *Context) ServerError() {
	c.ResponseWriter().WriteHeader(http.StatusInternalServerError)
}

// Contexter initializes a classic context for a request.
func Contexter() flamego.Handler {
	return func(ctx flamego.Context) {
		c := Context{
			Context: ctx,
		}

		c.Map(c)
	}
}
