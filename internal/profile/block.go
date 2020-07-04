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

func (blocks *Blocks) Filter() Blocks {
	newBlocks := make(Blocks, 0, len(*blocks))
	for _, b := range *blocks {
		i := newBlocks.index(&b)

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

func (blocks *Blocks) index(block *Block) int {
	for i, b := range *blocks {
		if b.StartLine == block.StartLine &&
			b.StartCol == block.StartCol &&
			b.EndLine == block.EndLine &&
			b.EndCol == block.EndCol {
			return i
		}
	}
	return -1
}
