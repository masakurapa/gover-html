package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/masakurapa/go-cover/internal/html"
	"github.com/masakurapa/go-cover/internal/profile"
)

var (
	input  = flag.String("i", "coverage.out", "coverage profile")
	output = flag.String("o", "coverage.html", "html file output")
)

func main() {
	parseFlags()

	f, err := os.Open(*input)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	profiles, err := profile.Read(f)
	if err != nil {
		panic(err)
	}

	// output
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
	fmt.Fprintln(os.Stderr, "Output coverage in HTML.")
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
