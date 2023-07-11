package helper

import (
	"strings"
	"testing"

	"github.com/masakurapa/gover-html/internal/option"
	"github.com/masakurapa/gover-html/internal/reader"
)

type mockReader struct {
	reader.Reader
}

func (m *mockReader) Exists(string) bool {
	return false
}

// GetOptionForDefault returns default option
func GetOptionForDefault(t *testing.T) option.Option {
	return makeOption(t, nil, nil, nil, nil, nil, nil)
}

// GetOptionForInclude returns option with include set
func GetOptionForInclude(t *testing.T, val []string) option.Option {
	return makeOption(t, nil, nil, nil, joinComma(val), nil, nil)
}

// GetOptionForExclude returns option with exclude set
func GetOptionForExclude(t *testing.T, val []string) option.Option {
	return makeOption(t, nil, nil, nil, nil, joinComma(val), nil)
}

// GetOptionForIncludeAndExclude returns option with include and exclude set
func GetOptionForIncludeAndExclude(t *testing.T, val []string) option.Option {
	return makeOption(t, nil, nil, nil, nil, joinComma(val), nil)
}

// GetOptionForExcludeFunc returns option with exclude set
func GetOptionForExcludeFunc(t *testing.T, val []string) option.Option {
	return makeOption(t, nil, nil, nil, nil, nil, joinComma(val))
}

func joinComma(val []string) *string {
	s := strings.Join(val, ",")
	return &s
}

func makeOption(
	t *testing.T,
	input *string,
	output *string,
	theme *string,
	include *string,
	exclude *string,
	excludeFunc *string,
) option.Option {
	opt, err := option.New(&mockReader{}).Generate(
		input,
		output,
		theme,
		include,
		exclude,
		excludeFunc,
	)
	if err != nil {
		t.Fatal(err)
	}
	return *opt
}
