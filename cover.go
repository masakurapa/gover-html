package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/masakurapa/gover-html/internal/cover"
	"github.com/masakurapa/gover-html/internal/cover/filter"
	"github.com/masakurapa/gover-html/internal/html"
)

var (
	input  = flag.String("i", "coverage.out", "coverage profile for input")
	output = flag.String("o", "coverage.html", "file for output")
	theme  = flag.String("theme", "dark", `HTML color theme to output ('dark' or 'light')
if the value is invalid, it will be 'dark'.
`)
	include = flag.String("include", "", `output only the specified directories.
multiple directories can be specified separated by commas.`)
)

func main() {
	parseFlags()

	f, err := os.Open(*input)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	profiles, err := cover.ReadProfile(f, filter.New(include))
	if err != nil {
		panic(err)
	}

	out, err := os.Create(*output)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if err = html.WriteTreeView(out, getTheme(), profiles); err != nil {
		panic(err)
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

	if input == nil || *input == "" {
		flag.Usage()
	}
	if output == nil || *output == "" {
		flag.Usage()
	}
}

func getTheme() string {
	if theme == nil || *theme == "" {
		return "dark"
	}
	if *theme != "dark" && *theme != "light" {
		return "dark"
	}
	return *theme
}
