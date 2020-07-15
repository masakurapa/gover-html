package profile

import (
	"math"
	"path"
	"path/filepath"
	"strings"
)

// Profiles is a type that represents a slice of Profile
type Profiles []Profile

// Profile is profiling data for each file
type Profile struct {
	ID        int
	Dir       string
	FileName  string
	Blocks    Blocks
	Functions Functions
}

// Functions is a type that represents a slice of Function
type Functions []Function

// Function is single func of profiling data
type Function struct {
	Name        string
	StartLine   int
	StartCol    int
	TypeEndLine int
	TypeEndCol  int
	Blocks      Blocks
}

// Blocks is a type that represents a slice of Block
type Blocks []Block

// Block is single block of profiling data
type Block struct {
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
	NumState  int
	Count     int
}

// IsRelativeOrAbsolute returns true if FileName is relative path or absolute path
func (prof *Profile) IsRelativeOrAbsolute() bool {
	return strings.HasPrefix(prof.FileName, ".") || filepath.IsAbs(prof.FileName)
}

// FilePath returns readable file path
func (prof *Profile) FilePath() string {
	if prof.IsRelativeOrAbsolute() {
		return prof.FileName
	}
	return filepath.Join(prof.Dir, path.Base(prof.FileName))
}

// Coverage returns covered ratio for file
func (blocks *Blocks) Coverage() float64 {
	var total, covered int64
	for _, b := range *blocks {
		total += int64(b.NumState)
		if b.Count > 0 {
			covered += int64(b.NumState)
		}
	}

	if total == 0 {
		return 0
	}

	return math.Round((float64(covered)/float64(total)*100)*10) / 10
}
