package reader

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/masakurapa/go-cover/internal/logger"
)

type Reader interface {
	Read(path string) ([]byte, error)
}

type reader struct {
	dirs map[string]string
}

func New() Reader {
	return &reader{dirs: make(map[string]string)}
}

func (r *reader) Read(path string) ([]byte, error) {
	file, err := r.find2(path)
	if err != nil {
		return nil, err
	}

	logger.L.Debug("read: %s", file)
	return ioutil.ReadFile(file)
}

func (r *reader) find(file string) (string, error) {
	dir, file := filepath.Split(file)

	if d, ok := r.dirs[dir]; ok {
		return filepath.Join(d, file), nil
	}

	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("can't find %q: %v", file, err)
	}

	r.dirs[dir] = pkg.Dir
	return filepath.Join(pkg.Dir, file), nil
}

func (r *reader) find2(file string) (string, error) {
	for k, v := range r.dirs {
		if !strings.HasPrefix(file, k) {
			continue
		}
		p := strings.TrimPrefix(file, k)
		return filepath.Join(v, p), nil
	}

	dir, file := filepath.Split(file)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("can't find %q: %v", file, err)
	}

	p := strings.TrimSuffix(pkg.ImportPath, strings.TrimPrefix(pkg.Dir, pkg.Root))
	r.dirs[p] = pkg.Root
	return filepath.Join(pkg.Dir, file), nil
}
