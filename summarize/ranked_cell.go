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
// returns the new formula cover and the cell cover
func (cell *RankedCell) recomputeFormulaCoverage(formula *Formula) (int, int) {
	before := cell.potential

	formulaCover := 0
	cell.potential = 0

	// doing this optimizations saves about 25% time
	if len(formula.tupleCover) <= len(cell.cell.cover) {
		for tuple, cover := range formula.tupleCover {
			covered, has := cell.cell.cover[tuple]
			if has {
				formulaCover += cover
				// no conflict
				if !covered {
					// and cell is not yet covered, great
					cell.potential++
					formulaCover++
				}
			}
		}
	} else {
		for tuple, covered := range cell.cell.cover {
			cover, has := formula.tupleCover[tuple]
			if has {
				formulaCover += cover
				// no conflict
				if !covered {
					// and cell is not yet covered, great
					cell.potential++
					formulaCover++
				}
			}
		}
	}

	if before < cell.potential {
		panic("not smaller")
	}

	if formulaCover < cell.potential {
		panic("potential is too large")
	}

	return formulaCover, cell.potential
}

func (cells CellHeap) String() string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "Cells (%d):\n", len(cells))
	for _, cell := range cells {
		fmt.Fprintf(&buffer, "%s\n", cell)
	}
	return buffer.String()
}

func (cell RankedCell) String() string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "%s (%d)", cell.cell, cell.potential)
	return buffer.String()
}
