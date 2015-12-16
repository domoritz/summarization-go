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
			rankedCell := RankedCell{cell, potential, potential}
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

	return bestCover > 0
}

// returns nil if no cell could be found that improves the formula
// requires cells to be a heap
func updateFormulaBestCellHeap(formulaCellHeap *CellHeap, formula *Formula) bool {
	// the largest change that a cell can do
	bestCover := 0

	for len(*formulaCellHeap) > 0 && formulaCellHeap.Peek().potential > bestCover {
		cell := formulaCellHeap.Peek()
		if cell.cell.attribute.attributeType == single && formula.usedSingleAttributes.Has(cell.cell.attribute.index) {
			// the formula already has a value assigned to this attribute
			heap.Pop(formulaCellHeap)
			continue
		}

		cellCover := cell.recomputeFormulaCoverage(formula)

		if cell.maxPotential <= 0 {
			// looks like there is no overlap between what tuples the formula and the cell cover
			heap.Pop(formulaCellHeap)
			continue
		}

		if cellCover > bestCover {
			bestCover = cellCover
		}
	}

	if bestCover <= 0 || len(*formulaCellHeap) == 0 {
		return false
	}

	fmt.Println("we should add", formulaCellHeap.Peek(), bestCover)

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
			} else {
				// have to reset the potentials because we will reduce the set of tuples that the formula covers
				for i := range formulaRankedCells {
					formulaRankedCells[i].potential = cell.maxPotential
				}
				heap.Init(&formulaRankedCells)
			}

			// add cell to formula
			cell := heap.Pop(&formulaRankedCells).(RankedCell)
			formula.AddCell(*cell.cell)
		}

		// set cover in index
		formula.CoverIndex(&relation)

		formulaCover = append(formulaCover, formula.cover)
		summaryCover += formula.cover

		// if the formula has only one cell, we can pop that one off the heap beacuse nothing can every use it again
		if len(formula.cells) == 1 {
			if rankedCells.Peek().cell.value != formula.cells[0].value {
				panic("assert")
			}
			heap.Pop(&rankedCells)
		}

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
	names = append(names, "#")
	names = append(names, "count")

	table.SetHeader(names)
	table.SetAutoWrapText(false)

	for i, name := range names[0:len(header)] {
		header[name] = i
	}

	for i, cells := range summary {
		values := make([]string, len(names))
		values[len(values)-2] = fmt.Sprintf("%d", i)
		for _, cell := range cells {
			key := fmt.Sprintf("%s (%s)", cell.attributeName, cell.attributeType)
			prefix := ""
			if len(values[header[key]]) > 0 {
				prefix = ", "
			}
			values[header[key]] += prefix + cell.value
		}
		values[len(values)-1] = fmt.Sprintf("%d", len(cells))
		table.Append(values)
	}

	table.Render()
}
