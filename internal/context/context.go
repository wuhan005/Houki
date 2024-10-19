// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package context

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flamego/flamego"
	"github.com/sirupsen/logrus"
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

func (c *Context) Success(data ...interface{}) error {
	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(http.StatusOK)

	var d interface{}
	if len(data) == 1 {
		d = data[0]
	} else {
		d = ""
	}

	if err := json.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"msg":  "success",
			"data": d,
		},
	); err != nil {
		logrus.WithContext(c.Request().Context()).WithError(err).Error("Failed to encode JSON response")
	}
	return nil
}

func (c *Context) ServerError() error {
	return c.Error(http.StatusInternalServerError, "Internal server error")
}

func (c *Context) Error(statusCode int, message string, v ...interface{}) error {
	c.ResponseWriter().Header().Set("Content-Type", "application/json")
	c.ResponseWriter().WriteHeader(statusCode)

	if len(v) != 0 {
		message = fmt.Sprintf(message, v...)
	}

	if err := json.NewEncoder(c.ResponseWriter()).Encode(
		map[string]interface{}{
			"msg": message,
		},
	); err != nil {
		logrus.WithContext(c.Request().Context()).WithError(err).Error("Failed to encode JSON response")
	}
	return nil
}
