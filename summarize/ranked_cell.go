package summarize

import (
	"bytes"
	"fmt"
)

// RankedCell is a cell pointer and a priority
type RankedCell struct {
	cell *Cell // pointer to cell

	potential    int // potential is what the cell can cover in the whole relation or in the context of a formula, constraint: potential must always be higher than actual cover
	maxPotential int // the maximum potential that the cell can have in the context of a formula, can be used to reset potential
	index        int // The index of the item in the heap.
}

// CellHeap is a heap of ranked cells
type CellHeap []*RankedCell

// Len is part of sort.Interface.
func (cells CellHeap) Len() int { return len(cells) }

// Swap is part of sort.Interface.
func (cells CellHeap) Swap(i, j int) {
	cells[i], cells[j] = cells[j], cells[i]
	cells[i].index = i
	cells[j].index = j
}

// Less is part of sort.Interface. Sort by Potential.
func (cells CellHeap) Less(i, j int) bool {
	p1 := cells[i]
	p2 := cells[j]
	if p1.potential == p2.potential {
		if p1.cell.attribute.index == p2.cell.attribute.index && p1.cell.attribute.attributeType == hierarchy && p2.cell.attribute.attributeType == hierarchy {
			// prefer shorter hierarchies
			return len(cells[i].cell.value) < len(cells[i].cell.value)
		}
	}
	return p1.potential > p2.potential
}

// Push pushes
func (cells *CellHeap) Push(x interface{}) {
	n := len(*cells)
	cell := x.(*RankedCell)
	cell.index = n
	*cells = append(*cells, cell)
}

// Pop pops
func (cells *CellHeap) Pop() interface{} {
	old := *cells
	n := len(old)
	cell := old[n-1]
	cell.index = -1 // for safety
	*cells = old[0 : n-1]
	return cell
}

// Peek returns a pointer to the best cell
func (cells CellHeap) Peek() *RankedCell {
	return cells[0]
}

// Valid checks whether the heap is a heap
func (cells CellHeap) Valid(i int) bool {
	n := cells.Len()
	j1 := 2*i + 1
	j2 := 2*i + 2
	if j1 < n {
		if cells.Less(j1, i) {
			fmt.Printf("heap invariant invalidated [%d] = %d > [%d] = %d", i, cells[i], j1, cells[j1])
			return false
		}
		if !cells.Valid(j1) {
			return false
		}
	}
	if j2 < n {
		if cells.Less(j2, i) {
			fmt.Printf("heap invariant invalidated [%d] = %d > [%d] = %d", i, cells[i], j1, cells[j2])
			return false
		}
		if !cells.Valid(j2) {
			return false
		}
	}

	return true
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
func (cell *RankedCell) recomputeFormulaCoverage(formula *Formula) int {
	before := cell.potential

	formulaCover := 0     // what we cover in the whole formula
	cell.maxPotential = 0 // what the cell can cover at most

	// compute cover in intersection, loops over smaller list
	// doing this optimizations saves about 25% time
	if len(formula.tupleCover) <= len(cell.cell.cover) {
		for tuple, cover := range formula.tupleCover {
			covered, has := cell.cell.cover[tuple]
			if has {
				formulaCover += cover
				// no conflict
				if !covered {
					// and cell is not yet covered, great
					cell.maxPotential++
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
					cell.maxPotential++
					formulaCover++
				}
			}
		}
	}

	// the potential to cover things in the context of a formula (so we subtract what cannot be covered if we add this cell)
	// this is what actually matters when we try to find a new cell but it also is only valid with respect to the current tupleCover
	cell.potential = formulaCover - formula.cover

	if before < cell.potential {
		panic("Coverage can only decrease")
	}

	return cell.potential
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
	fmt.Fprintf(&buffer, "%s (potential: %d, max: %d)", cell.cell, cell.potential, cell.maxPotential)
	return buffer.String()
}
