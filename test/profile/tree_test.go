package profile_test

import (
	"reflect"
	"testing"

	"github.com/masakurapa/gocover-html/internal/profile"
	"golang.org/x/tools/cover"
)

func TestProfiles_ToTree(t *testing.T) {
	prof := profile.Profiles{
		cover.Profile{
			Mode:     "set",
			FileName: "path/to/dir1/file0.go",
			Blocks: []cover.ProfileBlock{
				{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 52},
				{StartLine: 3, StartCol: 13, EndLine: 23, EndCol: 33, NumStmt: 43, Count: 53},
			},
		},
		cover.Profile{
			Mode:     "set",
			FileName: "path/to/dir1/file1.go",
			Blocks: []cover.ProfileBlock{
				{StartLine: 4, StartCol: 14, EndLine: 24, EndCol: 34, NumStmt: 44, Count: 54},
			},
		},
		cover.Profile{
			Mode:     "set",
			FileName: "path/to/dir2/file1.go",
			Blocks: []cover.ProfileBlock{
				{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 51},
			},
		},
	}

	tests := []struct {
		name string
		prof profile.Profiles
		want profile.Tree
	}{
		{
			name: "ディレクトリごとに階層化されたスライスが返却される",
			prof: prof,
			want: profile.Tree{
				{Name: "path", Profiles: profile.Profiles{}, SubDirs: profile.Tree{
					{Name: "to", Profiles: profile.Profiles{}, SubDirs: profile.Tree{
						{Name: "dir1", Profiles: profile.Profiles{prof[0], prof[1]}, SubDirs: profile.Tree{}},
						{Name: "dir2", Profiles: profile.Profiles{prof[2]}, SubDirs: profile.Tree{}},
					}},
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
