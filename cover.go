package main

import (
	"flag"
	"fmt"
	"os"

	cover "github.com/masakurapa/go-cover/internal"
	"github.com/masakurapa/go-cover/internal/html"
)

var (
	input  = flag.String("i", "coverage.out", "coverage profile for input")
	output = flag.String("o", "coverage.html", "file for output")
)

func main() {
	parseFlags()

	f, err := os.Open(*input)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	profiles, err := cover.ReadProfile(f)
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
	fmt.Fprintln(os.Stderr, `Usage of 'go-cover':
'go-cover' requires coverage profle by 'go test':
	go test -coverprofile=coverage.out

Write out HTML file:
	go-cover

Specify input file and output file:
	go-cover -i c.out -o c.html`)

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
