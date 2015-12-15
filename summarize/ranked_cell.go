package summarize

import (
	"bytes"
	"fmt"
)

// RankedCell is a cell pointer and a priority
type RankedCell struct {
	cell *Cell // pointer to cell

	// potential is what the cell can cover in the whole relation
	// constraint: potential must always be higher than actual cover
	potential int
}

// CellHeap is a heap of ranked cells
type CellHeap []RankedCell

// Len is part of sort.Interface.
func (cells CellHeap) Len() int { return len(cells) }

// Swap is part of sort.Interface.
func (cells CellHeap) Swap(i, j int) {
	cells[i], cells[j] = cells[j], cells[i]
}

// Less is part of sort.Interface. Sort by Potential.
func (cells CellHeap) Less(i, j int) bool {
	// todo: prefer shorter prefixes to break ties
	return cells[i].potential > cells[j].potential
}

// Push pushes
func (cells *CellHeap) Push(x interface{}) {
	*cells = append(*cells, x.(RankedCell))
}

// Pop pops
func (cells *CellHeap) Pop() interface{} {
	old := *cells
	n := len(old)
	item := old[n-1]
	*cells = old[0 : n-1]
	return item
}

// Peek returns a pointer to the best cell
func (cells CellHeap) Peek() *RankedCell {
	return &cells[0]
}

// recomputes how much the tuple covers
// returns the potential
func (cell *RankedCell) recomputeCoverage() int {
	cell.potential = 0

	for _, covered := range cell.cell.cover {
		if !covered {
			cell.potential++
		}
	}

	return cell.potential
}

// recomputes what this cell covers in the context of the formula
// returns the actual cell cover and the difference in cover for the formula
func (cell *RankedCell) recomputeFormulaCoverage(formula *Formula) (int, int) {
	before := cell.potential

	coverDiff := 0
	cell.potential = 0

	for tuple, cover := range formula.tupleCover {
		covered, has := cell.cell.cover[tuple]
		if has {
			// no conflict
			if !covered {
				// and cell is not yet covered, great
				coverDiff++
				cell.potential++
			}
		} else {
			// conflict, need to remove whatever we already have for this tuple
			coverDiff -= cover
		}
	}

	if before < cell.potential {
		panic("not smaller")
	}

	return cell.potential, coverDiff
}

func (cells CellHeap) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Cells (%d):\n", len(cells)))
	for _, cell := range cells {
		buffer.WriteString(fmt.Sprintf("%s\n", cell))
	}
	return buffer.String()
}

func (cell RankedCell) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s (%d)", cell.cell, cell.potential))
	return buffer.String()
}
