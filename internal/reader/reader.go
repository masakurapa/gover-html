package reader

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/masakurapa/go-cover/internal/logger"
)

type Reader interface {
	Read(path string) ([]byte, error)
}

type reader struct {
	pkg        *build.Package
	hasSrcRoot bool
	prefix     string
}

func New() Reader {
	return &reader{}
}

func (r *reader) Read(path string) ([]byte, error) {
	file, err := r.find(path)
	if err != nil {
		return nil, err
	}

	logger.L.Debug("read: %s", file)
	return ioutil.ReadFile(file)
}

func (r *reader) find(path string) (string, error) {
	dir, file := filepath.Split(path)

	if r.pkg == nil {
		pkg, err := build.Import(dir, ".", build.FindOnly)
		if err != nil {
			return "", fmt.Errorf("can't find %q: %v", path, err)
		}
		r.pkg = pkg

		src := filepath.Join(r.pkg.SrcRoot, path)
		if _, err := os.Stat(src); err == nil {
			r.hasSrcRoot = true
			return src, nil
		}

		src = filepath.Join(r.pkg.Dir, file)
		if _, err := os.Stat(src); err == nil {
			r.hasSrcRoot = false
			r.prefix = strings.TrimSuffix(r.pkg.ImportPath, strings.TrimPrefix(r.pkg.Dir, r.pkg.Root))
			return src, nil
		}
	}

	if r.hasSrcRoot {
		return filepath.Join(r.pkg.SrcRoot, path), nil
	}
	return filepath.Join(r.pkg.Root, strings.TrimPrefix(path, r.prefix)), nil
}
