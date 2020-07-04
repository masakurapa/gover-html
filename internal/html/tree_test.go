package html_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/masakurapa/go-cover/internal/html"
	"github.com/masakurapa/go-cover/internal/profile"
)

func TestWriteTreeView(t *testing.T) {
	reader := stubReader{}
	var writer bytes.Buffer
	profiles := profile.Profiles{
		{
			ID:       1,
			FileName: "path/to/example.go",
			Mode:     "set",
			Blocks: profile.Blocks{
				{StartLine: 5, StartCol: 31, EndLine: 7, EndCol: 16, NumStmt: 2, Count: 1},
				{StartLine: 11, StartCol: 2, EndLine: 11, EndCol: 13, NumStmt: 1, Count: 1},
				{StartLine: 17, StartCol: 2, EndLine: 17, EndCol: 17, NumStmt: 1, Count: 0},
				{StartLine: 7, StartCol: 16, EndLine: 9, EndCol: 3, NumStmt: 1, Count: 1},
				{StartLine: 11, StartCol: 13, EndLine: 13, EndCol: 3, NumStmt: 1, Count: 1},
				{StartLine: 13, StartCol: 8, EndLine: 13, EndCol: 20, NumStmt: 1, Count: 0},
				{StartLine: 13, StartCol: 20, EndLine: 15, EndCol: 3, NumStmt: 1, Count: 0},
			},
		},
	}
	tree := profiles.ToTree()

	err := html.WriteTreeView(&reader, &writer, profiles, tree)
	if err != nil {
		t.Fatal(err)
	}

	// got := writer.String()
	// if got != want {
	// 	t.Errorf("WriteTreeView() = %v, want %v", got, want)
	// }
}

// stub
type stubReader struct{}

func (stub *stubReader) Read(s string) ([]byte, error) {
	if s != "path/to/example.go" {
		return nil, fmt.Errorf("file not found")
	}

	return []byte(stubSrc), nil
}

const stubSrc = `package example

import "strconv"

func Example(s string) string {
	n, err := strconv.Atoi(s)
	if err != nil {
		return "error!!"
	}

	if n <= 10 {
		return "hello"
	} else if n <= 20 {
		return "world"
	}

	return "ninja!"
}
`

const want = ``
