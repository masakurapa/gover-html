package profile_test

import (
	"reflect"
	"testing"

	"github.com/masakurapa/go-cover/internal/profile"
)

func TestBlocks_Filter(t *testing.T) {
	blocks := profile.Blocks{
		{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 0},
		{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 51},
		{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 0},
	}

	// "開始／終了の行と位置が同じものは、count>1のデータのみにフィルタされること",
	want := profile.Blocks{
		{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 51},
		{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 0},
	}

	if got := blocks.Filter(); !reflect.DeepEqual(got, want) {
		t.Errorf("Blocks.Filter() = %v, want %v", got, want)
	}
}

func TestBlocks_Sort(t *testing.T) {
	blocks := profile.Blocks{
		{StartLine: 2, StartCol: 23, EndLine: 3, EndCol: 30},
		{StartLine: 2, StartCol: 12, EndLine: 2, EndCol: 22},
		{StartLine: 1, StartCol: 21, EndLine: 2, EndCol: 11},
		{StartLine: 1, StartCol: 11, EndLine: 1, EndCol: 20},
	}

	// 開始行 -> 開始位置の昇順にソートされること
	want := profile.Blocks{
		{StartLine: 1, StartCol: 11, EndLine: 1, EndCol: 20},
		{StartLine: 1, StartCol: 21, EndLine: 2, EndCol: 11},
		{StartLine: 2, StartCol: 12, EndLine: 2, EndCol: 22},
		{StartLine: 2, StartCol: 23, EndLine: 3, EndCol: 30},
	}

	blocks.Sort()
	if !reflect.DeepEqual(blocks, want) {
		t.Errorf("Profiles.Sort() = %#v, want %#v", blocks, want)
	}
}

func TestBlocks_Coverage(t *testing.T) {
	tests := []struct {
		name   string
		blocks profile.Blocks
		want   float64
	}{
		{
			name: "全レコードのcount > 1の場合100が返却される",
			blocks: profile.Blocks{
				{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 1},
				{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 2},
			},
			want: 100,
		},
		{
			name: "一部レコードがcount > 1の場合小数点を含む値が返却される",
			blocks: profile.Blocks{
				{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 1},
				{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 0},
			},
			want: 49.39759036144578, // 41 / (41 + 42) * 100
		},
		{
			name: "count > 1のレコードが無い場合0が返却される",
			blocks: profile.Blocks{
				{StartLine: 1, StartCol: 11, EndLine: 21, EndCol: 31, NumStmt: 41, Count: 0},
				{StartLine: 2, StartCol: 12, EndLine: 22, EndCol: 32, NumStmt: 42, Count: 0},
			},
			want: 0,
		},
		{
			name:   "Blocksが空の場合0が返却される",
			blocks: profile.Blocks{},
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
