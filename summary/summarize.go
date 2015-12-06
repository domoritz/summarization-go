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
	Potential int
}

type cellSlice []*Cell

// Len is part of sort.Interface.
func (d cellSlice) Len() int {
	return len(d)
}

// Swap is part of sort.Interface.
func (d cellSlice) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface. Sort by Potential.
func (d cellSlice) Less(i, j int) bool {
	return d[i].Potential > d[j].Potential
}

// DebugPrint prints cell slice
func (d *cellSlice) DebugPrint() {
	for _, cell := range *d {
		fmt.Println(*cell)
	}
	fmt.Println()
}

// Summarize returns a summary of size size
func (relation *Relation) Summarize(size int) Relation {
	var tuples []Tuple
	summary := Relation{tuples, relation.AttributeNames, relation.AttributeTypes, relation.GetSizes()}

	cells := relation.allCells()
	sort.Sort(cells)

	cells.DebugPrint()

	summary.Tuples = append(summary.Tuples, NewTupleFromCell(*cells[0], summary.GetSizes()))
	cells = cells[1:]

	for true {
		for _, cell := range cells {
			// TODO: check potential

			if len(summary.Tuples) < size {
				// try new formula
				summary.Tuples = append(summary.Tuples, NewTupleFromCell(*cell, summary.GetSizes()))
				// coverage = relation.Covers(summary)
			}

		}
		break
	}

	return summary
}

// returns all cells with potential
func (relation *Relation) allCells() cellSlice {
	cells := make(map[CellKey]*Cell)

	for _, tuple := range relation.Tuples {
		tuple.AddCells(&cells)
	}

	ret := make(cellSlice, 0, len(cells))

	for _, cell := range cells {
		ret = append(ret, cell)
	}

	return ret
}
