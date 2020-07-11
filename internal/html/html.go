package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"

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
	Tree  template.HTML
	Files []templateFile
}

type templateFile struct {
	ID       int
	Name     string
	Body     template.HTML
	Coverage float64
}

// WriteTreeView outputs coverage as a tree view HTML file
func WriteTreeView(out io.Writer, profiles []profile.Profile) error {
	data := templateData{
		Tree:  template.HTML(directoryTree(profiles)),
		Files: make([]templateFile, 0, len(profiles)),
	}

	for _, p := range profiles {
		b, err := ioutil.ReadFile(p.FilePath())
		if err != nil {
			return fmt.Errorf("can't read %q: %v", p.FileName, err)
		}

		var buf bytes.Buffer
		writeSource(&buf, b, &p)

		data.Files = append(data.Files, templateFile{
			ID:       p.ID,
			Name:     p.FileName,
			Body:     template.HTML(buf.String()),
			Coverage: p.Coverage(),
		})
	}

	return parsedTreeTemplate.Execute(out, data)
}

func directoryTree(profiles []profile.Profile) string {
	nodes := tree.Create(profiles)
	var buf bytes.Buffer
	writeDirectoryTree(&buf, nodes, 0)
	s := buf.String()
	return s
}

func writeDirectoryTree(buf *bytes.Buffer, nodes []tree.Node, indent int) {
	for _, node := range nodes {
		buf.WriteString(fmt.Sprintf(
			`<div style="padding-inline-start: %dpx">%s</div>`,
			indent*30,
			node.Name,
		))

		writeDirectoryTree(buf, node.Dirs, indent+1)

		for _, p := range node.Files {
			_, f := filepath.Split(p.FileName)
			buf.WriteString(fmt.Sprintf(
				`<div class="file" style="padding-inline-start: %dpx" id="tree%d" onclick="change(%d);">%s (%.1f%%)</div>`,
				(indent+1)*30,
				p.ID,
				p.ID,
				f,
				p.Coverage(),
			))
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
