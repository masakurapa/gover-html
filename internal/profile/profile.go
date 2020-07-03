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

	tree.mergeSingreDir()
	return tree
}

func (tree *Tree) add(paths []string, p *Profile) {
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

func (tree *Tree) mergeSingreDir() {
	tt := *tree
	for i, t := range *tree {
		t.SubDirs.mergeSingreDir()

		if len(t.Profiles) == 0 && len(t.SubDirs) == 1 {
			sub := t.SubDirs[0]
			tt[i].Name = filepath.Join(t.Name, sub.Name)
			tt[i].Profiles = sub.Profiles
			tt[i].SubDirs = sub.SubDirs
		}
	}
}
