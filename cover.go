package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/masakurapa/go-cover/internal/html"
	"github.com/masakurapa/go-cover/internal/logger"
	"github.com/masakurapa/go-cover/internal/profile"
	"github.com/masakurapa/go-cover/internal/reader"
)

var (
	input   = flag.String("i", "coverage.out", "coverage profile")
	output  = flag.String("o", "coverage.html", "html file output")
	isTree  = flag.Bool("tree", false, "output tree view")
	verbose = flag.Bool("v", false, "verbose output log")
)

func main() {
	start := time.Now()
	parseFlags()

	logger.New(*verbose)
	logger.L.Debug("start go-cover")

	logger.L.Debug("open profile: %s", *input)
	f, err := os.Open(*input)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}
	logger.L.Debug("profile size: %dbytes", fi.Size())

	buf := bufio.NewReader(f)
	profiles, err := profile.Scan(bufio.NewScanner(buf))
	if err != nil {
		panic(err)
	}

	// output
	out, err := os.Create(*output)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if *isTree {
		if err = html.WriteTreeView2(reader.New(), out, profiles, profiles.ToTree()); err != nil {
			panic(err)
		}
	} else {
		if err = html.Write(reader.New(), out, profiles); err != nil {
			panic(err)
		}
	}

	sec := time.Now().Sub(start).Milliseconds()
	logger.L.Debug("end go-cover. total: %dms", sec)
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
