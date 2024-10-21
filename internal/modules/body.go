// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package modules

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net/url"

	"github.com/google/cel-go/cel"
	"github.com/pkg/errors"

	"github.com/wuhan005/Houki/internal/expression"
)

// Body is the module body.
type Body struct {
	Req  *Request  `json:"request"`
	Resp *Response `json:"response"`

	Env *cel.Env `json:"-"`
}

var _ sql.Scanner = (*Body)(nil)

func (b *Body) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal([]byte(value.(string)), &b)
}

var _ driver.Valuer = (*Body)(nil)

func (b *Body) Value() (driver.Value, error) {
	body, err := json.Marshal(b)
	if err != nil {
		return nil, err
	}
	return string(body), nil
}

func (b *Body) Init() error {
	if b.Req == nil {
		b.Req = &Request{}
	}
	if b.Resp == nil {
		b.Resp = &Response{}
	}

	// If the `on` condition is empty, it means the module always enabled.
	if b.Req.On == "" {
		b.Req.On = "true"
	}
	if b.Resp.On == "" {
		b.Resp.On = "true"
	}

	if b.Req.Transmit != "" {
		transmitURL, err := url.Parse(b.Req.Transmit)
		if err != nil {
			return errors.Wrap(err, "parse transmit url")
		}
		b.Req.TransmitURL = transmitURL
	}

	env, err := expression.NewEnv()
	if err != nil {
		return errors.Wrap(err, "new env")
	}
	b.Env = env

	// Parse `on` expression.
	b.Req.OnPrg, err = b.parseExpression(b.Req.On)
	if err != nil {
		return errors.Wrap(err, "parse request `on`")
	}
	b.Resp.OnPrg, err = b.parseExpression(b.Resp.On)
	if err != nil {
		return errors.Wrap(err, "parse response `on`")
	}

	return nil
}

type Request struct {
	On    string      `json:"on"`
	OnPrg cel.Program `json:"-"`

	Transmit    string                 `json:"transmit"`
	TransmitURL *url.URL               `json:"-"`
	Headers     map[string]string      `json:"headers"`
	Body        map[string]interface{} `json:"-"` //FIXME
}

type Response struct {
	On    string      `json:"on"`
	OnPrg cel.Program `json:"-"`

	StatusCode int                    `json:"statusCode"`
	Headers    map[string]string      `json:"headers"`
	Body       map[string]interface{} `json:"body"` //FIXME
}

// parseExpression parses the module expression in itself environment.
func (b *Body) parseExpression(expression string) (cel.Program, error) {
	ast, issues := b.Env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return nil, errors.Wrap(issues.Err(), "type check")
	}

	prg, err := b.Env.Program(ast)
	if err != nil {
		return nil, errors.Wrap(err, "program construction")
	}
	return prg, nil
}
