package profile

type Blocks []Block

type Block struct {
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
	NumStmt   int
	Count     int
}

func (blocks *Blocks) Coverage() float64 {
	var total, covered int64
	for _, b := range *blocks {
		total += int64(b.NumStmt)
		if b.Count > 0 {
			covered += int64(b.NumStmt)
		}
	}

	if total == 0 {
		return 0
	}

	return float64(covered) / float64(total) * 100
}
