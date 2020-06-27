package reader

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"

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
	file, err := r.find(path)
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
