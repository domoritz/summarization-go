package summarize

import (
	"container/heap"
	"fmt"
	"log"
	"os"
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

func makeRankedCells(relation RelationIndex) CellHeap {
	rankedCells := make(CellHeap, 0, relation.numValues)
	uid := 0
	for _, attr := range relation.attrs {
		for _, cell := range attr.cells {
			rankedCell := RankedCell{&cell, len(cell.cover)}
			rankedCells = append(rankedCells, rankedCell)
			uid++
		}
	}
	return rankedCells
}

// returns the best cell form a list of cells with potentials
// requires that the cells are a sorted heap
func updateBestCellHeap(cellHeap *CellHeap) {
	bestValue := 0

	cell := (*cellHeap)[0]
	for cell.potential > bestValue {
		value := cell.recomputeCoverage()

		heap.Fix(cellHeap, 0)

		if value > bestValue {
			bestValue = value
		}
		cell = (*cellHeap)[0]
	}
}

// returns nil if no cell could be found that improves the formula
// requires cells to be a heap
func updateFormulaBestCellHeap(formulaCellHeap *CellHeap, formula *Formula) bool {
	// the best improvement in coverage for any cell
	bestDiff := 0
	bestValue := 0

	cell := (*formulaCellHeap)[0]
	for cell.potential > bestValue {
		if cell.cell.attribute.attributeType == single && formula.usedSingleAttributes.Contains(cell.cell.attribute.index) {
			// the formula already has a value assigned to this attribute
			dbg.Println("Ignoring single attribute cell", cell)
			formulaCellHeap.Pop()
			continue
		}

		value, valueDiff := cell.recomputeFormulaCoverage(formula)

		if cell.potential+valueDiff != value {
			fmt.Printf("%d + %d != %d", cell.potential, valueDiff, value)
		}

		heap.Fix(formulaCellHeap, 0)

		if value > bestValue {
			bestValue = value
		}
		cell = (*formulaCellHeap)[0]
	}

	if bestDiff <= 0 {
		info.Println("Looks like we cannot find a cell that should be added")
		return false
	}

	return true
}

// Summarize summarizes
func (relation RelationIndex) Summarize(size int) Summary {
	fmt.Println(relation)

	var summary Summary

	rankedCells := makeRankedCells(relation)
	heap.Init(&rankedCells)

	dbg.Println("Initial ranking")
	fmt.Println(rankedCells)

	for len(summary) < size {
		fmt.Println("==============================")

		// add new formula with best cell
		updateBestCellHeap(&rankedCells)

		if rankedCells[0].potential <= 0 {
			info.Println("Adding a new cell to the formula doesn't help. Let's stop right here.")
			break
		}

		// make a copy of the ranked cells, we can use this now in the context of a formula and remove elements and reorder
		// note that CellHeap has pointers so we can safely modify the slice but not the cells it points to
		formulaRankedCells := make(CellHeap, len(rankedCells))
		copy(formulaRankedCells, rankedCells)

		formula := NewFormula(*formulaRankedCells[0].cell)
		formulaRankedCells.Pop()

		dbg.Printf("Just added a new formula with cell %s\n", *formulaRankedCells[0].cell)
		fmt.Println(formula.tupleValue)

		// keep adding to formula
		for true {
			fmt.Println("--------------------------------")

			updateFormulaBestCellHeap(&formulaRankedCells, formula)

			if rankedCells[0].potential <= 0 {
				break
			}

			// add cell to formula
			formula.AddCell(*formulaRankedCells[0].cell)

			info.Printf("Just added a new cell (%s) to the formula\n", *formulaRankedCells[0].cell)
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
