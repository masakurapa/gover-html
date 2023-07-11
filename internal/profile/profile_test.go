package profile_test

import (
	"testing"

	"github.com/masakurapa/gover-html/internal/profile"
)

func TestProfile_IsRelativeOrAbsolute(t *testing.T) {
	type fields struct {
		ID       int
		FileName string
		Blocks   []profile.Block
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "ファイル名の先頭がドットの場合、trueを返す",
			fields: fields{FileName: "./path/to/file.go"},
			want:   true,
		},
		{
			name:   "ファイル名が絶対パスの場合、trueを返す",
			fields: fields{FileName: "/path/to/file.go"},
			want:   true,
		},
		{
			name:   "ファイル名がドット始まり、絶対パス以外の場合、falseを返す",
			fields: fields{FileName: "github.com/path/to/package"},
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prof := &profile.Profile{
				ID:       tt.fields.ID,
				FileName: tt.fields.FileName,
				Blocks:   tt.fields.Blocks,
			}
			if got := prof.IsRelativeOrAbsolute(); got != tt.want {
				t.Errorf("Profile.IsRelativeOrAbsolute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProfile_FilePath(t *testing.T) {
	type fields struct {
		Dir      string
		FileName string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "FileNameが相対パスの場合、FileNameと同じ値を返す",
			fields: fields{Dir: "/path/to/dir", FileName: "./path/to/file.go"},
			want:   "./path/to/file.go",
		},
		{
			name:   "FileNameが絶対パスの場合、FileNameと同じ値を返す",
			fields: fields{Dir: "/path/to/dir", FileName: "/path/to/file.go"},
			want:   "/path/to/file.go",
		},
		{
			name:   "FileNameが相対パス、絶対パスの以外の場合、Path+ファイル名を連結した値を返す",
			fields: fields{Dir: "/path/to/dir", FileName: "path/to/file.go"},
			want:   "/path/to/dir/file.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prof := &profile.Profile{
				Dir:      tt.fields.Dir,
				FileName: tt.fields.FileName,
			}
			if got := prof.FilePath(); got != tt.want {
				t.Errorf("Profile.FilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlocks_Coverage(t *testing.T) {
	tests := []struct {
		name   string
		blocks *profile.Blocks
		want   float64
	}{
		{
			name: "全ブロックのカウントが1以上の場合、100を返す",
			blocks: &profile.Blocks{
				{NumState: 41, Count: 1},
				{NumState: 42, Count: 2},
			},
			want: 100,
		},
		{
			name: "カウントが0のブロックが存在する場合、小数点第二位を四捨五入した値を返す",
			blocks: &profile.Blocks{
				{NumState: 41, Count: 1},
				{NumState: 42, Count: 0},
			},
			want: 49.4, // 41 / (41 + 42) * 100 = 49.39759036144578
		},
		{
			name: "全ブロックのカウントが0の場合、0を返す",
			blocks: &profile.Blocks{
				{NumState: 41, Count: 0},
				{NumState: 42, Count: 0},
			},
			want: 0,
		},
		{
			name:   "ブロックが空の場合、0を返す",
			blocks: &profile.Blocks{},
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.blocks.Coverage(); got != tt.want {
				t.Errorf("Blocks.Coverage() = %v, want %v", got, tt.want)
			}
		})
	}
}
