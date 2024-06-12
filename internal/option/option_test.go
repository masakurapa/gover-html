package option_test

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/masakurapa/gover-html/internal/option"
	"github.com/masakurapa/gover-html/internal/reader"
	"github.com/masakurapa/gover-html/test/helper"
)

type mockReader struct {
	reader.Reader
	mockRead   func(string) (io.Reader, error)
	mockExists func(string) bool
}

func (m *mockReader) Read(s string) (io.Reader, error) {
	return m.mockRead(s)
}

func (m *mockReader) Exists(s string) bool {
	return m.mockExists(s)
}

func TestNew(t *testing.T) {
	type args struct {
		input       *string
		inputFiles  *string
		output      *string
		theme       *string
		include     *string
		exclude     *string
		excludeFunc *string
	}
	type testCase []struct {
		name     string
		settings string
		args     args
		want     *option.Option
		wantErr  bool
	}

	t.Run("設定ファイルが存在しない", func(t *testing.T) {
		tests := testCase{
			{
				name: "全項目に設定値が存在(theme=dark)",
				args: args{
					input:       helper.P("example.out"),
					inputFiles:  helper.P("example1.out,example2.out"),
					output:      helper.P("example.html"),
					theme:       helper.P("dark"),
					include:     helper.P("path/to/dir1,path/to/dir2"),
					exclude:     helper.P("path/to/dir3,path/to/dir4"),
					excludeFunc: helper.P("(path/to/dir3).Func1,(path/to/dir4.Struct1).Func2"),
				},
				want: &option.Option{
					Input:      "example.out",
					InputFiles: []string{"example1.out", "example2.out"},
					Output:     "example.html",
					Theme:      "dark",
					Include:    []string{"path/to/dir1", "path/to/dir2"},
					Exclude:    []string{"path/to/dir3", "path/to/dir4"},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "path/to/dir3", Struct: "", Func: "Func1"},
						{Path: "path/to/dir4", Struct: "Struct1", Func: "Func2"},
					},
				},
				wantErr: false,
			},
			{
				name: "全項目に設定値が存在(theme=light)",
				args: args{
					input:       helper.P("example.out"),
					inputFiles:  helper.P("example1.out,example2.out"),
					output:      helper.P("example.html"),
					theme:       helper.P("light"),
					include:     helper.P("path/to/dir1,path/to/dir2"),
					exclude:     helper.P("path/to/dir3,path/to/dir4"),
					excludeFunc: helper.P("(path/to/dir3).Func1,(path/to/dir4.Struct1).Func2"),
				},
				want: &option.Option{
					Input:      "example.out",
					InputFiles: []string{"example1.out", "example2.out"},
					Output:     "example.html",
					Theme:      "light",
					Include:    []string{"path/to/dir1", "path/to/dir2"},
					Exclude:    []string{"path/to/dir3", "path/to/dir4"},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "path/to/dir3", Struct: "", Func: "Func1"},
						{Path: "path/to/dir4", Struct: "Struct1", Func: "Func2"},
					},
				},
				wantErr: false,
			},
			{
				name: "全項目に空文字を指定",
				args: args{
					input:       helper.P(""),
					inputFiles:  helper.P(""),
					output:      helper.P(""),
					theme:       helper.P(""),
					include:     helper.P(""),
					exclude:     helper.P(""),
					excludeFunc: helper.P(""),
				},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "全項目にnilを指定",
				args: args{
					input:       nil,
					inputFiles:  nil,
					output:      nil,
					theme:       nil,
					include:     nil,
					exclude:     nil,
					excludeFunc: nil,
				},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "input-filesに空の値を持つ",
				args: args{
					inputFiles: helper.P("example1.out,,example2.out,"),
				},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{"example1.out", "example2.out"},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "includeに空の値を持つ",
				args: args{
					include: helper.P("path/to/dir1,,path/to/dir2,,"),
				},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{"path/to/dir1", "path/to/dir2"},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "include./で始まるパスを指定",
				args: args{
					include: helper.P("./path/to/dir1"),
				},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{"path/to/dir1"},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "includeに/で終わるパスを指定",
				args: args{
					include: helper.P("path/to/dir1/"),
				},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{"path/to/dir1"},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "includeに/で始まるパスを指定",
				args: args{
					include: helper.P("/path/to/dir1"),
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "excludeに空の値を持つ",
				args: args{
					exclude: helper.P("path/to/dir3,,path/to/dir4,,"),
				},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{"path/to/dir3", "path/to/dir4"},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "excludeに./で始まるパスを指定",
				args: args{
					exclude: helper.P("./path/to/dir3"),
				},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{"path/to/dir3"},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "excludeに/で終わるパスを指定",
				args: args{
					exclude: helper.P("path/to/dir3/"),
				},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{"path/to/dir3"},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "excludeに/で始まるパスを指定",
				args: args{
					exclude: helper.P("/path/to/dir3"),
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "exclude-funcに空の値を持つ",
				args: args{
					excludeFunc: helper.P("(path/to/dir3).Func1,,(path/to/dir4.Struct1).Func2,,"),
				},
				want: &option.Option{
					Input:      "coverage.out",
					InputFiles: []string{},
					Output:     "coverage.html",
					Theme:      "dark",
					Include:    []string{},
					Exclude:    []string{},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "path/to/dir3", Struct: "", Func: "Func1"},
						{Path: "path/to/dir4", Struct: "Struct1", Func: "Func2"},
					},
				},
				wantErr: false,
			},
			{
				name: "exclude-funcに./で始まるパスを指定",
				args: args{
					excludeFunc: helper.P("(./path/to/dir3.Struct1).Func1"),
				},
				want: &option.Option{
					Input:      "coverage.out",
					InputFiles: []string{},
					Output:     "coverage.html",
					Theme:      "dark",
					Include:    []string{},
					Exclude:    []string{},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "path/to/dir3", Struct: "Struct1", Func: "Func1"},
					},
				},
				wantErr: false,
			},
			{
				name: "exclude-funcに/で終わるパスを指定",
				args: args{
					excludeFunc: helper.P("(path/to/dir3/).Func1"),
				},
				want: &option.Option{
					Input:      "coverage.out",
					InputFiles: []string{},
					Output:     "coverage.html",
					Theme:      "dark",
					Include:    []string{},
					Exclude:    []string{},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "path/to/dir3", Struct: "", Func: "Func1"},
					},
				},
				wantErr: false,
			},
			{
				name: "exclude-funcに/で終わるパスを指定(構造体名指定",
				args: args{
					excludeFunc: helper.P("(path/to/dir3/.Struct1).Func1"),
				},
				want: &option.Option{
					Input:      "coverage.out",
					InputFiles: []string{},
					Output:     "coverage.html",
					Theme:      "dark",
					Include:    []string{},
					Exclude:    []string{},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "path/to/dir3", Struct: "Struct1", Func: "Func1"},
					},
				},
				wantErr: false,
			},
			{
				name: "exclude-funcのパスにワイルドカードを指定",
				args: args{
					excludeFunc: helper.P("(*.Struct1).Func1"),
				},
				want: &option.Option{
					Input:      "coverage.out",
					InputFiles: []string{},
					Output:     "coverage.html",
					Theme:      "dark",
					Include:    []string{},
					Exclude:    []string{},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "", Struct: "Struct1", Func: "Func1"},
					},
				},
				wantErr: false,
			},
			{
				name: "exclude-funcに関数名のみを指定",
				args: args{
					excludeFunc: helper.P("Func1"),
				},
				want: &option.Option{
					Input:      "coverage.out",
					InputFiles: []string{},
					Output:     "coverage.html",
					Theme:      "dark",
					Include:    []string{},
					Exclude:    []string{},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "", Struct: "", Func: "Func1"},
					},
				},
				wantErr: false,
			},
			{
				name: "exclude-funcに/で始まるパスを指定",
				args: args{
					excludeFunc: helper.P("(/path/to/dir3).Func1"),
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "exclude-funcにパスのみ指定",
				args: args{
					excludeFunc: helper.P("(path/to/dir3)"),
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "exclude-funcにパス+構造体名のみ指定",
				args: args{
					excludeFunc: helper.P("(path/to/dir3.Struct1)"),
				},
				want:    nil,
				wantErr: true,
			},
			{
				name: "themeに期待値以外を設定",
				args: args{
					theme: helper.P("unknown"),
				},
				want:    nil,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				readerMock := &mockReader{
					mockExists: func(string) bool { return false },
				}

				got, err := option.New(readerMock).Generate(
					tt.args.input,
					tt.args.inputFiles,
					tt.args.output,
					tt.args.theme,
					tt.args.include,
					tt.args.exclude,
					tt.args.excludeFunc,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if d := cmp.Diff(got, tt.want); d != "" {
					t.Errorf(d)
				}
			})
		}
	})

	t.Run("設定ファイルが存在する", func(t *testing.T) {
		tests := testCase{
			{
				name: "全項目に設定値が存在(theme=dark)",
				settings: `
input: example.out
input-files:
  - example1.out
  - example2.out
output: example.html
theme: dark
include:
  - path/to/dir1
  - path/to/dir2
exclude:
  - path/to/dir3
  - path/to/dir4
exclude-func:
  - Func1
  - (path/to/dir3).Func1
  - (path/to/dir3.Struct1).Func1
`,
				args: args{},
				want: &option.Option{
					Input:      "example.out",
					InputFiles: []string{"example1.out", "example2.out"},
					Output:     "example.html",
					Theme:      "dark",
					Include:    []string{"path/to/dir1", "path/to/dir2"},
					Exclude:    []string{"path/to/dir3", "path/to/dir4"},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "", Struct: "", Func: "Func1"},
						{Path: "path/to/dir3", Struct: "", Func: "Func1"},
						{Path: "path/to/dir3", Struct: "Struct1", Func: "Func1"},
					},
				},
				wantErr: false,
			},
			{
				name: "全項目に設定値が存在(theme=light)",
				settings: `
input: example.out
input-files:
  - example1.out
  - example2.out
output: example.html
theme: light
include:
  - path/to/dir1
  - path/to/dir2
exclude:
  - path/to/dir3
  - path/to/dir4
exclude-func:
  - Func1
  - (path/to/dir3).Func1
  - (path/to/dir3.Struct1).Func1
`,
				args: args{},
				want: &option.Option{
					Input:      "example.out",
					InputFiles: []string{"example1.out", "example2.out"},
					Output:     "example.html",
					Theme:      "light",
					Include:    []string{"path/to/dir1", "path/to/dir2"},
					Exclude:    []string{"path/to/dir3", "path/to/dir4"},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "", Struct: "", Func: "Func1"},
						{Path: "path/to/dir3", Struct: "", Func: "Func1"},
						{Path: "path/to/dir3", Struct: "Struct1", Func: "Func1"},
					},
				},
				wantErr: false,
			},
			{
				name: "全項目に設定値が存在し、引数に全項目に設定値が存在",
				settings: `
input: example.out
input-files:
  - example1.out
  - example2.out
output: example.html
theme: dark
include:
  - path/to/dir1
  - path/to/dir2
exclude:
  - path/to/dir3
  - path/to/dir4
exclude-func:
  - Func1
  - (path/to/dir3).Func1
  - (path/to/dir3.Struct1).Func1
`,
				args: args{
					input:       helper.P("example2.out"),
					inputFiles:  helper.P("example1.out,example2.out"),
					output:      helper.P("example2.html"),
					theme:       helper.P("light"),
					include:     helper.P("path/to/dir5"),
					exclude:     helper.P("path/to/dir6"),
					excludeFunc: helper.P("Func2,Func3"),
				},
				want: &option.Option{
					Input:      "example2.out",
					InputFiles: []string{"example1.out", "example2.out"},
					Output:     "example2.html",
					Theme:      "light",
					Include:    []string{"path/to/dir5"},
					Exclude:    []string{"path/to/dir6"},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "", Struct: "", Func: "Func2"},
						{Path: "", Struct: "", Func: "Func3"},
					},
				},
				wantErr: false,
			},
			{
				name: "全項目に設定値が存在し、引数に全項目に空文字を設定",
				settings: `
input: example.out
input-files:
  - example1.out
  - example2.out
output: example.html
theme: light
include:
  - path/to/dir1
  - path/to/dir2
exclude:
  - path/to/dir3
  - path/to/dir4
exclude-func:
  - Func1
  - (path/to/dir3).Func1
  - (path/to/dir3.Struct1).Func1
`,
				args: args{
					input:       helper.P(""),
					inputFiles:  helper.P(""),
					output:      helper.P(""),
					theme:       helper.P(""),
					include:     helper.P(""),
					exclude:     helper.P(""),
					excludeFunc: helper.P(""),
				},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},

			{
				name: "全項目のキーのみが存在する",
				settings: `
input:
input-files:
output:
theme:
include:
exclude:
exclude-func:
`,
				args: args{},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "全項目のキーが存在しない",
				settings: `
# empty settings
`,
				args: args{},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "input-filesに空の値を持つ",
				settings: `
input-files:
  - example1.out
  -
  - example2.out
  -
  -
`,
				args: args{},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{"example1.out", "example2.out"},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},

			{
				name: "includeに空の値を持つ",
				settings: `
include:
  - path/to/dir1
  -
  - path/to/dir2
  -
  -
`,
				args: args{},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{"path/to/dir1", "path/to/dir2"},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "include./で始まるパスを指定",
				settings: `
include:
  - ./path/to/dir1
`,
				args: args{},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{"path/to/dir1"},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "includeに/で終わるパスを指定",
				settings: `
include:
  - path/to/dir1/
`,
				args: args{},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{"path/to/dir1"},
					Exclude:     []string{},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "includeに/で始まるパスを指定",
				settings: `
include:
  - /path/to/dir1
`,
				args:    args{},
				want:    nil,
				wantErr: true,
			},
			{
				name: "excludeに空の値を持つ",
				settings: `
exclude:
  - path/to/dir3
  -
  - path/to/dir4
  -
  -
`,
				args: args{},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{"path/to/dir3", "path/to/dir4"},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "excludeに./で始まるパスを指定",
				settings: `
exclude:
  - ./path/to/dir3
`,
				args: args{},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{"path/to/dir3"},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "excludeに/で終わるパスを指定",
				settings: `
exclude:
  - path/to/dir3/
`,
				args: args{},
				want: &option.Option{
					Input:       "coverage.out",
					InputFiles:  []string{},
					Output:      "coverage.html",
					Theme:       "dark",
					Include:     []string{},
					Exclude:     []string{"path/to/dir3"},
					ExcludeFunc: []option.ExcludeFuncOption{},
				},
				wantErr: false,
			},
			{
				name: "excludeに/で始まるパスを指定",
				settings: `
exclude:
  - /path/to/dir3
`,
				args:    args{},
				want:    nil,
				wantErr: true,
			},
			{
				name: "exclude-funcに空の値を持つ",
				settings: `
exclude-func:
  - Func1
  -
  - (path/to/dir3).Func1
  -
  - (path/to/dir3.Struct1).Func1
`,
				args: args{},
				want: &option.Option{
					Input:      "coverage.out",
					InputFiles: []string{},
					Output:     "coverage.html",
					Theme:      "dark",
					Include:    []string{},
					Exclude:    []string{},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "", Struct: "", Func: "Func1"},
						{Path: "path/to/dir3", Struct: "", Func: "Func1"},
						{Path: "path/to/dir3", Struct: "Struct1", Func: "Func1"},
					},
				},
				wantErr: false,
			},
			{
				name: "exclude-funcに./で始まるパスを指定",
				settings: `
exclude-func:
  - (./path/to/dir3).Func1
  - (./path/to/dir3.Struct1).Func1
`,
				args: args{},
				want: &option.Option{
					Input:      "coverage.out",
					InputFiles: []string{},
					Output:     "coverage.html",
					Theme:      "dark",
					Include:    []string{},
					Exclude:    []string{},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "path/to/dir3", Struct: "", Func: "Func1"},
						{Path: "path/to/dir3", Struct: "Struct1", Func: "Func1"},
					},
				},
				wantErr: false,
			},
			{
				name: "exclude-funcに/で終わるパスを指定",
				settings: `
exclude-func:
  - (path/to/dir3/).Func1
  - (path/to/dir3/.Struct1).Func1
`,
				args: args{},
				want: &option.Option{
					Input:      "coverage.out",
					InputFiles: []string{},
					Output:     "coverage.html",
					Theme:      "dark",
					Include:    []string{},
					Exclude:    []string{},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "path/to/dir3", Struct: "", Func: "Func1"},
						{Path: "path/to/dir3", Struct: "Struct1", Func: "Func1"},
					},
				},
				wantErr: false,
			},

			{
				name: "exclude-funcに/で始まるパスを指定",
				settings: `
exclude-func:
  - (/path/to/dir3).Func1
`,
				args:    args{},
				want:    nil,
				wantErr: true,
			},
			{
				name: "exclude-funcのパスにワイルドカードを指定",
				settings: `
exclude-func:
  - (*).Func1
  - (*.Struct1).Func1
`,
				args: args{},
				want: &option.Option{
					Input:      "coverage.out",
					InputFiles: []string{},
					Output:     "coverage.html",
					Theme:      "dark",
					Include:    []string{},
					Exclude:    []string{},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "", Struct: "", Func: "Func1"},
						{Path: "", Struct: "Struct1", Func: "Func1"},
					},
				},
				wantErr: false,
			},
			{
				name: "exclude-funcのパスにファイル名まで指定",
				settings: `
exclude-func:
  - (path/to/dir3/sample.go).Func1
  - (path/to/dir3/sample.go.Struct1).Func1
`,
				args: args{},
				want: &option.Option{
					Input:      "coverage.out",
					InputFiles: []string{},
					Output:     "coverage.html",
					Theme:      "dark",
					Include:    []string{},
					Exclude:    []string{},
					ExcludeFunc: []option.ExcludeFuncOption{
						{Path: "path/to/dir3/sample.go", Struct: "", Func: "Func1"},
						{Path: "path/to/dir3/sample.go", Struct: "Struct1", Func: "Func1"},
					},
				},
				wantErr: false,
			},
			{
				name: "exclude-funcにパスのみ指定",
				settings: `
exclude-func:
  - (/path/to/dir3)
`,
				args:    args{},
				want:    nil,
				wantErr: true,
			},
			{
				name: "exclude-funcにパス+構造体名のみ",
				settings: `
exclude-func:
  - (/path/to/dir3.Struct1)
`,
				args:    args{},
				want:    nil,
				wantErr: true,
			},
			{
				name: "themeに期待値以外を設定",
				settings: `
theme: unknown
`,
				args:    args{},
				want:    nil,
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				readerMock := &mockReader{
					mockExists: func(string) bool { return true },
					mockRead: func(string) (io.Reader, error) {
						return strings.NewReader(tt.settings), nil
					},
				}

				got, err := option.New(readerMock).Generate(
					tt.args.input,
					tt.args.inputFiles,
					tt.args.output,
					tt.args.theme,
					tt.args.include,
					tt.args.exclude,
					tt.args.excludeFunc,
				)
				if (err != nil) != tt.wantErr {
					t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if d := cmp.Diff(got, tt.want); d != "" {
					t.Errorf(d)
				}
			})
		}

	})
}
