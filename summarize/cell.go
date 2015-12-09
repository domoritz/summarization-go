package summarize

import (
	"bytes"
	"fmt"
)

// TupleCover is a map from tuple index to how much the tuple contributes
type TupleCover map[int]bool

// Cell is a cell
type Cell struct {
	uid       int         // unique id for cell
	attribute *Attribute  // attribute
	value     string      // attribute value
	cover     *TupleCover // what the attribute covers

	potential        int // what the cell can cover in the whole relation
	formulaPotential int // what the cell can cover in the context of a formula (has to be kept track of)
}

// CellPointers is a list of pointers to cells
type CellPointers []*Cell

// Len is part of sort.Interface.
func (cells CellPointers) Len() int { return len(cells) }

// Swap is part of sort.Interface.
func (cells CellPointers) Swap(i, j int) {
	cells[i], cells[j] = cells[j], cells[i]
}

// Less is part of sort.Interface. Sort by Potential.
func (cells CellPointers) Less(i, j int) bool {
	// todo: prefer shorter prefixes to break ties
	return cells[i].potential > cells[j].potential
}

// Push pushes
func (cells *CellPointers) Push(x interface{}) {
	*cells = append(*cells, x.(*Cell))
}

// Pop pops
func (cells *CellPointers) Pop() interface{} {
	old := *cells
	n := len(old)
	x := old[n-1]
	*cells = old[0 : n-1]
	return x
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

func (cells CellPointers) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Cells (%d):\n", len(cells)))
	for _, cell := range cells {
		buffer.WriteString(fmt.Sprintf("%s\n", cell))
	}
	return buffer.String()
}

func (cell Cell) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Attr %s: %s (%d)", cell.attribute.attributeName, cell.value, cell.potential))
	return buffer.String()
}
