package module

import (
	"net/url"
	"os"
	"path/filepath"

	"github.com/google/cel-go/cel"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/wuhan005/Houki/internal/expression"
)

type Request struct {
	On    string      `yaml:"on" json:"on"`
	OnPrg cel.Program `json:"-"`

	Transmit    string                 `yaml:"transmit" json:"transmit"`
	TransmitURL *url.URL               `json:"-"`
	Header      map[string]string      `yaml:"header" json:"header"`
	Body        map[string]interface{} `yaml:"body" json:"-"` //FIXME
}

type Response struct {
	On    string      `yaml:"on" json:"on"`
	OnPrg cel.Program `json:"-"`

	StatusCode int                    `yaml:"status_code" json:"status_code"`
	Header     map[string]string      `yaml:"header" json:"header"`
	Body       map[string]interface{} `yaml:"body" json:"-"` //FIXME
}

type module struct {
	Env *cel.Env `json:"-"`

	FileName string `yaml:"-" json:"file_name"`

	Title       string `yaml:"title" json:"title"`
	Author      string `yaml:"author" json:"author"`
	Description string `yaml:"description" json:"description"`
	ID          string `yaml:"id" json:"id"`
	Sign        string `yaml:"sign" json:"sign"`

	Req  *Request  `yaml:"request" json:"request"`
	Resp *Response `yaml:"response" json:"response"`
}

func NewModule(filePath string) (*module, error) {
	mod, err := ParseFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "parse module file")
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

	return mod, nil
}

func ParseFile(filePath string) (*module, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "read file")
	}

	var mod module
	err = yaml.Unmarshal(raw, &mod)
	if err != nil {
		return nil, errors.Wrap(err, "parse module yaml")
	}

	mod.FileName = filepath.Base(filePath)
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
