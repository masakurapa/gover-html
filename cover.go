package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/masakurapa/gocover-html/internal/html"
	"github.com/masakurapa/gocover-html/internal/profile"
)

func main() {
	path := "./coverage.out"
	profiles := read(path)
	html.Print(profiles, profiles.ToTree())
}

func read(path string) profile.Profiles {
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintf("failed to open %q: %s", path, err))
	}
	defer f.Close()

	buf := bufio.NewReader(f)

	profiles, err := profile.Scan(bufio.NewScanner(buf))
	if err != nil {
		panic(err)
	}

	return profiles
}
