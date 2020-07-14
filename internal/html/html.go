package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"

	"github.com/masakurapa/go-cover/internal/html/tree"
	"github.com/masakurapa/go-cover/internal/profile"
)

var (
	escapeChar = map[byte]string{
		'<':  "&lt;",
		'>':  "&gt;",
		'&':  "&amp;",
		'\t': "        ",
	}
)

type templateData struct {
	Tree  []templateTree
	Files []templateFile
}

type templateTree struct {
	ID       int
	Name     string
	Indent   int
	Coverage float64
	IsFile   bool
}

type templateFile struct {
	ID       int
	Name     string
	Body     template.HTML
	Coverage float64
}

// WriteTreeView outputs coverage as a tree view HTML file
func WriteTreeView(out io.Writer, profiles []profile.Profile) error {
	nodes := tree.Create(profiles)
	tree := make([]templateTree, 0)
	makeTemplateTree(&tree, nodes, 0)

	data := templateData{
		Tree:  tree,
		Files: make([]templateFile, 0, len(profiles)),
	}

	var buf bytes.Buffer
	for _, p := range profiles {
		b, err := ioutil.ReadFile(p.FilePath())
		if err != nil {
			return fmt.Errorf("can't read %q: %v", p.FileName, err)
		}

		writeSource(&buf, b, &p)
		data.Files = append(data.Files, templateFile{
			ID:       p.ID,
			Name:     p.FileName,
			Body:     template.HTML(buf.String()),
			Coverage: p.Coverage(),
		})
		buf.Reset()
	}

	return parsedTreeTemplate.Execute(out, data)
}

func makeTemplateTree(tree *[]templateTree, nodes []tree.Node, indent int) {
	for _, node := range nodes {
		*tree = append(*tree, templateTree{
			Name:   node.Name,
			Indent: indent,
		})

		makeTemplateTree(tree, node.Dirs, indent+1)

		for _, p := range node.Files {
			*tree = append(*tree, templateTree{
				ID:       p.ID,
				Name:     path.Base(p.FileName),
				Indent:   indent + 1,
				Coverage: p.Coverage(),
				IsFile:   true,
			})
		}
	}
}

func writeSource(buf *bytes.Buffer, src []byte, prof *profile.Profile) {
	bi := 0
	si := 0
	line := 1
	col := 1
	cov0 := false
	cov1 := false

	buf.WriteString("<ol>")

	for si < len(src) {
		if col == 1 {
			buf.WriteString("<li>")
			if cov0 {
				buf.WriteString(`<span class="cov0">`)
			}
			if cov1 {
				buf.WriteString(`<span class="cov1">`)
			}
		}

		if len(prof.Blocks) > bi {
			block := prof.Blocks[bi]
			if block.StartLine == line && block.StartCol == col {
				if block.Count == 0 {
					buf.WriteString(`<span class="cov0">`)
					cov0 = true
				} else {
					buf.WriteString(`<span class="cov1">`)
					cov1 = true
				}
			}
			if block.EndLine == line && block.EndCol == col || line > block.EndLine {
				buf.WriteString(`</span>`)
				bi++
				cov0 = false
				cov1 = false
				continue
			}
		}

		b := src[si]
		buf.WriteString(escape(b))

		if b == '\n' {
			if cov0 || cov1 {
				buf.WriteString("</span>")
			}
			buf.WriteString("</li>")
			line++
			col = 0
		}

		si++
		col++
	}

	buf.WriteString("</ol>")
}

func escape(b byte) string {
	if s, ok := escapeChar[b]; ok {
		return s
	}
	return string(b)
}
