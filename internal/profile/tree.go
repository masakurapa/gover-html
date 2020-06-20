package profile

import (
	"strings"

	"golang.org/x/tools/cover"
)

type Tree []Node

type Node struct {
	Name     string
	Profiles Profiles
	SubDirs  Tree
}

func (prof Profiles) ToTree() Tree {
	tree := make(Tree, 0)

	for _, p := range prof {
		tree.add(strings.Split(p.FileName, "/"), &p)
	}

	return tree
}

func (tree *Tree) add(paths []string, p *cover.Profile) {
	name := paths[0]
	nextPaths := paths[1:]

	idx, ok := tree.findNode(name)
	if !ok {
		*tree = append(*tree, Node{
			Name:     name,
			Profiles: make(Profiles, 0),
			SubDirs:  make(Tree, 0),
		})
		idx = len(*tree) - 1
	}

	t := *tree
	if len(nextPaths) == 1 {
		t[idx].Profiles = append(t[idx].Profiles, *p)
		return
	}
	t[idx].SubDirs.add(nextPaths, p)
}

func (tree *Tree) findNode(name string) (int, bool) {
	for i, t := range *tree {
		if t.Name == name {
			return i, true
		}
	}
	return 0, false
}
