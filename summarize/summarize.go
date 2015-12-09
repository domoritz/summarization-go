package summarize

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"

	"golang.org/x/tools/container/intsets"
)

var info = log.New(os.Stdout, "INFO: ", log.Lshortfile)
var dbg = log.New(os.Stdout, "DEBUG: ", log.Lshortfile)

// TupleCover is a map from tuple index to whether it is covered or not
type TupleCover map[int]bool

type tupleCover []int

// Formula is a map from attribute id to lists of cells
type Formula map[int][]Cell // TODO

// Summary is a summary
type Summary []Formula

func makeRankedCells(relation RelationIndex) CellPointers {
	rankedCells := make(CellPointers, 0, relation.numValues)
	uid := 0
	for i, attr := range relation.attrs {
		for value, cover := range attr.tuples {
			cell := Cell{uid, &relation.attrs[i], value, cover, len(*cover), 0}
			rankedCells = append(rankedCells, &cell)
			uid++
		}
	}
	return rankedCells
}

// returns the best cell form a list of cells with potentials
// requires them to be sorted and requires that the true potential of a cell is less than the given potential
func getBestCell(sortedCells CellPointers) *Cell {
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

// Summarize summarizes
func (relation RelationIndex) Summarize(size int) Summary {
	fmt.Println(relation)

	var summary Summary

	rankedCells := makeRankedCells(relation)
	sort.Sort(rankedCells)

	dbg.Println("Initial ranking")
	fmt.Println(rankedCells)

	for len(summary) < size {
		// add new formula with best cell
		cell := getBestCell(rankedCells)

		if cell.potential < 0 {
			info.Println("Adding a new cell to the formula doesn't help. Let's stop right here.")
			break
		}

		formula := make(Formula)
		formula[cell.attribute.index] = []Cell{*cell}

		// how much does the current formula contribute to the coverage
		theTupleCover := make(tupleCover, relation.numTuples)
		tuplesInFormula := make(Set)

		for tuple, covered := range *cell.cover {
			if !covered {
				theTupleCover[tuple]++
			}
			tuplesInFormula.Add(tuple)
		}

		dbg.Printf("Just added a new formula with cell %s, here is the tuple cover\n", cell)
		fmt.Println(theTupleCover)

		// make a copy of the ranked cells, we can use this now in the context of a formula and remove elements and reorder
		// note that CellPointers has pointers so we can safely modify the slice but not the cells it points to
		formulaRankedCells := make(CellPointers, len(rankedCells))
		copy(formulaRankedCells, rankedCells)

		// which cells to skip in formulaRankedCells, should be reset for each formula
		var skipTheseCells intsets.Sparse

		// keep adding to formula
		for true {
			var bestCell *Cell

			// the best improvement in coverage for any cell
			bestDiff := 0

			for _, cell := range formulaRankedCells {
				if skipTheseCells.Has(cell.uid) {
					dbg.Println("Skipping cell")
				}

				if cell.attribute.attributeType == single && len(formula[cell.attribute.index]) > 0 {
					// the formula already has a value assigned to this attribute
					// todo: cound delete here
					dbg.Println("Ignoring single attribute cell", cell)
					skipTheseCells.Insert(cell.uid)
					continue
				}

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

			break

			if bestDiff == 0 {
				// we could not improve the coverage so let's give up
				info.Println("Looks like we cannot find a cell that should be added")
				break
			}

			skipTheseCells.Insert(bestCell.uid)
			dbg.Printf("Now skipping %d cells\n", skipTheseCells.Len())

			// add cell to formula
			idx := bestCell.attribute.index
			formula[idx] = append(formula[idx], *bestCell)

			// shrink the relevant tuples
			for tuple := range tuplesInFormula {
				if _, has := (*bestCell.cover)[tuple]; !has {
					delete(tuplesInFormula, tuple)
				} else {
					theTupleCover[tuple]++
				}
			}

			dbg.Println("Relevant tuples:", tuplesInFormula)

			info.Printf("Just added a new cell (%s) to the formula\n", bestCell)
			fmt.Println(theTupleCover)
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
		break
	}

	return summary
}

func (cover tupleCover) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Cover (%d):\n", len(cover)))
	for i, cover := range cover {
		buffer.WriteString(fmt.Sprintf("%d: %d\n", i, cover))
	}
	return buffer.String()
}