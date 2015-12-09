package summarize

import (
	"fmt"
	"log"
	"os"
	"sort"

	"golang.org/x/tools/container/intsets"
)

var info = log.New(os.Stdout, "INFO: ", log.Lshortfile)
var dbg = log.New(os.Stdout, "DEBUG: ", log.Lshortfile)

// Value is an assignment for the summary
type Value struct {
	attributeType Type   // attribute type
	attributeName string // attribute name
	value         string // value
}

// Summary is a summary
// type Summary [][]Value
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

func getBestFormulaCell(formulaRankedCells CellPointers, formula Formula, skipTheseCells *intsets.Sparse) *Cell {
	// the best improvement in coverage for any cell
	bestDiff := 0

	var bestCell *Cell

	for _, cell := range formulaRankedCells {
		if skipTheseCells.Has(cell.uid) {
			dbg.Println("Skipping cell")
		}

		if cell.attribute.attributeType == single && formula.usedSingleAttributes.Contains(cell.attribute.index) {
			// the formula already has a value assigned to this attribute
			// todo: cound delete here
			dbg.Println("Ignoring single attribute cell", cell)
			skipTheseCells.Insert(cell.uid)
			continue
		}

		// how much does adding the cell to the formula change the coverage
		valueDiff := 0

		for tuple, value := range formula.tupleValue {
			covered, has := (*cell.cover)[tuple]
			if has {
				// no conflict
				if !covered {
					// and cell is not yet covered, great
					valueDiff++
				}
			} else {
				// conflict, need to remove whatever we already have for this tuple
				valueDiff -= value
			}
		}

		dbg.Printf("Adding cell %s adds %d", cell, valueDiff)

		if valueDiff > bestDiff {
			bestCell = cell
			bestDiff = valueDiff
		}
	}

	if bestDiff == 0 {
		// we could not improve the coverage so let's give up
		info.Println("Looks like we cannot find a cell that should be added")
		return nil
	}

	return bestCell
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
		fmt.Println("==============================")

		// add new formula with best cell
		cell := getBestCell(rankedCells)

		if cell.potential < 0 {
			info.Println("Adding a new cell to the formula doesn't help. Let's stop right here.")
			break
		}

		formula := NewFormula(cell)

		dbg.Printf("Just added a new formula with cell %s\n", cell)
		fmt.Println(formula.tupleValue)

		// make a copy of the ranked cells, we can use this now in the context of a formula and remove elements and reorder
		// note that CellPointers has pointers so we can safely modify the slice but not the cells it points to
		formulaRankedCells := make(CellPointers, len(rankedCells))
		copy(formulaRankedCells, rankedCells)

		// which cells to skip in formulaRankedCells, should be reset for each formula
		var skipTheseCells intsets.Sparse

		// keep adding to formula
		for true {
			bestCell := getBestFormulaCell(formulaRankedCells, formula, &skipTheseCells)

			if bestCell == nil {
				break
			}

			// we should skip the best cell in the next iteration
			skipTheseCells.Insert(bestCell.uid)

			// dbg.Printf("Now skipping %d cells\n", skipTheseCells.Len())

			// add cell to formula
			formula.AddCell(bestCell)

			info.Printf("Just added a new cell (%s) to the formula\n", bestCell)

			break
		}

		// set cover in index
		formula.CoverIndex(&relation)

		info.Println("Formula")
		fmt.Println(formula.cells)
		info.Println("Relation")
		fmt.Println(relation)

		summary = append(summary, formula)
		break
	}

	return summary
}
