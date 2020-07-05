package profile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
)

type Packages map[string]*Package

type Package struct {
	Dir        string
	ImportPath string
	Error      *struct {
		Err string
	}
}

func (profs *Profiles) ToPackages() (Packages, error) {
	pkgs := make(Packages)
	dirs := make([]string, 0, len(*profs))
	for _, prof := range *profs {
		if prof.IsRelativeOrAbsolute() {
			continue
		}

		dir := path.Dir(prof.FileName)
		if _, ok := pkgs[dir]; !ok {
			pkgs[dir] = nil
			dirs = append(dirs, dir)
		}
	}

	if len(dirs) == 0 {
		return pkgs, nil
	}

	cmdName := filepath.Join(runtime.GOROOT(), "bin/go")
	args := append([]string{"list", "-e", "-json"}, dirs...)
	cmd := exec.Command(cmdName, args...)

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(bytes.NewReader(stdout))
	for {
		var pkg Package
		err := dec.Decode(&pkg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decoding go list json: %v", err)
		}
		pkgs[pkg.ImportPath] = &pkg
	}

	return pkgs, nil
}
