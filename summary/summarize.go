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
	Attributes []Counter
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

// Summarize returns a summary of the relation
func (relation *Relation) Summarize(size int) Relation {
	var tuples []Tuple
	summary := Relation{tuples, relation.AttributeNames, relation.AttributeTypes, relation.GetSizes()}

	cells := relation.allCells()
	sort.Sort(cells)

	cells.DebugPrint()

	// we can only add a formula so we know it's going to be the one from the best cell
	cell := *cells[0]
	formula := NewTupleFromCell(cell, summary.GetSizes())
	summary.Tuples = append(summary.Tuples, formula)
	relation.IncreaseCounts(cell)

	cells = cells[1:]

	relation.PrintDebugString()

	// for true {
	// 	for _, cell := range cells {
	// 		// TODO: check potential
	//
	// 		if len(summary.Tuples) < size {
	// 			// try new formula and see how much of uncovered space it can cover
	// 			formula := NewTupleFromCell(*cell, summary.GetSizes())
	// 		}
	//
	// 		// try adding to existing formula and see how it changes coverage
	// 		for i, existing := range summary.Tuples {
	// 		}
	//
	// 	}
	// 	break
	// }

	return summary
}

// IncreaseCounts increases the counts
func (relation *Relation) IncreaseCounts(cell Cell) {
	for _, counter := range cell.Attributes {
		(*counter)++
	}
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
