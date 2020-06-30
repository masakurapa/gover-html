package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	"github.com/masakurapa/go-cover/internal/logger"
	"github.com/masakurapa/go-cover/internal/profile"
	"github.com/masakurapa/go-cover/internal/reader"
)

func WriteTreeView(
	reader reader.Reader,
	out io.Writer,
	profiles profile.Profiles,
	tree profile.Tree,
) error {
	logger.L.Debug("start generating tree view html")

	data := TreeTemplateData{
		Tree:  template.HTML(genDirectoryTree(tree)),
		Files: make([]TemplateFile, 0, len(profiles)),
	}

	for _, p := range profiles {
		b, err := reader.Read(p.FileName)
		if err != nil {
			return fmt.Errorf("can't read %q: %v", p.FileName, err)
		}

		buf := writePool.Get().(*bytes.Buffer)
		genSource(buf, b, &p)

		data.Files = append(data.Files, TemplateFile{
			Name:     p.FileName,
			Body:     template.HTML(buf.String()),
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
