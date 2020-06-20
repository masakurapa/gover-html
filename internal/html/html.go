package html

import (
	"fmt"
	"go/build"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/masakurapa/gocover-html/internal/profile"
	"golang.org/x/tools/cover"
)

func Print(profiles profile.Profiles, tree profile.Tree) {
	files := make([]TemplateFile, 0)
	for _, p := range profiles {

		file := findFile(p.FileName)
		b, err := ioutil.ReadFile(file)
		if err != nil {
			panic(fmt.Sprintf("can't read %q: %v", p.FileName, err))
		}

		files = append(files, TemplateFile{
			Name:     p.FileName,
			Body:     template.HTML(b),
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
	return "<ul>" + makeDirectoryTree(tree) + "</ul>"
}

// ディレクトリツリーの生成
func makeDirectoryTree(tree profile.Tree) string {
	tag := ""
	for _, t := range tree {
		tag += "<li>" + t.Name
		if len(t.Child) > 0 {
			tag += "<ul>"
		}

		tag += makeDirectoryTree(t.Child)

		if len(t.Child) > 0 {
			tag += "</ul>"
		}
		tag += "</li>"
	}
	return tag
}

func coverage(blocks []cover.ProfileBlock) float64 {
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
