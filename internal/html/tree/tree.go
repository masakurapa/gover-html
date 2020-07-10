package tree

import (
	"path/filepath"
	"strings"

	"github.com/masakurapa/go-cover/internal/profile"
)

// Node is single node of directory tree
type Node struct {
	Name  string
	Files []profile.Profile
	Dirs  []Node
}

// Create returns directory tree
func Create(profiles []profile.Profile) []Node {
	nodes := make([]Node, 0)
	for _, p := range profiles {
		addNode(&nodes, strings.Split(p.FileName, "/"), &p)
	}

	return mergeSingreDir(nodes)
}

func addNode(nodes *[]Node, paths []string, p *profile.Profile) {
	name := paths[0]
	nextPaths := paths[1:]

	idx := index(*nodes, name)
	if idx == -1 {
		*nodes = append(*nodes, Node{
			Name:  name,
			Files: make([]profile.Profile, 0),
			Dirs:  make([]Node, 0),
		})
		idx = len(*nodes) - 1
	}

	n := *nodes
	if len(nextPaths) == 1 {
		n[idx].Files = append(n[idx].Files, *p)
		return
	}

	addNode(&n[idx].Dirs, nextPaths, p)
}

func index(nodes []Node, name string) int {
	for i, t := range nodes {
		if t.Name == name {
			return i
		}
	}
	return -1
}

func mergeSingreDir(nodes []Node) []Node {
	for i, t := range nodes {
		mergeSingreDir(t.Dirs)

		if len(t.Files) == 0 && len(t.Dirs) == 1 {
			sub := t.Dirs[0]
			nodes[i].Name = filepath.Join(t.Name, sub.Name)
			nodes[i].Files = sub.Files
			nodes[i].Dirs = sub.Dirs
		}
	}

	return nodes
}
