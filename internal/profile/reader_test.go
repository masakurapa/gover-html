package profile_test

import (
	"io"
	"testing"

	"github.com/masakurapa/gover-html/internal/option"
	"github.com/masakurapa/gover-html/internal/profile"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name string
		opt  option.Option
		want string
	}{
		{
			name: "InputFilesを指定しない場合は1ファイルのみ読み込んだ結果が返却される",
			opt: option.Option{
				Input:      "../../testdata/coverage_test1.out",
				InputFiles: []string{},
			},
			want: `mode: set
github.com/masakurapa/gover-html/internal/html/html.go:54.87,66.29 6 1
github.com/masakurapa/gover-html/internal/html/html.go:66.29,68.17 2 1
github.com/masakurapa/gover-html/internal/html/html.go:68.17,70.4 1 0
github.com/masakurapa/gover-html/internal/html/html.go:72.3,82.34 4 1
github.com/masakurapa/gover-html/internal/html/html.go:82.34,88.4 1 0
`,
		},
		{
			name: "InputFilesを指定した場合は複数ファイルがマージされた結果が返却される",
			opt: option.Option{
				Input: "../../testdata/coverage_test1.out",
				InputFiles: []string{
					"../../testdata/coverage_test1.out",
					"../../testdata/coverage_test2.out",
					"../../testdata/coverage_test3.out",
				},
			},
			want: `mode: set
github.com/masakurapa/gover-html/internal/html/html.go:54.87,66.29 6 1
github.com/masakurapa/gover-html/internal/html/html.go:66.29,68.17 2 1
github.com/masakurapa/gover-html/internal/html/html.go:68.17,70.4 1 0
github.com/masakurapa/gover-html/internal/html/html.go:72.3,82.34 4 1
github.com/masakurapa/gover-html/internal/html/html.go:82.34,88.4 1 0
github.com/masakurapa/gover-html/internal/html/tree/tree.go:18.45,20.28 2 1
github.com/masakurapa/gover-html/internal/html/tree/tree.go:20.28,22.3 1 1
github.com/masakurapa/gover-html/internal/html/tree/tree.go:23.2,23.28 1 1
github.com/masakurapa/gover-html/internal/cover/func.go:31.56,32.26 1 1
github.com/masakurapa/gover-html/internal/cover/func.go:33.21,34.20 1 1
github.com/masakurapa/gover-html/internal/cover/func.go:34.20,35.9 1 0
github.com/masakurapa/gover-html/internal/cover/func.go:37.3,46.17 5 1
github.com/masakurapa/gover-html/internal/cover/func.go:46.17,47.14 1 0
github.com/masakurapa/gover-html/internal/cover/func.go:49.3,50.17 2 1
github.com/masakurapa/gover-html/internal/cover/func.go:50.17,51.14 1 0
github.com/masakurapa/gover-html/internal/cover/func.go:54.3,55.44 2 1
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := profile.Read(tt.opt)
			if err != nil {
				t.Fatal(err)
			}
			body, err := io.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			if got := string(body); got != tt.want {
				t.Errorf("profile.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
