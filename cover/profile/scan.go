package profile

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var reg = regexp.MustCompile(`^(.+):([0-9]+).([0-9]+),([0-9]+).([0-9]+) ([0-9]+) ([0-9]+)$`)

func Scan(s *bufio.Scanner) (Profiles, error) {
	files := make(map[string]*Profile)
	mode := ""

	for s.Scan() {
		line := s.Text()

		if mode == "" {
			const p = "mode: "
			if !strings.HasPrefix(line, p) || line == p {
				return nil, fmt.Errorf("bad mode line: %q", line)
			}
			mode = line[len(p):]
			continue
		}

		m := reg.FindStringSubmatch(line)
		if m == nil {
			return nil, fmt.Errorf("line %q does not match expected format: %v", line, reg)
		}

		fileName := m[1]
		p := files[fileName]

		if p == nil {
			p = &Profile{FileName: fileName, Mode: mode}
			files[fileName] = p
		}

		p.Blocks = append(p.Blocks, Block{
			StartLine: toInt(m[2]),
			StartCol:  toInt(m[3]),
			EndLine:   toInt(m[4]),
			EndCol:    toInt(m[5]),
			NumStmt:   toInt(m[6]),
			Count:     toInt(m[7]),
		})
	}

	return toProfiles(files), nil
}

func toProfiles(files map[string]*Profile) Profiles {
	profiles := make([]Profile, 0, len(files))
	for _, p := range files {
		profiles = append(profiles, *p)
	}

	sort.SliceStable(profiles, func(i, j int) bool {
		return profiles[i].FileName < profiles[j].FileName
	})

	for i, p := range profiles {
		profiles[i].ID = i

		blocks := p.Blocks.Filter()
		sort.SliceStable(blocks, func(i, j int) bool {
			bi, bj := blocks[i], blocks[j]
			return bi.StartLine < bj.StartLine || bi.StartLine == bj.StartLine && bi.StartCol < bj.StartCol
		})
		profiles[i].Blocks = blocks
	}

	return profiles
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
