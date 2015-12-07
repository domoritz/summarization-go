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
	Tuples     []*Tuple  // pointers to tuple
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
		relation.IncreaseCounts(cell)
	}

	relation.PrintDebugString()

	for true {
		// how much we can cover
		bestPotential := 0
		var bestCell *Cell
		var formulaToAdd *Tuple

		for _, cell := range cells {
			if cell.Potential < bestPotential {
				fmt.Printf("No way to get better from here on")
				break
			}

			{
				// how much new space can we cover with a new formula
				// a new formula can not have any conflicts so this is easy

				// TODO: we might already know the potential because it is calculated
				potential := 0
				for _, c := range cell.Attributes {
					// count non-covered cells
					if *c == 0 {
						potential++
					}
				}

				fmt.Printf("Update potential of %s from %d to %d\n", cell, cell.Potential, potential)
				cell.Potential = potential

				if potential > bestPotential && len(summary.Tuples) < size {
					bestPotential = potential
					bestCell = cell
					formulaToAdd = nil
				}
			}

			if cell.Potential < bestPotential {
				fmt.Printf("No way to get better from here on")
				break
			}

			// how about adding to an existing formula?
			for _, formula := range summary.Tuples {
				if formula.SatisfiesCell(cell) {
					// formula already has attribute, so skip it
					continue
				}

				// remove counts, let's assume we add a new formula, wht do we get
				for _, fCell := range formula.Cells {
					relation.DecreaseCounts(fCell)
				}

				// how much can we cover if we add this formula again
				potential := 0
				for _, tuple := range relation.Tuples {
					if tuple.Satisfies(formula) && tuple.SatisfiesCell(cell) {
						for _, fCell := range formula.Cells {
							switch fCell.Type {
							case single:
								if *tuple.Single[fCell.Attribute].covered == 0 {
									potential++
								}
							case set:
								if *tuple.Set[fCell.Attribute].values[fCell.Value] == 0 {
									potential++
								}
							}
						}
					}
				}

				fmt.Printf("Adding %s to %p has potential %d\n", cell, formula, potential)

				if potential > bestPotential {
					potential = bestPotential
					bestCell = cell
					formulaToAdd = formula
				}

				// add them back
				for _, fCell := range formula.Cells {
					relation.IncreaseCounts(fCell)
				}
			}
		}

		if bestPotential == 0 {
			fmt.Println("Done")
			break
		}

		fmt.Printf("Adding %s has best potential %d in formula %p\n", bestCell, bestPotential, formulaToAdd)

		if formulaToAdd == nil {
			formula := NewTupleFromCell(bestCell, summary.GetSizes())
			summary.Tuples = append(summary.Tuples, &formula)
			relation.IncreaseCounts(bestCell)
		} else {
			for _, fCell := range formulaToAdd.Cells {
				relation.DecreaseCounts(fCell)
			}

			formulaToAdd.Cells = append(formulaToAdd.Cells, bestCell)

			switch bestCell.Type {
			case single:
				attr := NewSingle(bestCell.Value)
				c := 1
				attr.covered = &c
				formulaToAdd.Single[bestCell.Attribute] = &attr
			case set:
				attr := formulaToAdd.Set[bestCell.Attribute]
				c := 1
				if attr == nil {
					set := NewSet(Set{bestCell.Value: &c})
					formulaToAdd.Set[bestCell.Attribute] = &set
				} else {
					attr.values[bestCell.Value] = &c
				}
			}

			formulaToAdd.Tuples = nil
			for _, tuple := range relation.Tuples {
				if tuple.Satisfies(formulaToAdd) {
					formulaToAdd.Tuples = append(formulaToAdd.Tuples, tuple)
					for _, fCell := range formulaToAdd.Cells {
						switch fCell.Type {
						case single:
							*tuple.Single[fCell.Attribute].covered++
						case set:
							*tuple.Set[fCell.Attribute].values[fCell.Value]++
						}
					}
				}
			}

		}

		// TODO: No need to update potential if we add a single

		sort.Sort(cells)
	}

	return summary
}

// IncreaseCounts increases the counts
func (relation *Relation) IncreaseCounts(cell *Cell) {
	for _, counter := range cell.Attributes {
		(*counter)++
	}
}

// DecreaseCounts increases the counts
func (relation *Relation) DecreaseCounts(cell *Cell) {
	for _, counter := range cell.Attributes {
		(*counter)--
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
