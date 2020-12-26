package option

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/masakurapa/gover-html/internal/reader"
	"gopkg.in/yaml.v2"
)

const (
	fileName = ".gover.yml"

	optionSeparator = ","

	inputDefault  = "coverage.out"
	outputDefault = "coverage.html"

	themeDark    = "dark"
	themeLight   = "light"
	themeDefault = themeDark
)

// Option is option data
type Option struct {
	Input   string
	Output  string
	Theme   string
	Include []string
	Exclude []string
}

// Generator is generates option
type Generator struct {
	r reader.Reader
}

// New is initialize the option generator
func New(r reader.Reader) *Generator {
	return &Generator{r: r}
}

// Generate returns option
func (g *Generator) Generate(
	input *string,
	output *string,
	theme *string,
	include *string,
	exclude *string,
) (*Option, error) {
	opt := Option{}
	if g.r.Exists(fileName) {
		fileOpt, err := g.readOptionFile()
		if err != nil {
			return nil, err
		}
		opt = *fileOpt
	}

	opt.Input = g.stringValue(input, opt.Input)
	opt.Output = g.stringValue(output, opt.Output)
	opt.Theme = g.stringValue(theme, opt.Theme)
	opt.Include = g.stringsValue(include, opt.Include)
	opt.Exclude = g.stringsValue(exclude, opt.Exclude)
	return g.getValidatedOption(opt)
}

func (g *Generator) readOptionFile() (*Option, error) {
	r, err := g.r.Read(fileName)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	opt := Option{}
	if err := yaml.Unmarshal(b, &opt); err != nil {
		return nil, err
	}

	return &opt, nil
}

func (g *Generator) stringValue(arg *string, opt string) string {
	if arg == nil {
		return opt
	}
	return *arg
}

func (g *Generator) stringsValue(arg *string, opt []string) []string {
	if arg == nil {
		return opt
	}
	return strings.Split(*arg, optionSeparator)
}

func (g *Generator) getValidatedOption(opt Option) (*Option, error) {
	if err := g.validate(opt); err != nil {
		return nil, err
	}
	return g.getOptionWithDefaultValue(opt), nil
}

func (g *Generator) validate(opt Option) error {
	errs := make(optionErrors, 0)

	if !g.isEmpty(opt.Theme) && opt.Theme != themeDark && opt.Theme != themeLight {
		errs = append(errs, fmt.Errorf("theme must be %q or %q", themeDark, themeLight))
	}

	if es := g.validateFilter("include", opt.Include); len(es) > 0 {
		errs = append(errs, es...)
	}
	if es := g.validateFilter("exclude", opt.Exclude); len(es) > 0 {
		errs = append(errs, es...)
	}

	if len(errs) > 0 {
		return &errs
	}
	return nil
}

func (g *Generator) validateFilter(f string, values []string) optionErrors {
	errs := make(optionErrors, 0)
	for _, v := range values {
		if g.isEmpty(v) {
			continue
		}

		if strings.HasPrefix(v, "/") {
			errs = append(errs, fmt.Errorf("%s value (%q) must not be an absolute path", f, v))
		}
	}
	return errs
}

func (g *Generator) getOptionWithDefaultValue(opt Option) *Option {
	if g.isEmpty(opt.Input) {
		opt.Input = inputDefault
	}
	if g.isEmpty(opt.Output) {
		opt.Output = outputDefault
	}
	if g.isEmpty(opt.Theme) {
		opt.Theme = themeDefault
	}

	opt.Include = g.convertFilterValue(opt.Include)
	opt.Exclude = g.convertFilterValue(opt.Exclude)
	return &opt
}

func (g *Generator) isEmpty(s string) bool {
	return s == ""
}

func (g *Generator) convertFilterValue(values []string) []string {
	ret := make([]string, 0, len(values))
	for _, v := range values {
		s := strings.TrimSpace(v)
		if g.isEmpty(s) {
			continue
		}

		s = strings.TrimPrefix(s, "./")
		s = strings.TrimSuffix(s, "/")
		ret = append(ret, s)
	}
	return ret
}
