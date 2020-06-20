package profile

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/tools/cover"
)

var reg = regexp.MustCompile(`^(.+):([0-9]+).([0-9]+),([0-9]+).([0-9]+) ([0-9]+) ([0-9]+)$`)

type Profiles []cover.Profile

func Scan(s *bufio.Scanner) (Profiles, error) {
	files := make(map[string]*cover.Profile)
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
			p = &cover.Profile{
				FileName: fileName,
				Mode:     mode,
			}
			files[fileName] = p
		}

		p.Blocks = append(p.Blocks, cover.ProfileBlock{
			StartLine: toInt(m[2]),
			StartCol:  toInt(m[3]),
			EndLine:   toInt(m[4]),
			EndCol:    toInt(m[5]),
			NumStmt:   toInt(m[6]),
			Count:     toInt(m[7]),
		})
	}

	return sortProfile(files), nil
}

func sortProfile(files map[string]*cover.Profile) Profiles {
	profiles := make([]cover.Profile, 0, len(files))
	for _, p := range files {
		profiles = append(profiles, *p)
	}

	sort.SliceStable(profiles, func(i, j int) bool {
		return profiles[i].FileName < profiles[j].FileName
	})

	return profiles
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
