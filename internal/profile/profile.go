package profile

import (
	"path/filepath"
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
