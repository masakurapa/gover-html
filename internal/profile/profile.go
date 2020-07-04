package profile

import (
	"strings"
)

type Profiles []Profile

type Profile struct {
	ID       int
	FileName string
	Mode     string
	Blocks   Blocks
}

func (prof *Profiles) ToTree() Tree {
	tree := make(Tree, 0)
	for _, p := range *prof {
		tree.add(strings.Split(p.FileName, "/"), &p)
	}

	tree.mergeSingreDir()
	return tree
}
