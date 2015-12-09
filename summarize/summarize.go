package summarize

import (
	"container/heap"
	"fmt"
	"log"
	"os"

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
type Summary [][]Value

func makeRankedCells(relation RelationIndex) *CellPointers {
	rankedCells := make(CellPointers, 0, relation.numValues)
	uid := 0
	for i, attr := range relation.attrs {
		for value, cover := range attr.tuples {
			cell := Cell{uid, &relation.attrs[i], value, cover, len(*cover), 0}
			rankedCells = append(rankedCells, &cell)
			uid++
		}
	}
	return &rankedCells
}

// returns the best cell form a list of cells with potentials
// requires that the cells are a sorted heap
func getBestCell(cellHeap *CellPointers) *Cell {
	bestCoverage := 0

	cell := (*cellHeap)[0]
	for cell.potential > bestCoverage {
		coverage := cell.recomputeCoverage()
		heap.Fix(cellHeap, 0)
		if coverage > bestCoverage {
			bestCoverage = coverage
		}
		cell = (*cellHeap)[0]
	}

	return (*cellHeap)[0]
}

// returns nil if no cell could be found that improves the formula
// requires cells to be a heap
func getBestFormulaCell(formulaRankedCells CellPointers, formula Formula, skipTheseCells *intsets.Sparse) *Cell {
	// the best improvement in coverage for any cell
	bestDiff := 0

	var bestCell *Cell

	for _, cell := range formulaRankedCells {
		if skipTheseCells.Has(cell.uid) {
			dbg.Println("Skipping cell")
			continue
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
		cell.formulaPotential = 0

		// todo: what if we do the inverse?
		for tuple, value := range formula.tupleValue {
			covered, has := (*cell.cover)[tuple]
			if has {
				// no conflict
				if !covered {
					// and cell is not yet covered, great
					valueDiff++
					cell.formulaPotential++
				}
			} else {
				// conflict, need to remove whatever we already have for this tuple
				valueDiff -= value
			}
		}

		dbg.Printf("Adding cell %s adds %d and has potential %d", cell, valueDiff, cell.formulaPotential)

		if valueDiff > bestDiff {
			bestCell = cell
			bestDiff = valueDiff
		}
	}

	if bestDiff == 0 {
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
	heap.Init(rankedCells)

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
		formulaRankedCells := make(CellPointers, len(*rankedCells))
		copy(formulaRankedCells, *rankedCells)

		// which cells to skip in formulaRankedCells, should be reset for each formula
		var skipTheseCells intsets.Sparse

		// keep adding to formula
		for true {
			fmt.Println("--------------------------------")

			bestCell := getBestFormulaCell(formulaRankedCells, formula, &skipTheseCells)

			if bestCell == nil {
				break
			}

			// we should skip the best cell in the next iteration
			skipTheseCells.Insert(bestCell.uid)

			dbg.Printf("Now skipping %d cells\n", skipTheseCells.Len())

			// add cell to formula
			formula.AddCell(bestCell)

			info.Printf("Just added a new cell (%s) to the formula\n", bestCell)
		}

		fmt.Println("#############")

		// set cover in index
		formula.CoverIndex(&relation)

		info.Println("Formula")
		fmt.Println(formula.cells)
		info.Println("Relation")
		fmt.Println(relation)

		var values []Value
		for _, cell := range formula.cells {
			value := Value{cell.attribute.attributeType, cell.attribute.attributeName, cell.value}
			values = append(values, value)
		}
		summary = append(summary, values)
	}

	return summary
}
