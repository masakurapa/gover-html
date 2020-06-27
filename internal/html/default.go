package html

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/masakurapa/go-cover/internal/logger"
	"github.com/masakurapa/go-cover/internal/profile"
	"github.com/masakurapa/go-cover/internal/reader"
)

func Write(reader reader.Reader, out io.Writer, profiles profile.Profiles) error {
	logger.L.Debug("start generating html")

	files := make([]TemplateFile, 0, len(profiles))
	for _, p := range profiles {
		b, err := reader.Read(p.FileName)
		if err != nil {
			return fmt.Errorf("can't read %q: %v", p.FileName, err)
		}

		buf := writePool.Get().(*bytes.Buffer)
		genSource(buf, b, &p)

		files = append(files, TemplateFile{
			Name:     p.FileName,
			Body:     template.HTML(buf.String()),
			Coverage: p.Blocks.Coverage(),
		})

		buf.Reset()
		writePool.Put(buf)
	}

	logger.L.Debug("write html")
	return parsedDefaultTemplate.Execute(out, DefaultTemplateData{
		Files: files,
	})
}
