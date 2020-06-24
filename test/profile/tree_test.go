package profile_test

import (
	"reflect"
	"testing"

	"github.com/masakurapa/go-cover/internal/profile"
)

func TestProfiles_ToTree(t *testing.T) {
	prof := profile.Profiles{
		profile.Profile{FileName: "path/to/dir1/file0.go"},
		profile.Profile{FileName: "path/to/dir1/file1.go"},
		profile.Profile{FileName: "path/to/dir2/file1.go"},
		profile.Profile{FileName: "path/to/dir3/sub/file1.go"},
	}

	tests := []struct {
		name string
		prof profile.Profiles
		want profile.Tree
	}{
		{
			name: "ファイルが無いディレクトリはマージされ、ディレクトリごとに階層化されたスライスが返却される",
			prof: prof,
			want: profile.Tree{
				{Name: "path/to", Profiles: profile.Profiles{}, SubDirs: profile.Tree{
					{Name: "dir1", Profiles: profile.Profiles{
						profile.Profile{FileName: "path/to/dir1/file0.go"},
						profile.Profile{FileName: "path/to/dir1/file1.go"},
					}, SubDirs: profile.Tree{}},
					{Name: "dir2", Profiles: profile.Profiles{
						profile.Profile{FileName: "path/to/dir2/file1.go"},
					}, SubDirs: profile.Tree{}},
					{Name: "dir3/sub", Profiles: profile.Profiles{
						profile.Profile{FileName: "path/to/dir3/sub/file1.go"},
					}, SubDirs: profile.Tree{}},
				}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.prof.ToTree(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Profiles.ToTree() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func BenchmarkProfiles_ToTree(b *testing.B) {
	prof := profile.Profiles{
		profile.Profile{FileName: "path/to/dir1/file0.go"},
		profile.Profile{FileName: "path/to/dir1/file1.go"},
		profile.Profile{FileName: "path/to/dir2/file1.go"},
		profile.Profile{FileName: "path/to/dir3/sub/file1.go"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		prof.ToTree()
	}
}
