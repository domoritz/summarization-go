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
	var tuples []*Tuple
	summary := Relation{tuples, relation.AttributeNames, relation.AttributeTypes, relation.GetSizes()}

	cells := relation.allCells()
	sort.Sort(cells)

	cells.DebugPrint()

	{
		// we can only add a formula so we know it's going to be the one from the best cell
		cell := cells[0]
		formula := NewTupleFromCell(cell, summary.GetSizes())
		summary.Tuples = append(summary.Tuples, &formula)
		relation.SetToCovered(cell)
	}

	relation.PrintDebugString()

	for true {

		// TODO: No need to update potential if we add a single

		sort.Sort(cells)

		break
	}

	return summary
}

// SetToCovered sets to covered
func (relation Relation) SetToCovered(cell Cell) {
	for _, covered := range cell.Attributes {
		*covered = true
	}
}

func (c Cell) String() string {
	return fmt.Sprintf("{%s %d %s (%d)}", c.Type, c.Attribute, c.Value, c.Potential)
}

// returns all cells with potential
func (relation *Relation) allCells() cellSlice {
	cells := make(map[CellKey]*Cell)

	for _, tuple := range relation.Tuples {
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
