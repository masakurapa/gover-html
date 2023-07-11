package option

import (
	"fmt"
	"io"
	"regexp"
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

var (
	// 関数除外設定の検証用の正規表現（ここでは緩い検証にする）
	excludeFuncFormat = regexp.MustCompile(`^\((\./)?([a-zA-Z\d/\-_]+)(\.[a-zA-Z].+)?\)\.([a-zA-Z].+)$`)
)

type optionConfig struct {
	Input       string
	Output      string
	Theme       string
	Include     []string // include fire or directories
	Exclude     []string // exclude fire or directories
	ExcludeFunc []string // exclude functions
}

type Option struct {
	Input       string
	Output      string
	Theme       string
	Include     []string            // include fire or directories
	Exclude     []string            // exclude fire or directories
	ExcludeFunc []ExcludeFuncOption // exclude functions
}

type ExcludeFuncOption struct {
	Package string
	Struct  string
	Func    string
}

type Generator struct {
	r reader.Reader
}

func New(r reader.Reader) *Generator {
	return &Generator{r: r}
}

func (g *Generator) Generate(
	input *string,
	output *string,
	theme *string,
	include *string,
	exclude *string,
	excludeFunc *string,
) (*Option, error) {
	opt := &optionConfig{}
	if g.r.Exists(fileName) {
		fileOpt, err := g.readOptionFile()
		if err != nil {
			return nil, err
		}
		opt = fileOpt
	}

	opt.Input = g.stringValue(input, opt.Input)
	opt.Output = g.stringValue(output, opt.Output)
	opt.Theme = g.stringValue(theme, opt.Theme)
	opt.Include = g.stringsValue(include, opt.Include)
	opt.Exclude = g.stringsValue(exclude, opt.Exclude)
	opt.ExcludeFunc = g.stringsValue(excludeFunc, opt.ExcludeFunc)
	return g.getValidatedOption(opt)
}

func (g *Generator) readOptionFile() (*optionConfig, error) {
	r, err := g.r.Read(fileName)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	opt := optionConfig{}
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

func (g *Generator) getValidatedOption(opt *optionConfig) (*Option, error) {
	if err := g.validate(opt); err != nil {
		return nil, err
	}
	return g.getOptionWithDefaultValue(opt), nil
}

func (g *Generator) validate(opt *optionConfig) error {
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
	if es := g.validateExcludeFunc(opt.ExcludeFunc); len(es) > 0 {
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

func (g *Generator) validateExcludeFunc(values []string) optionErrors {
	errs := make(optionErrors, 0)
	for _, v := range values {
		if g.isEmpty(v) {
			continue
		}

		if strings.HasPrefix(v, "(/") {
			errs = append(errs, fmt.Errorf("exclude-func value (%q) must not be an absolute path", v))
			continue
		}

		// ()が含まれない場合は関数名のみとみなしてOK
		if !strings.Contains(v, "(") && !strings.Contains(v, ")") {
			continue
		}

		if !excludeFuncFormat.MatchString(v) {
			errs = append(errs, fmt.Errorf("exclude-func value (%q) format is invalid", v))
		}
	}
	return errs
}

func (g *Generator) getOptionWithDefaultValue(opt *optionConfig) *Option {
	newOpt := &Option{
		Input:   opt.Input,
		Output:  opt.Output,
		Theme:   opt.Theme,
		Include: opt.Include,
		Exclude: opt.Exclude,
	}

	if g.isEmpty(newOpt.Input) {
		newOpt.Input = inputDefault
	}
	if g.isEmpty(newOpt.Output) {
		newOpt.Output = outputDefault
	}
	if g.isEmpty(newOpt.Theme) {
		newOpt.Theme = themeDefault
	}

	newOpt.Include = g.convertFilterValue(newOpt.Include)
	newOpt.Exclude = g.convertFilterValue(newOpt.Exclude)
	newOpt.ExcludeFunc = g.convertExcludeFuncOption(opt.ExcludeFunc)
	return newOpt
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

func (g *Generator) convertExcludeFuncOption(values []string) []ExcludeFuncOption {
	ret := make([]ExcludeFuncOption, 0, len(values))
	for _, v := range values {
		s := strings.TrimSpace(v)
		if g.isEmpty(s) {
			continue
		}

		// ()が含まれない場合は関数名のみとみなして終了
		if !strings.Contains(s, "(") && !strings.Contains(s, ")") {
			ret = append(ret, ExcludeFuncOption{Func: s})
			continue
		}

		matches := excludeFuncFormat.FindStringSubmatch(s)
		structName := matches[3]
		if structName != "" {
			structName = structName[1:]
		}

		ret = append(ret, ExcludeFuncOption{
			Package: strings.TrimSuffix(strings.TrimPrefix(matches[2], "./"), "/"),
			Struct:  structName,
			Func:    matches[4],
		})
	}
	return ret
}
