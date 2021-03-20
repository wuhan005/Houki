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

			// Response
			decls.NewVar("status_code", decls.Int),

			decls.NewVar("headers", decls.Dyn),
			decls.NewVar("body", decls.String),
		),
	)
}
