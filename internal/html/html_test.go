package html_test

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/masakurapa/go-cover/internal/html"
	"github.com/masakurapa/go-cover/internal/profile"
)

func TestWriteTreeView(t *testing.T) {
	dir, err := filepath.Abs("../../testdata/ex.go")
	if err != nil {
		t.Fatal(err)
	}

	var writer bytes.Buffer
	profiles := []profile.Profile{
		{
			ID:       1,
			Dir:      dir,
			FileName: dir,
			Blocks: []profile.Block{
				{StartLine: 5, StartCol: 31, EndLine: 7, EndCol: 16, NumState: 2, Count: 1},
				{StartLine: 11, StartCol: 2, EndLine: 11, EndCol: 13, NumState: 1, Count: 1},
				{StartLine: 17, StartCol: 2, EndLine: 17, EndCol: 17, NumState: 1, Count: 0},
				{StartLine: 7, StartCol: 16, EndLine: 9, EndCol: 3, NumState: 1, Count: 1},
				{StartLine: 11, StartCol: 13, EndLine: 13, EndCol: 3, NumState: 1, Count: 1},
				{StartLine: 13, StartCol: 8, EndLine: 13, EndCol: 20, NumState: 1, Count: 0},
				{StartLine: 13, StartCol: 20, EndLine: 15, EndCol: 3, NumState: 1, Count: 0},
			},
		},
	}

	err = html.WriteTreeView(&writer, profiles)
	if err != nil {
		t.Fatal(err)
	}

	// got := writer.String()
	// if got != want {
	// 	t.Errorf("WriteTreeView() = %v, want %v", got, want)
	// }
}

func BenchmarkWriteTreeView(b *testing.B) {
	dir, err := filepath.Abs("../../testdata/ex.go")
	if err != nil {
		b.Fatal(err)
	}

	var writer bytes.Buffer
	profiles := []profile.Profile{
		{
			ID:       1,
			Dir:      dir,
			FileName: dir,
			Blocks: []profile.Block{
				{StartLine: 5, StartCol: 31, EndLine: 7, EndCol: 16, NumState: 2, Count: 1},
				{StartLine: 11, StartCol: 2, EndLine: 11, EndCol: 13, NumState: 1, Count: 1},
				{StartLine: 17, StartCol: 2, EndLine: 17, EndCol: 17, NumState: 1, Count: 0},
				{StartLine: 7, StartCol: 16, EndLine: 9, EndCol: 3, NumState: 1, Count: 1},
				{StartLine: 11, StartCol: 13, EndLine: 13, EndCol: 3, NumState: 1, Count: 1},
				{StartLine: 13, StartCol: 8, EndLine: 13, EndCol: 20, NumState: 1, Count: 0},
				{StartLine: 13, StartCol: 20, EndLine: 15, EndCol: 3, NumState: 1, Count: 0},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		html.WriteTreeView(&writer, profiles)
		writer.Reset()
	}
}
