package profile

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/masakurapa/gover-html/internal/option"
)

var (
	newline = []byte("\n")
)

func Read(opt option.Option) (io.Reader, error) {
	// InputFilesが指定されていない場合だけInputを使う
	if len(opt.InputFiles) == 0 {
		b, err := os.ReadFile(opt.Input)
		if err != nil {
			return nil, err
		}
		return bytes.NewReader(b), nil
	}

	// 1行目は "mode: set" を固定にする
	buf := bytes.NewBuffer([]byte("mode: set"))
	buf.Write(newline)

	for _, in := range opt.InputFiles {
		if err := read(buf, in); err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func read(buf *bytes.Buffer, in string) error {
	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()

	// 1行目には "mode: xxx" が入っているはずなので無視して2行目から読み込む
	skip := true
	r := bufio.NewScanner(f)
	for r.Scan() {
		if skip {
			skip = false
			continue
		}
		buf.Write(r.Bytes())
		buf.Write(newline)
	}

	return nil
}
