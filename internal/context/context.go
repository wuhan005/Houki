// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"fmt"
	"net/http"

	"github.com/flamego/flamego"
	jsoniter "github.com/json-iterator/go"
	log "unknwon.dev/clog/v2"
)

// Context represents context of a request.
type Context struct {
	flamego.Context
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

func (c *Context) Success(data ...interface{}) {
	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)

	var d interface{}
	if len(data) == 1 {
		d = data[0]
	} else {
		d = ""
	}

	err := jsoniter.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"error": 0,
			"data":  d,
		},
	)
	if err != nil {
		log.Error("Failed to encode: %v", err)
	}
}

func (c *Context) ServerError() {
	c.Error(http.StatusInternalServerError*100, "Internal server error")
}

func (c *Context) Error(errorCode uint, message string, v ...interface{}) {
	statusCode := int(errorCode / 100)

	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(statusCode)

	if len(v) != 0 {
		message = fmt.Sprintf(message, v...)
	}

	err := jsoniter.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"error": errorCode,
			"msg":   message,
		},
	)
	if err != nil {
		log.Error("Failed to encode: %v", err)
	}
}
