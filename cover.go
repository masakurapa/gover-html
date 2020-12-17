package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/masakurapa/gover-html/internal/cover"
	"github.com/masakurapa/gover-html/internal/html"
)

var (
	input  = flag.String("i", "coverage.out", "coverage profile for input")
	output = flag.String("o", "coverage.html", "file for output")
	filter = flag.String("filter", "", `output only the specified directories.
multiple directories can be specified separated by commas.`)
)

func main() {
	parseFlags()

	f, err := os.Open(*input)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	profiles, err := cover.ReadProfile(f, getFilters())
	if err != nil {
		panic(err)
	}

	out, err := os.Create(*output)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if err = html.WriteTreeView(out, profiles); err != nil {
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

func getFilters() []string {
	if filter == nil || *filter == "" {
		return []string{}
	}
	return strings.Split(*filter, ",")
}
