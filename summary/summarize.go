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
	Attributes []Counter // pointers to counter in attribute
	Tuples     []*Tuple   // pointers to tuple
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
	formula.Cover = cell.Potential
	summary.Tuples = append(summary.Tuples, formula)
	relation.IncreaseCounts(cell)

	// only need the remaining cells
	cells = cells[1:]

	relation.PrintDebugString()

	relation.IncreaseCounts(*cells[0])

	relation.PrintDebugString()

	for true {
		bestPotential := 0
		var bestCell *Cell

		for _, cell := range cells {
			// TODO: check potential

			if len(summary.Tuples) < size {
				// how much new space can we cover with a new formula
				// a new formula can not have any conflicts so this is easy

				// TODO: we might already know the potential because it is calculated
				potential := 0
				for _, c := range cell.Attributes {
					if *c == 0 {
						potential++
					}
				}

				fmt.Printf("Adding %s has potential %d\n", cell, potential)

				if potential > bestPotential {
					bestPotential = potential
					bestCell = cell
				}
			}

			// how about adding to an existing formula?
			for _, formula := range summary.Tuples {
				if cell.Type == single {
					if formula.Single[cell.Attribute] != nil {
						// formula already has attribute, so skip it
						continue
					}

				}
			}
		}

		fmt.Printf("Adding %s has best potential %d\n", bestCell, bestPotential)

		// No need to update potential if we add a single

		break
	}

	return summary
}

// IncreaseCounts increases the counts
func (relation *Relation) IncreaseCounts(cell Cell) {
	for _, counter := range cell.Attributes {
		(*counter)++
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
			ret = append(ret, cell)
		}
	}

	return ret
}
