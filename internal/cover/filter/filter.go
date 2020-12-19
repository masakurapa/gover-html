package filter

import "strings"

const filterSeparator = ","

// Filter is filter the output directory
type Filter interface {
	IsOutputTarget(string) bool
}

type filter struct {
	include []string
}

// New is initialize the filter
func New(include *string) Filter {
	return &filter{
		include: parse(include),
	}
}

func parse(value *string) []string {
	if value == nil || *value == "" {
		return []string{}
	}
	return convert(strings.Split(*value, filterSeparator))
}

func convert(values []string) []string {
	newFilter := make([]string, 0, len(values))
	for _, f := range values {
		s := strings.TrimSpace(f)
		s = strings.TrimPrefix(s, "./")

		if !strings.HasSuffix(s, "/") {
			s += "/"
		}
		newFilter = append(newFilter, s)
	}
	return newFilter
}

// IsOutputTarget returns true if output target
// The "relativePath" must be relative to the base path
func (f *filter) IsOutputTarget(relativePath string) bool {
	// absolute path is always NG
	if strings.HasPrefix(relativePath, "/") {
		return false
	}
	if len(f.include) == 0 {
		return true
	}

	path := f.convertPathForValidate(relativePath)
	for _, f := range f.include {
		if strings.HasPrefix(path, f) {
			return true
		}
	}
	return false
}

func (f *filter) convertPathForValidate(relativePath string) string {
	path := strings.TrimPrefix(relativePath, "./")
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}
