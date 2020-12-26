package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path"

	"github.com/masakurapa/gover-html/internal/html/tree"
	"github.com/masakurapa/gover-html/internal/option"
	"github.com/masakurapa/gover-html/internal/profile"
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
	Theme string
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
	ID        int
	Name      string
	Body      template.HTML
	Coverage  float64
	Functions []templateFunc
}

type templateFunc struct {
	Name     string
	Line     int
	Coverage float64
}

// WriteTreeView outputs coverage as a tree view HTML file
func WriteTreeView(out io.Writer, profiles profile.Profiles, opt option.Option) error {
	nodes := tree.Create(profiles)
	tree := make([]templateTree, 0)
	makeTemplateTree(&tree, nodes, 0)

	data := templateData{
		Theme: opt.Theme,
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
		f := templateFile{
			ID:        p.ID,
			Name:      p.FileName,
			Body:      template.HTML(buf.String()),
			Coverage:  p.Blocks.Coverage(),
			Functions: make([]templateFunc, 0, len(p.Functions)),
		}
		buf.Reset()

		for _, fn := range p.Functions {
			f.Functions = append(f.Functions, templateFunc{
				Name:     fn.Name,
				Line:     fn.StartLine,
				Coverage: fn.Blocks.Coverage(),
			})
		}

		data.Files = append(data.Files, f)
	}

	return parsedTreeTemplate.Execute(out, data)
}

func makeTemplateTree(tree *[]templateTree, nodes []tree.Node, indent int) {
	for _, node := range nodes {
		childBlocks := node.ChildBlocks()
		*tree = append(*tree, templateTree{
			Name:     node.Name,
			Indent:   indent,
			Coverage: childBlocks.Coverage(),
		})

		makeTemplateTree(tree, node.Dirs, indent+1)

		for _, p := range node.Files {
			*tree = append(*tree, templateTree{
				ID:       p.ID,
				Name:     path.Base(p.FileName),
				Indent:   indent + 1,
				Coverage: p.Blocks.Coverage(),
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
			buf.WriteString(fmt.Sprintf(`<li id="file%d-%d">`, prof.ID, line))
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
		writeChar(buf, b)

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

func writeChar(buf *bytes.Buffer, b byte) {
	if s, ok := escapeChar[b]; ok {
		buf.WriteString(s)
		return
	}
	buf.WriteByte(b)
}
