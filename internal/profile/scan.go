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
	id := 1

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
			p = &Profile{
				FileName: fileName,
				Mode:     mode,
			}
			files[fileName] = p
			id++
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
		p.Blocks = filterBlocks(p.Blocks)
		profiles = append(profiles, *p)
	}

	sort.SliceStable(profiles, func(i, j int) bool {
		return profiles[i].FileName < profiles[j].FileName
	})

	// ID振る
	for i := range profiles {
		profiles[i].ID = i
	}

	return profiles
}

func filterBlocks(blocks Blocks) Blocks {
	pbm := make(map[string]Block)
	for _, b := range blocks {
		// TODO: やり方考える
		k := fmt.Sprintf("%d-%d-%d-%d", b.StartLine, b.StartCol, b.EndLine, b.EndCol)
		if _, ok := pbm[k]; !ok {
			pbm[k] = b
		}
		if b.Count > 0 {
			pbm[k] = b
		}
	}

	pbs := make(Blocks, 0, len(pbm))
	for _, b := range pbm {
		pbs = append(pbs, b)
	}

	sort.SliceStable(pbs, func(i, j int) bool {
		bi, bj := pbs[i], pbs[j]
		return bi.StartLine < bj.StartLine || bi.StartLine == bj.StartLine && bi.StartCol < bj.StartCol
	})

	return pbs
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}
