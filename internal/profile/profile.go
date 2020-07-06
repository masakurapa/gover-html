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
	"strings"
)

type Profiles []Profile

type Profile struct {
	ID       int
	FileName string
	Mode     string
	Blocks   Blocks
}

func (profs *Profiles) ToTree() Tree {
	tree := make(Tree, 0)
	for _, p := range *profs {
		tree.add(strings.Split(p.FileName, "/"), &p)
	}

	tree.mergeSingreDir()
	return tree
}

func (prof *Profile) IsRelativeOrAbsolute() bool {
	return strings.HasPrefix(prof.FileName, ".") || filepath.IsAbs(prof.FileName)
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
