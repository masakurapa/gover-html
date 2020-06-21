package html

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/masakurapa/go-cover/internal/profile"
)

func Print(profiles profile.Profiles, tree profile.Tree) {
	files := make([]TemplateFile, 0)
	for _, p := range profiles {

		file := findFile(p.FileName)
		b, err := ioutil.ReadFile(file)
		if err != nil {
			panic(fmt.Sprintf("can't read %q: %v", p.FileName, err))
		}

		var buf bytes.Buffer
		if err = genSource(&buf, b, p); err != nil {
			panic(err)
		}

		files = append(files, TemplateFile{
			Name:     p.FileName,
			Body:     template.HTML(buf.String()),
			Coverage: coverage(p.Blocks),
		})
	}

	var out *os.File
	out, err := os.Create("./my_coverage.html")
	if err != nil {
		panic(err)
	}

	err = parsedTemplate.Execute(out, TemplateData{
		Tree:  template.HTML(genDirectoryTree(tree)),
		Files: files,
	})
	if err != nil {
		panic(err)
	}
	defer out.Close()
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

func coverage(blocks []profile.Block) float64 {
	var total, covered int64
	for _, b := range blocks {
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}

	if total == 0 {
		return 0
	}

	return float64(covered) / float64(total) * 100
}

func findFile(file string) string {
	dir, file := filepath.Split(file)
	pkg, err := build.Import(dir, ".", build.FindOnly)
	if err != nil {
		panic(fmt.Sprintf("can't find %q: %v", file, err))
	}
	return filepath.Join(pkg.Dir, file)
}

func genSource(w io.Writer, src []byte, prof profile.Profile) error {
	dst := bufio.NewWriter(w)

	bi := 0
	si := 0
	line := 1
	col := 1

	for si < len(src) {
		if len(prof.Blocks) > bi {
			block := prof.Blocks[bi]
			if block.StartLine == line && block.StartCol == col {
				if block.Count == 0 {
					dst.WriteString(`<span class="uncov">`)
				} else {
					dst.WriteString(`<span class="cov">`)
				}
			}
			if block.EndLine == line && block.EndCol == col || line > block.EndLine {
				dst.WriteString(`</span>`)
				bi++
				continue
			}
		}

		b := src[si]
		switch b {
		case '<':
			dst.WriteString("&lt;")
		case '>':
			dst.WriteString("&gt;")
		case '&':
			dst.WriteString("&amp;")
		case '\t':
			dst.WriteString("        ")
		default:
			dst.WriteByte(b)
		}

		if b == '\n' {
			line++
			col = 0
		}

		si++
		col++
	}

	return dst.Flush()
}
