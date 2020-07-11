package tree_test

import (
	"reflect"
	"testing"

	"github.com/masakurapa/go-cover/internal/html/tree"
	"github.com/masakurapa/go-cover/internal/profile"
)

func TestCreate(t *testing.T) {
	type args struct {
		profiles []profile.Profile
	}
	tests := []struct {
		name string
		args args
		want []tree.Node
	}{
		{
			name: "ファイルが無いディレクトリはマージされ、ディレクトリごとに階層化されたスライスが返却される",
			args: args{
				profiles: []profile.Profile{
					{FileName: "path/to/dir1/file0.go"},
					{FileName: "path/to/dir1/file1.go"},
					{FileName: "path/to/dir2/file1.go"},
					{FileName: "path/to/dir3/sub/file1.go"},
				},
			},
			want: []tree.Node{
				{Name: "path/to", Dirs: []tree.Node{
					{Name: "dir1", Files: []profile.Profile{
						{FileName: "path/to/dir1/file0.go"},
						{FileName: "path/to/dir1/file1.go"},
					}},
					{Name: "dir2", Files: []profile.Profile{
						{FileName: "path/to/dir2/file1.go"},
					}},
					{Name: "dir3/sub", Files: []profile.Profile{
						{FileName: "path/to/dir3/sub/file1.go"},
					}},
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tree.Create(tt.args.profiles); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCreate(b *testing.B) {
	profiles := []profile.Profile{
		{FileName: "path/to/dir1/file0.go"},
		{FileName: "path/to/dir1/file1.go"},
		{FileName: "path/to/dir2/file1.go"},
		{FileName: "path/to/dir3/sub/file1.go"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tree.Create(profiles)
	}
}
