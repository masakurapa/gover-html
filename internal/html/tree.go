package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"sync"

	"github.com/masakurapa/go-cover/internal/logger"
	"github.com/masakurapa/go-cover/internal/profile"
	"github.com/masakurapa/go-cover/internal/reader"
)

var (
	writePool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

type TreeTemplateData struct {
	Tree  template.HTML
	Files []File
}

type File struct {
	Name     string
	Body     template.HTML
	Line     int
	Coverage float64
}

func WriteTreeView(
	reader reader.Reader,
	out io.Writer,
	profiles profile.Profiles,
	tree profile.Tree,
) error {
	logger.L.Debug("start generating tree view html")

	data := TreeTemplateData{
		Tree:  template.HTML(genDirectoryTree(tree)),
		Files: make([]File, 0, len(profiles)),
	}

	for _, p := range profiles {
		b, err := reader.Read(p.FileName)
		if err != nil {
			return fmt.Errorf("can't read %q: %v", p.FileName, err)
		}

		buf := writePool.Get().(*bytes.Buffer)
		line := genSource(buf, b, &p)

		data.Files = append(data.Files, File{
			Name:     p.FileName,
			Body:     template.HTML(buf.String()),
			Line:     line,
			Coverage: p.Blocks.Coverage(),
		})

		buf.Reset()
		writePool.Put(buf)
	}

	logger.L.Debug("write html")
	return parsedTreeTemplate.Execute(out, data)
}

func genDirectoryTree(tree profile.Tree) string {
	buf := writePool.Get().(*bytes.Buffer)
	makeDirectoryTree(buf, tree, 0)
	s := buf.String()
	buf.Reset()
	writePool.Put(buf)
	return s
}

// ディレクトリツリーの生成
func makeDirectoryTree(buf *bytes.Buffer, tree profile.Tree, indent int) {
	for _, t := range tree {
		buf.WriteString(fmt.Sprintf(
			`<div style="padding-inline-start: %dpx">%s</div>`,
			indent*30,
			t.Name,
		))

		makeDirectoryTree(buf, t.SubDirs, indent+1)

		for _, p := range t.Profiles {
			_, f := filepath.Split(p.FileName)
			buf.WriteString(fmt.Sprintf(
				`<div class="file" style="padding-inline-start: %dpx" id="tree%d" onclick="change(%d);">%s (%.1f%%)</div>`,
				(indent+1)*30,
				p.ID,
				p.ID,
				f,
				p.Blocks.Coverage(),
			))
		}
	}
}

func genSource(buf *bytes.Buffer, src []byte, prof *profile.Profile) int {
	bi := 0
	si := 0
	line := 1
	col := 1

	for si < len(src) {
		if len(prof.Blocks) > bi {
			block := prof.Blocks[bi]
			if block.StartLine == line && block.StartCol == col {
				if block.Count == 0 {
					buf.WriteString(`<span class="cov0">`)
				} else {
					buf.WriteString(`<span class="cov1">`)
				}
			}
			if block.EndLine == line && block.EndCol == col || line > block.EndLine {
				buf.WriteString(`</span>`)
				bi++
				continue
			}
		}

		b := src[si]
		switch b {
		case '<':
			buf.WriteString("&lt;")
		case '>':
			buf.WriteString("&gt;")
		case '&':
			buf.WriteString("&amp;")
		case '\t':
			buf.WriteString("        ")
		default:
			buf.WriteByte(b)
		}

		if b == '\n' {
			line++
			col = 0
		}

		si++
		col++
	}

	return line
}
