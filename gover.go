package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/masakurapa/gover-html/internal/cover"
	"github.com/masakurapa/gover-html/internal/cover/filter"
	"github.com/masakurapa/gover-html/internal/html"
	"github.com/masakurapa/gover-html/internal/option"
	"github.com/masakurapa/gover-html/internal/profile"
	"github.com/masakurapa/gover-html/internal/reader"
)

var (
	input      = flag.String("i", "coverage.out", "coverage profile for input")
	inputFiles = flag.String("input-files", "", `Specify multiple input coverage profiles.
multiple files can be specified separated by commas.`)
	output = flag.String("o", "coverage.html", "file for output")
	theme  = flag.String("theme", "dark", `HTML color theme to output ('dark' or 'light')
if the value is invalid, it will be 'dark'.`)
	include = flag.String("include", "", `output only the specified directories.
multiple directories can be specified separated by commas.

if "exclude" is also specified, "exclude" option takes precedence.`)
	exclude = flag.String("exclude", "", `output expect the specified directories.
multiple directories can be specified separated by commas.

if "include" is also specified, this option takes precedence.`)
	excludeFunc = flag.String("exclude-func", "", `exclude specified function from output.
multiple directories can be specified separated by commas.`)
)

func main() {
	opt := getOption()

	r, err := profile.Read(opt)
	if err != nil {
		exitError(err)
	}

	profiles, err := cover.ReadProfile(r, filter.New(opt))
	if err != nil {
		exitError(err)
	}

	out, err := os.Create(opt.Output)
	if err != nil {
		exitError(err)
	}
	defer out.Close()

	if err = html.WriteTreeView(out, profiles, opt); err != nil {
		exitError(err)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `Usage of 'gover-html':
  'gover-html' requires coverage profle by 'go test':
        go test -coverprofile=coverage.out

  Write out HTML file:
        gover-html

  Specify input file and output file:
        gover-html -i c.out -o c.html`)

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Flags:")
	flag.PrintDefaults()
	os.Exit(2)
}

func parseFlags() {
	flag.Usage = usage
	flag.Parse()
}

func getOption() option.Option {
	parseFlags()

	// make options with command line arguments
	opt, err := option.New(reader.New()).Generate(input, inputFiles, output, theme, include, exclude, excludeFunc)
	if err != nil {
		exitError(err)
	}

	return *opt
}

func exitError(err error) {
	fmt.Println(err)
	os.Exit(1)
}
