package summarize

import (
	"bytes"
	"fmt"
	"sort"
)

// Cell is a cell
type Cell struct {
	attribute int         // attribute
	value     string      // attribute value
	cover     *TupleCover // what the attribute covers
	potential int         // what the cell can cover in the whole relation
}

type cellSlice []Cell

// Len is part of sort.Interface.
func (cells cellSlice) Len() int {
	return len(cells)
}

// Swap is part of sort.Interface.
func (cells cellSlice) Swap(i, j int) {
	cells[i], cells[j] = cells[j], cells[i]
}

// Less is part of sort.Interface. Sort by Potential.
func (cells cellSlice) Less(i, j int) bool {
	// todo: prefer shorter prefixes to break ties
	return cells[i].potential > cells[j].potential
}

// recomputes how much the tuple covers
func (cell *Cell) recomputeCoverage() int {
	cell.potential = 0
	for _, covered := range *cell.cover {
		if !covered {
			cell.potential++
		}
	}
	return cell.potential
}

// returns the best cell form a list of cells with potentials
// requires them to be sorted and requires that the true potential of a cell is less than the given potential
func getBestCell(sortedCells cellSlice) Cell {
	if !sort.IsSorted(sortedCells) {
		panic("Not sorted")
	}

	n := len(sortedCells)

	bestCoverage := 0
	for i, cell := range sortedCells {
		if cell.potential > bestCoverage {
			coverage := cell.recomputeCoverage()
			if coverage > bestCoverage {
				bestCoverage = coverage
			}
		} else {
			// potential is lower than the best so far
			n = i
			break
		}
	}

	// sort the range where we recomputed things, the rest is definitely lower
	sort.Sort(sortedCells[0:n])
	return sortedCells[0]
}

func (cells cellSlice) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Cells (%d):\n", len(cells)))
	for _, cell := range cells {
		buffer.WriteString(fmt.Sprintf("%s\n", cell))
	}
	return buffer.String()
}

func (cell Cell) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Attr %d: %s (%d)", cell.attribute, cell.value, cell.potential))
	return buffer.String()
}
