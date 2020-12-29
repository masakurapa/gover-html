package filter

import (
	"strings"

	"github.com/masakurapa/gover-html/internal/option"
)

// Filter is filter the output directory
type Filter interface {
	IsOutputTarget(string, string) bool
}

type filter struct {
	include []string
	exclude []string
}

// New is initialize the filter
func New(opt option.Option) Filter {
	return &filter{
		include: opt.Include,
		exclude: opt.Exclude,
	}
}

// IsOutputTarget returns true if output target
// The "relativePath" must be relative to the base path
func (f *filter) IsOutputTarget(relativePath, fileName string) bool {
	// absolute path is always NG
	if strings.HasPrefix(relativePath, "/") {
		return false
	}

	path := f.convertPathForValidate(relativePath)

	for _, s := range f.exclude {
		if f.hasPrefix(path, fileName, s) {
			return false
		}
	}

	if len(f.include) == 0 {
		return true
	}

	for _, s := range f.include {
		if f.hasPrefix(path, fileName, s) {
			return true
		}
	}
	return false
}

func (f *filter) hasPrefix(path, fileName, prefix string) bool {
	if path == prefix || strings.HasPrefix(path, prefix+"/") {
		return true
	}
	return path+"/"+fileName == prefix
}

func (f *filter) convertPathForValidate(relativePath string) string {
	path := strings.TrimPrefix(relativePath, "./")
	return strings.TrimSuffix(path, "/")
}
