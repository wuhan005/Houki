package module

import (
	"net/url"
	"os"

	"github.com/google/cel-go/cel"
	"github.com/pkg/errors"
	"github.com/wuhan005/Houki/expression"
	"gopkg.in/yaml.v2"
)

type Request struct {
	On    string `yaml:"on"`
	OnPrg cel.Program

	Transmit    string `yaml:"transmit"`
	TransmitURL *url.URL
	Header      map[string]string      `yaml:"header"`
	Body        map[string]interface{} `yaml:"body"`
}

type Response struct {
	On    string `yaml:"on"`
	OnPrg cel.Program

	StatusCode int                    `yaml:"status_code"`
	Header     map[string]string      `yaml:"header"`
	Body       map[string]interface{} `yaml:"body"`
}

type module struct {
	Env *cel.Env

	Title       string `yaml:"title"`
	Author      string `yaml:"author"`
	Description string `yaml:"description"`
	ID          string `yaml:"id"`
	Sign        string `yaml:"sign"`

	Req  *Request  `yaml:"request"`
	Resp *Response `yaml:"response"`
}

func NewModule(filePath string) (*module, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	var mod module
	err = yaml.Unmarshal(raw, &mod)
	if err != nil {
		return nil, errors.Wrap(err, "parse module yaml")
	}

	if mod.Req.On == "" {
		mod.Req.On = "true"
	}

	if mod.Resp.On == "" {
		mod.Resp.On = "true"
	}

	if mod.Req.Transmit != "" {
		transmitURL, err := url.Parse(mod.Req.Transmit)
		if err != nil {
			return nil, errors.Wrap(err, "parse transmit url")
		}
		mod.Req.TransmitURL = transmitURL
	}

	env, err := expression.NewEnv()
	if err != nil {
		return nil, errors.Wrap(err, "new env")
	}
	mod.Env = env

	// Parse `on` program.
	mod.Req.OnPrg, err = mod.parse(mod.Req.On)
	if err != nil {
		return nil, errors.Wrap(err, "parse request `on`")
	}
	mod.Resp.OnPrg, err = mod.parse(mod.Resp.On)
	if err != nil {
		return nil, errors.Wrap(err, "parse response `on`")
	}

	return &mod, nil
}

func (m *module) parse(expression string) (cel.Program, error) {
	ast, issues := m.Env.Compile(expression)
	if issues != nil && issues.Err() != nil {
		return nil, errors.Wrap(issues.Err(), "type check")
	}
	prg, err := m.Env.Program(ast)
	if err != nil {
		return nil, errors.Wrap(err, "program construction")
	}
	return prg, err
}
