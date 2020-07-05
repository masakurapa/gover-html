package reader

import (
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/masakurapa/go-cover/internal/profile"
)

type Reader interface {
	Read(profile.Packages, string) ([]byte, error)
}

type reader struct{}

func New() Reader {
	return &reader{}
}

func (r *reader) Read(pkgs profile.Packages, path string) ([]byte, error) {
	file, err := r.find(pkgs, path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(file)
}

func (r *reader) find(pkgs profile.Packages, file string) (string, error) {
	pkg := pkgs[path.Dir(file)]
	if pkg != nil {
		if pkg.Dir != "" {
			return filepath.Join(pkg.Dir, path.Base(file)), nil
		}
		if pkg.Error != nil {
			return "", fmt.Errorf(pkg.Error.Err)
		}
	}
	return "", fmt.Errorf("file not found. %s", file)
}
