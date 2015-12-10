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
	for _, attr := range relation.attrs {
		for i := range attr.cells {
			cell := &attr.cells[i]
			potential := len(cell.cover)
			// todo; we may be able to ignore if we add regularization
			rankedCell := RankedCell{cell, potential}
			rankedCells = append(rankedCells, rankedCell)
		}
	}
	return rankedCells
}

// returns the best cell form a list of cells with potentials
// requires that the cells are a sorted heap
func updateBestCellHeap(cellHeap *CellHeap) bool {
	bestCover := 0

	for cellHeap.Peek().potential > bestCover {
		cell := cellHeap.Peek()

		fmt.Print(cell.cell, "From ", cell.potential)
		cover := cell.recomputeCoverage()
		fmt.Println(" to ", cover, cell.potential)

		heap.Fix(cellHeap, 0)

		if cover > bestCover {
			bestCover = cover
		}
	}

	return bestCover != 0
}

// returns nil if no cell could be found that improves the formula
// requires cells to be a heap
func updateFormulaBestCellHeap(formulaCellHeap *CellHeap, formula *Formula) bool {
	// diff of the coverage of the whole formula
	bestDiff := 0

	// cover of a single cell
	bestCover := 0

	for len(*formulaCellHeap) > 0 && formulaCellHeap.Peek().potential > bestCover {
		cell := formulaCellHeap.Peek()
		if cell.cell.attribute.attributeType == single && formula.usedSingleAttributes.Contains(cell.cell.attribute.index) {
			// the formula already has a value assigned to this attribute
			dbg.Println("Ignoring single attribute cell", cell)
			heap.Pop(formulaCellHeap)
			continue
		}

		cover, coverDiff := cell.recomputeFormulaCoverage(formula)

		heap.Fix(formulaCellHeap, 0)

		if cover > bestCover {
			// update cover so that we can escape early
			bestCover = cover
		}

		if coverDiff > bestDiff {
			bestDiff = coverDiff
		}
	}

	if bestDiff <= 0 || len(*formulaCellHeap) == 0 {
		return false
	}

	return true
}

// Summarize summarizes
func (relation RelationIndex) Summarize(size int) Summary {
	fmt.Println(relation)

	var formulaCover []int
	summaryCover := 0
	var summary Summary

	rankedCells := makeRankedCells(relation)
	heap.Init(&rankedCells)

	dbg.Println("Initial ranking")
	fmt.Println(rankedCells)

	for len(summary) < size {
		fmt.Println("==============================")

		// add new formula with best cell
		goodFormula := updateBestCellHeap(&rankedCells)

		if !goodFormula {
			info.Println("Adding a formula doesn't help. Let's stop right here.")
			break
		}

		// make a copy of the ranked cells, we can use this now in the context of a formula and remove elements and reorder
		// note that CellHeap has pointers so we can safely modify the slice but not the cells it points to
		formulaRankedCells := make(CellHeap, len(rankedCells))
		copy(formulaRankedCells, rankedCells)

		cell := heap.Pop(&formulaRankedCells).(RankedCell)
		formula := NewFormula(*cell.cell)

		dbg.Printf("Just added a new formula with cell %s\n", cell.cell)
		fmt.Println(formula.tupleCover)

		dbg.Println("Popped. Ranking is now")
		fmt.Println(formulaRankedCells)

		// keep adding to formula
		for true {
			fmt.Println("--------------------------------")

			improved := updateFormulaBestCellHeap(&formulaRankedCells, formula)

			if !improved {
				info.Println("Adding a new cell to the formula doesn't help.")
				break
			}

			// add cell to formula
			cell := heap.Pop(&formulaRankedCells).(RankedCell)
			formula.AddCell(*cell.cell)

			info.Printf("Just added a new cell (%s) to the formula\n", cell)

			dbg.Println("Popped. Ranking is now")
			fmt.Println(formulaRankedCells)
		}

		fmt.Println("#############")

		// set cover in index
		formula.CoverIndex(&relation)

		info.Println("Formula")
		fmt.Println(formula.cells)

		info.Println("Relation")
		fmt.Println(relation)

		cover := 0
		for _, tupleCover := range formula.tupleCover {
			cover += tupleCover
		}
		formulaCover = append(formulaCover, cover)
		summaryCover += cover

		var values []Value
		for _, cell := range formula.cells {
			value := Value{cell.attribute.attributeType, cell.attribute.attributeName, cell.value}
			values = append(values, value)
		}
		summary = append(summary, values)
	}

	fmt.Printf("Formulas cover %v = %d\n", formulaCover, summaryCover)
	return summary
}
