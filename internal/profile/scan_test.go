package profile_test

import (
	"bufio"
	"reflect"
	"strings"
	"testing"

	"github.com/masakurapa/go-cover/internal/profile"
)

func TestScan(t *testing.T) {
	toScanner := func(s string) *bufio.Scanner {
		return bufio.NewScanner(strings.NewReader(s))
	}

	type args struct {
		s *bufio.Scanner
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
				s: toScanner(`mode: set
path/to/dir2/file1.go:1.11,21.31 41 51
path/to/dir1/file0.go:2.12,22.32 42 52
path/to/dir1/file0.go:3.13,23.33 43 53
path/to/dir1/file1.go:4.14,24.34 44 54
`),
			},
			want: profile.Profiles{
				profile.Profile{
					ID:       0,
					Mode:     "set",
					FileName: "path/to/dir1/file0.go",
					Blocks: []profile.Block{
						{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 52},
						{StartLine: 3, StartCol: 13, EndLine: 23, EndCol: 33, NumStmt: 43, Count: 53},
					},
				},
				profile.Profile{
					ID:       1,
					Mode:     "set",
					FileName: "path/to/dir1/file1.go",
					Blocks: []profile.Block{
						{StartLine: 4, StartCol: 14, EndLine: 24, EndCol: 34, NumStmt: 44, Count: 54},
					},
				},
				profile.Profile{
					ID:       2,
					Mode:     "set",
					FileName: "path/to/dir2/file1.go",
					Blocks: []profile.Block{
						{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 51},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "開始／終了の行とカラムが同じものは、count>1のデータのみにフィルタされること",
			args: args{
				s: toScanner(`mode: set
path/to/dir1/file0.go:2.12,22.32 42 0
path/to/dir1/file0.go:2.12,22.32 42 52
path/to/dir1/file0.go:3.13,23.33 43 53
`),
			},
			want: profile.Profiles{
				profile.Profile{
					ID:       0,
					Mode:     "set",
					FileName: "path/to/dir1/file0.go",
					Blocks: []profile.Block{
						{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 52},
						{StartLine: 3, StartCol: 13, EndLine: 23, EndCol: 33, NumStmt: 43, Count: 53},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "modeが無い場合はエラーが返却される",
			args: args{
				s: toScanner(`first line
path/to/dir1/file0.go:2.12,22.32 42 52
`),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "2行目以降のフォーマットが不正な場合はエラーが返却される",
			args: args{
				s: toScanner(`mode: set
path/to/dir1/file0.go,2.12,22.32 42 52
`),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := profile.Scan(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Scan() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scan() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
