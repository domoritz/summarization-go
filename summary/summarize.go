package summary

import (
	"fmt"
	"sort"
)

// CellKey can be used to idenfity a cell
type CellKey struct {
	Type      Type
	Attribute int
	Value     string
}

// Cell is a cell
type Cell struct {
	CellKey
	Potential  int
	Attributes []*bool // pointers to whether attribute value is covered
	Tuples     []*Tuple
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
	return cells[i].Potential > cells[j].Potential
}

// DebugPrint prints cell slice
func (cells cellSlice) DebugPrint() {
	for _, cell := range cells {
		fmt.Println(cell)
	}
	fmt.Println()
}

// Summarize returns a summary of the relation
func (relation Relation) Summarize(size int) Relation {
	summary := Relation{
		make([]*Tuple, 0),
		relation.AttributeNames,
		relation.AttributeTypes,
		relation.GetSizes(),
	}

	mainCells := allCells(relation.Tuples)
	sort.Sort(mainCells)
	mainCells.DebugPrint()

	for len(summary.Tuples) < size {
		// add a new formula
		firstCell := mainCells[0]
		formula := NewTupleFromCell(firstCell, summary.GetSizes())
		summary.Tuples = append(summary.Tuples, &formula)
		SetToCovered(firstCell)
		tuples := firstCell.Tuples

		relation.PrintDebugString()

		// the cells that are relevant for the current set of tuples
		cells := allCells(tuples)
		sort.Sort(cells)

		// the second cell is in the subset
		secondCell := cells[0]
		formula.AddCell(secondCell)
		SetToCovered(secondCell)

		// cannot add the currently best cell anyway
		cells = cells[1:]

		// cell := cells.findBestCell(tuples)

		// TODO: No need to update potential if we add a single
		sort.Sort(mainCells)

		break
	}

	return summary
}

func (cells cellSlice) findBestCell(tuples []Tuple) Cell {
	bestCover := 0
	var bestCell Cell
	for _, cell := range cells {
		if cell.Potential < bestCover {
			bestCell = cell
			break
		}

	}
	return bestCell
}

// SetToCovered sets to covered
func SetToCovered(cell Cell) {
	for _, covered := range cell.Attributes {
		*covered = true
	}
}

func (c Cell) String() string {
	return fmt.Sprintf("{%s %d %s (%d)}", c.Type, c.Attribute, c.Value, c.Potential)
}

func allCells(tuples []*Tuple) cellSlice {
	cells := make(map[CellKey]*Cell)

	for _, tuple := range tuples {
		tuple.AddCells(&cells)
	}

	ret := make(cellSlice, 0)

	for _, cell := range cells {
		// ignore cells with potential < 2
		if cell.Potential > 1 {
			ret = append(ret, *cell)
		}
	}

	return ret
}
