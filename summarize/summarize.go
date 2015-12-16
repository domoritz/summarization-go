package summarize

import (
	"container/heap"
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
)

// var info = log.New(os.Stdout, "INFO: ", log.Lshortfile)
// var dbg = log.New(os.Stdout, "DEBUG: ", log.Lshortfile)

// Value is an assignment for the summary
type Value struct {
	attributeType Type   // attribute type
	attributeName string // attribute name
	value         string // value
}

// Summary is a summary
type Summary [][]Value

func makeRankedCells(relation RelationIndex) CellHeap {
	var rankedCells CellHeap
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
		cover := cell.recomputeCoverage()
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
	bestCellCover := 0

	for len(*formulaCellHeap) > 0 && formulaCellHeap.Peek().potential > bestCellCover {
		cell := formulaCellHeap.Peek()
		if cell.cell.attribute.attributeType == single && formula.usedSingleAttributes.Has(cell.cell.attribute.index) {
			// the formula already has a value assigned to this attribute
			heap.Pop(formulaCellHeap)
			continue
		}

		formulaCover, cellCover := cell.recomputeFormulaCoverage(formula)

		heap.Fix(formulaCellHeap, 0)

		if cellCover > bestCellCover {
			// update cover so that we can escape early
			bestCellCover = cellCover
			bestDiff = formulaCover - formula.cover
		}
	}

	if bestDiff <= 0 || len(*formulaCellHeap) == 0 {
		return false
	}

	return true
}

// Summarize summarizes
func (relation RelationIndex) Summarize(size int) Summary {
	var formulaCover []int
	summaryCover := 0
	var summary Summary

	rankedCells := makeRankedCells(relation)
	heap.Init(&rankedCells)

	for len(summary) < size {
		// add new formula with best cell
		goodFormula := updateBestCellHeap(&rankedCells)

		if !goodFormula {
			break
		}

		// make a copy of the ranked cells, we can use this now in the context of a formula and remove elements and reorder
		// note that CellHeap has pointers so we can safely modify the slice but not the cells it points to
		formulaRankedCells := make(CellHeap, len(rankedCells))
		copy(formulaRankedCells, rankedCells)

		cell := heap.Pop(&formulaRankedCells).(RankedCell)
		formula := NewFormula(*cell.cell)

		// keep adding to formula
		for true {
			improved := updateFormulaBestCellHeap(&formulaRankedCells, formula)

			if !improved {
				break
			}

			// add cell to formula
			cell := heap.Pop(&formulaRankedCells).(RankedCell)
			formula.AddCell(*cell.cell)
		}

		// set cover in index
		formula.CoverIndex(&relation)

		formulaCover = append(formulaCover, formula.cover)
		summaryCover += formula.cover

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

// DebugPrint prints a summary
func (summary Summary) DebugPrint() {
	table := tablewriter.NewWriter(os.Stdout)

	// provides positions
	header := make(map[string]int)

	for _, cells := range summary {
		for _, cell := range cells {
			key := fmt.Sprintf("%s (%s)", cell.attributeName, cell.attributeType)
			header[key] = 0
		}
	}

	names := make([]string, 0, len(header))
	for name := range header {
		names = append(names, name)
	}
	sort.Strings(names)

	table.SetHeader(names)

	for i, name := range names {
		header[name] = i
	}

	for _, cells := range summary {
		values := make([]string, len(names))
		for _, cell := range cells {
			key := fmt.Sprintf("%s (%s)", cell.attributeName, cell.attributeType)
			values[header[key]] += cell.value + " "
		}
		table.Append(values)
	}

	table.Render()
}
