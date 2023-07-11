package cover

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"

	"github.com/masakurapa/gover-html/internal/cover/filter"
	"github.com/masakurapa/gover-html/internal/profile"
)

type funcExtent struct {
	structName string
	funcName   string
	name       string
	startLine  int
	startCol   int
	endLine    int
	endCol     int
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

		var structName string
		if n.Recv != nil && len(n.Recv.List) > 0 {
			structName = n.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name
		}

		fe := &funcExtent{
			structName: structName,
			funcName:   n.Name.Name,
			name:       string(src),
			startLine:  start.Line,
			startCol:   start.Column,
			endLine:    end.Line,
			endCol:     end.Column,
		}
		v.funcs = append(v.funcs, fe)
	}
	return v
}

func makeNewProfile(prof *profile.Profile, f filter.Filter) (*profile.Profile, error) {
	exts, err := findFuncs(prof.FilePath())
	if err != nil {
		return nil, err
	}

	return newProfile(prof, exts, f), nil
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

func newProfile(prof *profile.Profile, exts []*funcExtent, f filter.Filter) *profile.Profile {
	fncs := make(profile.Functions, 0, len(exts))
	blocks := make(profile.Blocks, 0, len(prof.Blocks))

	bi := 0
	for _, e := range exts {
		isCoverageBlock := f.IsOutputTargetFunc(prof.RemoveModulePathFromFileName(), e.structName, e.funcName)

		fnc := profile.Function{
			Name:      e.name,
			StartLine: e.startLine,
			StartCol:  e.startCol,
		}

		for bi < len(prof.Blocks) {
			b := prof.Blocks[bi]

			if b.StartLine < e.startLine {
				if isCoverageBlock {
					blocks = append(blocks, b)
				}
				bi++
				continue
			}

			if b.StartLine >= e.startLine &&
				b.EndLine <= e.endLine {
				fnc.Blocks = append(fnc.Blocks, b)
				if isCoverageBlock {
					blocks = append(blocks, b)
				}
				bi++
				continue
			}
			break
		}

		if isCoverageBlock {
			fncs = append(fncs, fnc)
		}
	}

	prof.Blocks = blocks
	prof.Functions = fncs
	return prof
}
