package reader

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"
)

type Reader interface {
	Read(path string) ([]byte, error)
}

type reader struct{}

func New() Reader {
	return &reader{}
}

func (r *reader) Read(path string) ([]byte, error) {
	file, err := r.find(path)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(file)
}

func (*reader) find(file string) (string, error) {
	dir, file := filepath.Split(file)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		return "", fmt.Errorf("can't find %q: %v", file, err)
	}
	return filepath.Join(pkg.Dir, file), nil
}
