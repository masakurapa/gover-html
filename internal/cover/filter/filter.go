package filter

import (
	"path/filepath"
	"strings"

	"github.com/masakurapa/gover-html/internal/option"
)

// Filter is filter the output directory
type Filter interface {
	IsOutputTarget(string, string) bool
	IsOutputTargetFunc(string, string, string) bool
}

type filter struct {
	opt option.Option
}

// New is initialize the filter
func New(opt option.Option) Filter {
	return &filter{opt: opt}
}

// IsOutputTarget returns true if output target
// The "relativePath" must be relative to the base path
func (f *filter) IsOutputTarget(relativePath, fileName string) bool {
	// absolute path is always NG
	if strings.HasPrefix(relativePath, "/") {
		return false
	}

	path := f.convertPathForValidate(relativePath)

	for _, s := range f.opt.Exclude {
		if f.hasPrefix(path, fileName, s) {
			return false
		}
	}

	if len(f.opt.Include) == 0 {
		return true
	}

	for _, s := range f.opt.Include {
		if f.hasPrefix(path, fileName, s) {
			return true
		}
	}
	return false
}

func (f *filter) IsOutputTargetFunc(relativePath, structName, funcName string) bool {
	path := f.convertPathForValidate(relativePath)
	// ファイル名が入っているので削除
	path = filepath.Dir(path)

	for _, opt := range f.opt.ExcludeFunc {
		if opt.Package != "" && opt.Package != path {
			continue
		}
		if opt.Struct != "" && opt.Struct != structName {
			continue
		}
		if opt.Func != funcName {
			continue
		}
		return false
	}
	return true
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
