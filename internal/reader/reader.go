package reader

import (
	"io"
	"os"
)

// Reader is file reader
type Reader interface {
	Read(string) (io.Reader, error)
	Exists(string) bool
}

type fileReader struct{}

// New is initialize the file reader
func New() Reader {
	return &fileReader{}
}

func (r *fileReader) Read(file string) (io.Reader, error) {
	return os.Open(file)
}

func (r *fileReader) Exists(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}
