// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package expression

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/checker/decls"
)

func NewEnv() (*cel.Env, error) {
	return cel.NewEnv(
		cel.Declarations(
			// Request
			decls.NewVar("method", decls.String),
			decls.NewVar("url", decls.String),
			decls.NewVar("host", decls.String),

			// Response
			decls.NewVar("status_code", decls.Int),

			decls.NewVar("headers", decls.Dyn),
			decls.NewVar("body", decls.String),
		),
	)
}
