package cover

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/masakurapa/gover-html/internal/cover/filter"
	"github.com/masakurapa/gover-html/internal/profile"
)

var reg = regexp.MustCompile(`^(.+):([0-9]+)\.([0-9]+),([0-9]+)\.([0-9]+) ([0-9]+) ([0-9]+)$`)

type importDir struct {
	modulePath string
	relative   string
	dir        string
}

// ReadProfile is reads profiling data
func ReadProfile(r io.Reader, f filter.Filter) (profile.Profiles, error) {
	files := make(map[string]*profile.Profile)
	modeOk := false
	id := 1

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()

		// first line must be "mode: xxx"
		if !modeOk {
			const p = "mode: "
			if !strings.HasPrefix(line, p) || line == p {
				return nil, fmt.Errorf("first line must be mode: %q", line)
			}
			modeOk = true
			continue
		}

		matches := reg.FindStringSubmatch(line)
		if matches == nil {
			return nil, fmt.Errorf("%q does not match expected format: %v", line, reg)
		}

		fileName := matches[1]
		p := files[fileName]

		if p == nil {
			p = &profile.Profile{ID: id, FileName: fileName}
			files[fileName] = p
			id++
		}

		p.Blocks = append(p.Blocks, profile.Block{
			StartLine: toInt(matches[2]),
			StartCol:  toInt(matches[3]),
			EndLine:   toInt(matches[4]),
			EndCol:    toInt(matches[5]),
			NumState:  toInt(matches[6]),
			Count:     toInt(matches[7]),
		})
	}

	return toProfiles(files, f)
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func toProfiles(files map[string]*profile.Profile, f filter.Filter) (profile.Profiles, error) {
	dirs, err := makeImportDirMap(files)
	if err != nil {
		return nil, err
	}

	profiles := make(profile.Profiles, 0, len(files))
	for _, p := range files {
		p.Blocks = filterBlocks(p.Blocks)
		sort.SliceStable(p.Blocks, func(i, j int) bool {
			bi, bj := p.Blocks[i], p.Blocks[j]
			return bi.StartLine < bj.StartLine || bi.StartLine == bj.StartLine && bi.StartCol < bj.StartCol
		})

		d := dirs[path.Dir(p.FileName)]
		if !f.IsOutputTarget(d.relative, filepath.Base(p.FileName)) {
			continue
		}

		p.Dir = d.dir
		p.ModulePath = d.modulePath

		pp, err := makeNewProfile(p, f)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, *pp)
	}

	sort.SliceStable(profiles, func(i, j int) bool {
		return profiles[i].FileName < profiles[j].FileName
	})

	return profiles, nil
}

func filterBlocks(blocks []profile.Block) []profile.Block {
	index := func(bs []profile.Block, b *profile.Block) int {
		for i, bb := range bs {
			if bb.StartLine == b.StartLine &&
				bb.StartCol == b.StartCol &&
				bb.EndLine == b.EndLine &&
				bb.EndCol == b.EndCol {
				return i
			}
		}
		return -1
	}

	newBlocks := make([]profile.Block, 0, len(blocks))
	for _, b := range blocks {
		i := index(newBlocks, &b)
		if i == -1 {
			newBlocks = append(newBlocks, b)
			continue
		}
		if b.Count > 0 {
			newBlocks[i] = b
		}
	}

	return newBlocks
}

func makeImportDirMap(files map[string]*profile.Profile) (map[string]importDir, error) {
	stdout, err := execGoList(files)
	if err != nil {
		return nil, err
	}

	pkgs := make(map[string]importDir)
	if len(stdout) == 0 {
		return pkgs, nil
	}

	type pkg struct {
		Dir    string
		Module *struct {
			Path string
			Dir  string
		}
		ImportPath string
		Error      *struct {
			Err string
		}
	}

	dec := json.NewDecoder(bytes.NewReader(stdout))
	for {
		var p pkg
		err := dec.Decode(&p)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("decoding go list json: %v", err)
		}
		if p.Error != nil {
			return nil, fmt.Errorf(p.Error.Err)
		}
		// should have the same result for "pkg.ImportPath" and "path.Dir(Profile.FileName)"
		pkgs[p.ImportPath] = importDir{
			modulePath: p.Module.Path,
			relative:   strings.TrimPrefix(p.Dir, p.Module.Dir+"/"),
			dir:        p.Dir,
		}
	}

	return pkgs, nil
}

// execute "go list" command
func execGoList(files map[string]*profile.Profile) ([]byte, error) {
	dirs := make([]string, 0, len(files))
	m := make(map[string]struct{})

	for _, p := range files {
		if p.IsRelativeOrAbsolute() {
			continue
		}
		dir := path.Dir(p.FileName)
		if _, ok := m[dir]; !ok {
			m[dir] = struct{}{}
			dirs = append(dirs, dir)
		}
	}

	if len(dirs) == 0 {
		return make([]byte, 0), nil
	}

	cmdName := filepath.Join(runtime.GOROOT(), "bin/go")
	args := append([]string{"list", "-e", "-json"}, dirs...)
	cmd := exec.Command(cmdName, args...)

	return cmd.Output()
}
