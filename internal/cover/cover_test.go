package cover_test

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/masakurapa/gover-html/internal/cover"
	"github.com/masakurapa/gover-html/internal/cover/filter"
	"github.com/masakurapa/gover-html/internal/profile"
)

type mockFilter struct {
	filter.Filter
	mockIsOutputTarget     func(string, string) bool
	mockIsOutputTargetFunc func(string, string, string) bool
}

func (f *mockFilter) IsOutputTarget(path, fileName string) bool {
	return f.mockIsOutputTarget(path, fileName)
}

func (f *mockFilter) IsOutputTargetFunc(s1, s2, s3 string) bool {
	return f.mockIsOutputTargetFunc(s1, s2, s3)
}

func TestReadProfile(t *testing.T) {
	absDir1, err := filepath.Abs("../../testdata/dir1")
	if err != nil {
		t.Fatal(err)
	}
	absDir2, err := filepath.Abs("../../testdata/dir2")
	if err != nil {
		t.Fatal(err)
	}
	absDir3, err := filepath.Abs("../../testdata/dir2/dir3")
	if err != nil {
		t.Fatal(err)
	}

	defaultFilter := &mockFilter{
		mockIsOutputTarget: func(string, string) bool {
			return true
		},
		mockIsOutputTargetFunc: func(s1 string, s2 string, s3 string) bool {
			return true
		},
	}

	type args struct {
		r io.Reader
		f filter.Filter
	}
	tests := []struct {
		name    string
		args    args
		want    profile.Profiles
		wantErr bool
	}{
		{
			name: "ファイル名の昇順でProfileが返却される",
			args: args{
				r: strings.NewReader(`mode: set
github.com/masakurapa/gover-html/testdata/dir2/file1.go:1.11,21.31 41 51
github.com/masakurapa/gover-html/testdata/dir1/file0.go:2.12,22.32 42 52
github.com/masakurapa/gover-html/testdata/dir1/file0.go:3.13,23.33 43 53
github.com/masakurapa/gover-html/testdata/dir1/file1.go:4.14,24.34 44 54
github.com/masakurapa/gover-html/testdata/dir2/dir3/file3.go:5.15,25.35 45 55
`),
				f: defaultFilter,
			},
			want: profile.Profiles{
				{
					ID:         2,
					ModulePath: "github.com/masakurapa/gover-html",
					Dir:        absDir1,
					FileName:   "github.com/masakurapa/gover-html/testdata/dir1/file0.go",
					Blocks: []profile.Block{
						{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumState: 42, Count: 52},
						{StartLine: 3, StartCol: 13, EndLine: 23, EndCol: 33, NumState: 43, Count: 53},
					},
					Functions: profile.Functions{
						{Name: "func Func1() string", StartLine: 3, StartCol: 1},
						{Name: "func Func2(s string) int", StartLine: 7, StartCol: 1},
					},
				},
				{
					ID:         3,
					ModulePath: "github.com/masakurapa/gover-html",
					Dir:        absDir1,
					FileName:   "github.com/masakurapa/gover-html/testdata/dir1/file1.go",
					Blocks: []profile.Block{
						{StartLine: 4, StartCol: 14, EndLine: 24, EndCol: 34, NumState: 44, Count: 54},
					},
					Functions: profile.Functions{
						{Name: "func Func3()", StartLine: 7, StartCol: 1},
						{Name: "func Func4(s string) bool", StartLine: 25, StartCol: 1},
					},
				},
				{
					ID:         4,
					ModulePath: "github.com/masakurapa/gover-html",
					Dir:        absDir3,
					FileName:   "github.com/masakurapa/gover-html/testdata/dir2/dir3/file3.go",
					Blocks: []profile.Block{
						{StartLine: 5, StartCol: 15, EndLine: 25, EndCol: 35, NumState: 45, Count: 55},
					},
					Functions: profile.Functions{
						{Name: "func (s *Struct2) Func1()", StartLine: 9, StartCol: 1},
						{Name: "func (s *Struct2) Func4(ss string) bool", StartLine: 27, StartCol: 1},
					},
				},
				{
					ID:         1,
					ModulePath: "github.com/masakurapa/gover-html",
					Dir:        absDir2,
					FileName:   "github.com/masakurapa/gover-html/testdata/dir2/file1.go",
					Blocks: []profile.Block{
						{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumState: 41, Count: 51},
					},
					Functions: profile.Functions{
						{Name: "func (*Struct1) Func1() string", StartLine: 5, StartCol: 1},
						{Name: "func (*Struct1) Func2(s string) int", StartLine: 9, StartCol: 1},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "開始／終了の行とカラムが同じものは、count>1のデータのみにフィルタされること",
			args: args{
				r: strings.NewReader(`mode: set
github.com/masakurapa/gover-html/testdata/dir1/file0.go:2.12,22.32 42 0
github.com/masakurapa/gover-html/testdata/dir1/file0.go:2.12,22.32 42 52
github.com/masakurapa/gover-html/testdata/dir1/file0.go:3.13,23.33 43 53
`),
				f: defaultFilter,
			},
			want: profile.Profiles{
				{
					ID:         1,
					ModulePath: "github.com/masakurapa/gover-html",
					Dir:        absDir1,
					FileName:   "github.com/masakurapa/gover-html/testdata/dir1/file0.go",
					Blocks: []profile.Block{
						{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumState: 42, Count: 52},
						{StartLine: 3, StartCol: 13, EndLine: 23, EndCol: 33, NumState: 43, Count: 53},
					},
					Functions: profile.Functions{
						{Name: "func Func1() string", StartLine: 3, StartCol: 1},
						{Name: "func Func2(s string) int", StartLine: 7, StartCol: 1},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "フィルタでtrueを返すファイルのProfileのみ返却される",
			args: args{
				r: strings.NewReader(`mode: set
github.com/masakurapa/gover-html/testdata/dir2/file1.go:1.11,21.31 41 51
github.com/masakurapa/gover-html/testdata/dir1/file0.go:2.12,22.32 42 52
github.com/masakurapa/gover-html/testdata/dir1/file0.go:3.13,23.33 43 53
github.com/masakurapa/gover-html/testdata/dir1/file1.go:4.14,24.34 44 54
github.com/masakurapa/gover-html/testdata/dir2/dir3/file3.go:5.15,25.35 45 55
`),
				f: &mockFilter{
					mockIsOutputTarget: func(path, fileName string) bool {
						if path == "testdata/dir1" && fileName == "file0.go" {
							return false
						}
						if path == "testdata/dir1" && fileName == "file1.go" {
							return false
						}
						if path == "testdata/dir2" && fileName == "file1.go" {
							return true
						}
						if path == "testdata/dir2/dir3" && fileName == "file3.go" {
							return true
						}

						t.Errorf("Unexpected parameters. path: %q, fileName: %q", path, fileName)
						return false
					},
					mockIsOutputTargetFunc: func(s string, s2 string, s3 string) bool {
						return true
					},
				},
			},
			want: profile.Profiles{
				{
					ID:         4,
					ModulePath: "github.com/masakurapa/gover-html",
					Dir:        absDir3,
					FileName:   "github.com/masakurapa/gover-html/testdata/dir2/dir3/file3.go",
					Blocks: []profile.Block{
						{StartLine: 5, StartCol: 15, EndLine: 25, EndCol: 35, NumState: 45, Count: 55},
					},
					Functions: profile.Functions{
						{Name: "func (s *Struct2) Func1()", StartLine: 9, StartCol: 1},
						{Name: "func (s *Struct2) Func4(ss string) bool", StartLine: 27, StartCol: 1},
					},
				},
				{
					ID:         1,
					ModulePath: "github.com/masakurapa/gover-html",
					Dir:        absDir2,
					FileName:   "github.com/masakurapa/gover-html/testdata/dir2/file1.go",
					Blocks: []profile.Block{
						{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumState: 41, Count: 51},
					},
					Functions: profile.Functions{
						{Name: "func (*Struct1) Func1() string", StartLine: 5, StartCol: 1},
						{Name: "func (*Struct1) Func2(s string) int", StartLine: 9, StartCol: 1},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "modeが無い場合はエラーが返却される",
			args: args{
				r: strings.NewReader(`first line
github.com/masakurapa/gover-html/testdata/dir1/file0.go:2.12,22.32 42 52
`),
				f: defaultFilter,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "2行目以降のフォーマットが不正な場合はエラーが返却される",
			args: args{
				r: strings.NewReader(`mode: set
github.com/masakurapa/gover-html/testdata/dir1/file0.go,2.12,22.32 42 52
`),
				f: defaultFilter,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cover.ReadProfile(tt.args.r, tt.args.f)
			_, _ = got, err
			if (err != nil) != tt.wantErr {
				t.Fatalf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func BenchmarkReadProfile(b *testing.B) {
	r := strings.NewReader(`mode: set
github.com/masakurapa/gover-html/testdata/dir2/file1.go:1.11,21.31 41 51
github.com/masakurapa/gover-html/testdata/dir1/file0.go:2.12,22.32 42 52
github.com/masakurapa/gover-html/testdata/dir1/file0.go:3.13,23.33 43 53
github.com/masakurapa/gover-html/testdata/dir1/file1.go:4.14,24.34 44 54
`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cover.ReadProfile(r, &mockFilter{
			mockIsOutputTarget: func(string, string) bool {
				return true
			},
		})
	}
}
