package tree_test

import (
	"reflect"
	"testing"

	"github.com/masakurapa/gover-html/internal/html/tree"
	"github.com/masakurapa/gover-html/internal/profile"
)

func TestNode_ChildBlocks(t *testing.T) {
	node := tree.Node{
		Files: profile.Profiles{
			{Blocks: profile.Blocks{
				{StartLine: 1, StartCol: 2, EndLine: 3, EndCol: 4},
				{StartLine: 11, StartCol: 12, EndLine: 13, EndCol: 14},
			}},
		},
		Dirs: []tree.Node{
			{
				Files: profile.Profiles{
					{Blocks: profile.Blocks{
						{StartLine: 21, StartCol: 22, EndLine: 23, EndCol: 24},
						{StartLine: 31, StartCol: 32, EndLine: 33, EndCol: 34},
					}},
				},
				Dirs: []tree.Node{},
			},
			{
				Files: profile.Profiles{},
				Dirs: []tree.Node{
					{
						Files: profile.Profiles{
							{Blocks: profile.Blocks{
								{StartLine: 41, StartCol: 42, EndLine: 43, EndCol: 44},
								{StartLine: 51, StartCol: 52, EndLine: 53, EndCol: 54},
							}},
						},
						Dirs: []tree.Node{},
					},
				},
			},
		},
	}

	want := profile.Blocks{
		{StartLine: 1, StartCol: 2, EndLine: 3, EndCol: 4},
		{StartLine: 11, StartCol: 12, EndLine: 13, EndCol: 14},
		{StartLine: 21, StartCol: 22, EndLine: 23, EndCol: 24},
		{StartLine: 31, StartCol: 32, EndLine: 33, EndCol: 34},
		{StartLine: 41, StartCol: 42, EndLine: 43, EndCol: 44},
		{StartLine: 51, StartCol: 52, EndLine: 53, EndCol: 54},
	}

	if got := node.ChildBlocks(); !reflect.DeepEqual(got, want) {
		t.Errorf("ChildBlocks() = %v, want %v", got, want)
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		profiles profile.Profiles
	}
	tests := []struct {
		name string
		args args
		want []tree.Node
	}{
		{
			name: "ファイルが無いディレクトリはマージされ、ディレクトリごとに階層化されたスライスが返却される",
			args: args{
				profiles: profile.Profiles{
					{
						ModulePath: "github.com/masakurapa/gover-html",
						FileName:   "github.com/masakurapa/gover-html/path/to/dir1/file0.go",
					},
					{
						ModulePath: "github.com/masakurapa/gover-html",
						FileName:   "github.com/masakurapa/gover-html/path/to/dir1/file1.go",
					},
					{
						ModulePath: "github.com/masakurapa/gover-html",
						FileName:   "github.com/masakurapa/gover-html/path/to/dir2/file1.go",
					},
					{
						ModulePath: "github.com/masakurapa/gover-html",
						FileName:   "github.com/masakurapa/gover-html/path/to/dir3/sub/file1.go",
					},
					{
						ModulePath: "github.com/masakurapa/gover-html2",
						FileName:   "github.com/masakurapa/gover-html2/path/to/dir1/file1.go",
					},
				},
			},
			want: []tree.Node{
				{
					Name: "github.com/masakurapa/gover-html",
					Dirs: []tree.Node{
						{
							Name: "path/to",
							Dirs: []tree.Node{
								{
									Name: "dir1",
									Files: profile.Profiles{
										{
											ModulePath: "github.com/masakurapa/gover-html",
											FileName:   "github.com/masakurapa/gover-html/path/to/dir1/file0.go",
										},
										{
											ModulePath: "github.com/masakurapa/gover-html",
											FileName:   "github.com/masakurapa/gover-html/path/to/dir1/file1.go",
										},
									},
								},
								{
									Name: "dir2",
									Files: profile.Profiles{
										{
											ModulePath: "github.com/masakurapa/gover-html",
											FileName:   "github.com/masakurapa/gover-html/path/to/dir2/file1.go",
										},
									},
								},
								{
									Name: "dir3/sub",
									Files: profile.Profiles{
										{
											ModulePath: "github.com/masakurapa/gover-html",
											FileName:   "github.com/masakurapa/gover-html/path/to/dir3/sub/file1.go",
										},
									},
								},
							},
						},
					},
				},
				{
					Name: "github.com/masakurapa/gover-html2",
					Dirs: []tree.Node{
						{
							Name: "path/to/dir1",
							Files: profile.Profiles{
								{
									ModulePath: "github.com/masakurapa/gover-html2",
									FileName:   "github.com/masakurapa/gover-html2/path/to/dir1/file1.go",
								},
							},
						},
					},
				},
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
	profiles := profile.Profiles{
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
