package tree

import (
	"path/filepath"
	"strings"

	"github.com/masakurapa/gover-html/internal/profile"
)

// Node is single node of directory tree
type Node struct {
	Name  string
	Files profile.Profiles
	Dirs  []Node
}

// ChildBlocks returns all child Blocks for Node
func (n *Node) ChildBlocks() profile.Blocks {
	blocks := make(profile.Blocks, 0)
	for _, f := range n.Files {
		blocks = append(blocks, f.Blocks...)
	}
	for _, ch := range n.Dirs {
		blocks = append(blocks, ch.ChildBlocks()...)
	}
	return blocks
}

// Create returns directory tree
func Create(profiles profile.Profiles) []Node {
	nodes := make([]Node, 0)
	for _, p := range profiles {
		idx := index(nodes, p.ModulePath)
		if idx == -1 {
			nodes = append(nodes, Node{Name: p.ModulePath})
			idx = len(nodes) - 1
		}

		addNode(&nodes[idx].Dirs, strings.Split(p.RemoveModulePathFromFileName(), "/"), &p)
	}

	for i, node := range nodes {
		nodes[i].Dirs = mergeSingreDir(node.Dirs)
	}
	return nodes
}

func addNode(nodes *[]Node, paths []string, p *profile.Profile) {
	name := paths[0]
	nextPaths := paths[1:]

	idx := index(*nodes, name)
	if idx == -1 {
		*nodes = append(*nodes, Node{Name: name})
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

// merge directories with no files and only one child directory
//
//	path/
//		to/
//		file.go
//
// to
//
//	path/to/
//		file.go
func mergeSingreDir(nodes []Node) []Node {
	for i, n := range nodes {
		if len(n.Dirs) == 0 {
			continue
		}

		mergeSingreDir(n.Dirs)
		if len(n.Files) > 0 || len(n.Dirs) != 1 {
			continue
		}

		sub := n.Dirs[0]
		nodes[i].Name = filepath.Join(n.Name, sub.Name)
		nodes[i].Files = sub.Files
		nodes[i].Dirs = sub.Dirs
	}
	return nodes
}
