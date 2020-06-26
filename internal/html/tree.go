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

func WriteTreeView2(
	reader reader.Reader,
	out io.Writer,
	profiles profile.Profiles,
	tree profile.Tree,
) error {
	var buf bytes.Buffer
	buf.WriteString(`<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
	<style>
		.content {
			display: flex;
			padding: 8px;
		}
		.tree {
			width: 30%;
			white-space: nowrap;
			overflow-x: scroll;
		}
		.cover {
			width: 70%;
			margin-left: 32px;
		}
		.source {
			white-space: nowrap;
			overflow-x: scroll;
		}
		ul {
			list-style: none;
			margin-block-start: unset;
			margin-block-end: unset;
		}
		li .file {
			cursor: pointer;
		}
		pre {
			font-family: Menlo, monospace;
			font-weight: bold;
			color: rgb(80, 80, 80);
		}
		.uncov {
			color: rgb(192, 0, 0);
		}
		.cov {
			color: rgb(44, 212, 149);
		}
	</style>
</head>
<body>
	<div class="content">
		<div class="tree">`)

	buf.WriteString(genDirectoryTree(tree))
	buf.WriteString(`</div>
		<div class="cover">
`)

	for i, p := range profiles {
		b, err := reader.Read(p.FileName)
		if err != nil {
			return fmt.Errorf("can't read %q: %v", p.FileName, err)
		}

		buf.WriteString(fmt.Sprintf(`<div id="file%d" style="display: none">
		<div>%s</div>
		<div>%.1f%%</div>
		<div class="source">
			<pre>`, i, p.FileName, p.Blocks.Coverage()))

		if err = genSource(&buf, b, &p); err != nil {
			return err
		}

		buf.WriteString(`</pre>
		</div>
	</div>`)
	}

	buf.WriteString(`</div>
	</div>

	<script>
		let current;

		function select(f) {
			if (current) {
				current.style.display = 'none';
			}

			current = document.getElementById(f)
			if (!current) {
				return;
			}

			current.style.display = 'block';
		}
		function change(f) {
			select(f);
			window.scrollTo(0, 0);
		}
	</script>
</body>
</html>`)

	out.Write(buf.Bytes())
	return nil
}

func WriteTreeView(
	reader reader.Reader,
	out io.Writer,
	profiles profile.Profiles,
	tree profile.Tree,
) error {
	logger.L.Debug("start generating tree view html")

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

	logger.L.Debug("write html")
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
