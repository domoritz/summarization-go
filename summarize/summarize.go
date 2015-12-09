package summarize

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
)

var info = log.New(os.Stdout, "INFO: ", log.Lshortfile)

// TupleCover is a map from tuple index to whether it is covered or not
type TupleCover map[int]bool

type tupleCover []int

// Formula is a map from attribute to lists of cells
type Formula map[int][]Cell // TODO

// Summary is a summary
type Summary []Formula

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
