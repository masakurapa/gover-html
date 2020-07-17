package cover

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/masakurapa/go-cover/internal/profile"
)

type funcExtent struct {
	name      string
	startLine int
	startCol  int
	endLine   int
	endCol    int
}

type funcVisitor struct {
	fset    *token.FileSet
	name    string
	astFile *ast.File
	funcs   []*funcExtent
}

func (v *funcVisitor) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if n.Body == nil {
			break
		}
		start := v.fset.Position(n.Pos())
		end := v.fset.Position(n.End())

		var buf bytes.Buffer
		err := format.Node(&buf, token.NewFileSet(), &ast.FuncDecl{
			Name: n.Name,
			Recv: n.Recv,
			Type: n.Type,
		})
		if err != nil {
			panic(err)
		}
		src, err := format.Source(buf.Bytes())
		if err != nil {
			panic(err)
		}

		fe := &funcExtent{
			name:      string(src),
			startLine: start.Line,
			startCol:  start.Column,
			endLine:   end.Line,
			endCol:    end.Column,
		}
		v.funcs = append(v.funcs, fe)
	}
	return v
}

func makeFuncs(prof profile.Profile) (profile.Functions, error) {
	exts, err := findFuncs(prof.FilePath())
	if err != nil {
		return nil, err
	}

	return toFunctions(prof, exts), nil
}

func findFuncs(name string) ([]*funcExtent, error) {
	fset := token.NewFileSet()
	parsedFile, err := parser.ParseFile(fset, name, nil, 0)
	if err != nil {
		return nil, err
	}
	visitor := &funcVisitor{
		fset:    fset,
		name:    name,
		astFile: parsedFile,
	}
	ast.Walk(visitor, visitor.astFile)
	return visitor.funcs, nil
}

func toFunctions(prof profile.Profile, exts []*funcExtent) profile.Functions {
	fncs := make(profile.Functions, 0, len(exts))

	bi := 0
	for _, e := range exts {
		fnc := profile.Function{
			Name:      e.name,
			StartLine: e.startLine,
			StartCol:  e.startCol,
		}

		for bi < len(prof.Blocks) {
			b := prof.Blocks[bi]

			if b.StartLine < e.startLine {
				bi++
				continue
			}

			if b.StartLine >= e.startLine &&
				b.EndLine <= e.endLine {
				fnc.Blocks = append(fnc.Blocks, b)
				bi++
				continue
			}
			break
		}

		fncs = append(fncs, fnc)
	}

	return fncs
}
