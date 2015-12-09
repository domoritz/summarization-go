package summarize

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

var info = log.New(os.Stdout, "INFO: ", log.Lshortfile)

// Type is the attribute type
type Type int

const (
	single Type = iota
	set
	hierarchy
)

func (t Type) String() string {
	switch t {
	case single:
		return "single"
	case set:
		return "set"
	case hierarchy:
		return "hierarchy"
	default:
		return "unknown"
	}
}

// TupleCover is a map from tuple index to whether it is covered or not
type TupleCover map[int]bool

// Attribute is an attribute
type Attribute struct {
	attributeType Type                   // attribute type
	name          string                 // attribute name
	tuples        map[string]*TupleCover // TODO: make slice
}

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

// RelationIndex is an inverted index
type RelationIndex struct {
	attrs     []Attribute
	numTuples int
	numValues int
}

type tupleCover []int

// Summary is a summary
type Summary []Formula

// Formula is a map from attribute to lists of cells
type Formula map[int][]Cell

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
func getBestCell(rankedCells cellSlice) Cell {
	n := len(rankedCells)

	bestCoverage := 0
	for i, cell := range rankedCells {
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
	sort.Sort(rankedCells[0:n])
	return rankedCells[0]
}

// Summarize summarizes
func (relation RelationIndex) Summarize(size int) Summary {
	fmt.Println(relation)

	var summary Summary

	rankedCells := make(cellSlice, 0)
	for i, attr := range relation.attrs {
		for value, cover := range attr.tuples {
			cell := Cell{i, value, cover, len(*cover)}
			rankedCells = append(rankedCells, cell)
		}
	}

	sort.Sort(rankedCells)
	info.Println("Initial ranking")
	fmt.Println(rankedCells)

	for len(summary) < size {
		// add new formula with best cell
		cell := getBestCell(rankedCells)

		if cell.potential < 0 {
			log.Println("Adding a new cell to the formula doesn't help. Let's stop right here.")
			break
		}

		formula := make(Formula)
		formula[cell.attribute] = []Cell{cell}

		// how much does the current formula contribute to the coverage
		theTupleCover := make(tupleCover, relation.numTuples)
		tuplesInFormula := make(map[int]bool)

		for tuple, covered := range *cell.cover {
			if !covered {
				theTupleCover[tuple]++
			}
			tuplesInFormula[tuple] = true
		}

		info.Println("Just added a new formula, here is the tuple cover")
		fmt.Println(theTupleCover)

		// keep adding to formula
		for true {
			var bestCell Cell

			// the best improvement in coverage for any cell
			bestDiff := 0

			for _, cell := range rankedCells {
				// todo: ignore cells for the same attribute if it is single
				// delete it from slice

				// how much does adding the cell to the formula change the coverage
				coverageDiff := 0

				for tuple := range tuplesInFormula {
					covered, has := (*cell.cover)[tuple]
					if has {
						// no conflict
						if !covered {
							// and cell is not yet covered, great
							coverageDiff++
						}
					} else {
						// conflict, need to remove whatever we already have for this tuple
						coverageDiff -= theTupleCover[tuple]
					}
				}

				if coverageDiff > bestDiff {
					bestCell = cell
					bestDiff = coverageDiff
				}
			}

			if bestDiff == 0 {
				// we could not improve the coverage so let's give up
				info.Println("Looks like we cannot find a cell that should be added")
				break
			}

			// add cell to formula
			formulaCells := formula[bestCell.attribute]
			formulaCells = append(formulaCells, bestCell)

			// shrink the relevant tuples
			for tuple := range tuplesInFormula {
				if _, has := (*bestCell.cover)[tuple]; !has {
					delete(tuplesInFormula, tuple)
				} else {
					theTupleCover[tuple]++
				}
			}

			info.Println("Relevant tuples:", tuplesInFormula)

			info.Printf("Just added a new cell (%s) to the formula\n", bestCell)
			fmt.Println(theTupleCover)

			break
		}

		// update the cover so that in the next iteration the same tuples are not covered again
		for _, formulaCells := range formula {
			for _, cell := range formulaCells {
				for tuple := range tuplesInFormula {
					if covered, has := (*cell.cover)[tuple]; has && !covered {
						// set uncovered to covered
						(*cell.cover)[tuple] = true
					}
				}
			}
		}

		info.Printf("After adding a new formula %s the relation looks like\n", formula)
		fmt.Println(relation)

		summary = append(summary, formula)
	}

	return summary
}

func (cells cellSlice) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Cells (%d):\n", len(cells)))
	for _, cell := range cells {
		buffer.WriteString(fmt.Sprintf("%s\n", cell))
	}
	return buffer.String()
}

func (cover tupleCover) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Cover (%d):\n", len(cover)))
	for i, cover := range cover {
		buffer.WriteString(fmt.Sprintf("%d: %d\n", i, cover))
	}
	return buffer.String()
}

func (cell Cell) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Attr %d: %s (%d)", cell.attribute, cell.value, cell.potential))
	return buffer.String()
}

func (relation RelationIndex) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Relation Index (%d attributes, %d tuples, %d values):\n", len(relation.attrs), relation.numTuples, relation.numValues))
	for _, attribute := range relation.attrs {
		buffer.WriteString(fmt.Sprintf("Attribute %s (%s):\n", attribute.name, attribute.attributeType))
		for value, cell := range attribute.tuples {
			buffer.WriteString(fmt.Sprintf("Value %s covers: [", value))
			var tuples []string
			for tuple, covered := range *cell {
				tuples = append(tuples, fmt.Sprintf("%d: %t", tuple, covered))
			}

			buffer.WriteString(strings.Join(tuples, " "))

			buffer.WriteString("]\n")
		}
	}

	return buffer.String()
}
