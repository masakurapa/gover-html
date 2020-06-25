package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	"github.com/masakurapa/go-cover/internal/profile"
	"github.com/masakurapa/go-cover/internal/reader"
)

func WriteTreeView(
	reader reader.Reader,
	out io.Writer,
	profiles profile.Profiles,
	tree profile.Tree,
) error {
	files := make([]TemplateFile, 0, len(profiles))
	for _, p := range profiles {
		b, err := reader.Read(p.FileName)
		if err != nil {
			return fmt.Errorf("can't read %q: %v", p.FileName, err)
		}

		var buf bytes.Buffer
		if err = genSource(&buf, b, &p); err != nil {
			return err
		}

		files = append(files, TemplateFile{
			Name:     p.FileName,
			Body:     template.HTML(buf.String()),
			Coverage: p.Blocks.Coverage(),
		})
	}

	return parsedTreeTemplate.Execute(out, TreeTemplateData{
		Tree:  template.HTML(genDirectoryTree(tree)),
		Files: files,
	})
}

func genDirectoryTree(tree profile.Tree) string {
	var buf bytes.Buffer
	buf.WriteString(`<ul style="padding-left: 0;">`)
	makeDirectoryTree(&buf, tree)
	buf.WriteString("</ul>")
	return buf.String()
}

// ディレクトリツリーの生成
func makeDirectoryTree(buf *bytes.Buffer, tree profile.Tree) {
	for _, t := range tree {
		buf.WriteString("<li>")
		buf.WriteString(t.Name)

		if len(t.Profiles) > 0 || len(t.SubDirs) > 0 {
			buf.WriteString("<ul>")
		}

		makeDirectoryTree(buf, t.SubDirs)

		for _, p := range t.Profiles {
			_, f := filepath.Split(p.FileName)
			buf.WriteString(fmt.Sprintf(
				`<li class="file" onclick="change('file%d');">%s (%.1f%%)</li>`,
				p.ID,
				f,
				p.Blocks.Coverage(),
			))
		}

		if len(t.Profiles) > 0 || len(t.SubDirs) > 0 {
			buf.WriteString("</ul>")
		}
		buf.WriteString("</li>")
	}
}
