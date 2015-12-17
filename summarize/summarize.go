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

// SummaryResult packs a summary
type SummaryResult struct {
	Summary      Summary // the summary
	FormulaCover []int   // how much each formula covers
	SummaryCover int     // sum of tupleCover
}

func makeRankedCells(relation RelationIndex) CellHeap {
	var rankedCells CellHeap
	index := 0
	for _, attr := range relation.attrs {
		for i := range attr.cells {
			cell := &attr.cells[i]
			potential := len(cell.cover)
			// todo; we may be able to ignore if we add regularization
			rankedCell := RankedCell{cell, potential, potential, index}
			rankedCells = append(rankedCells, &rankedCell)
			index++
		}
	}
	return rankedCells
}

// makes copies of pointer
func copyRankedCells(rankedCells CellHeap) CellHeap {
	cellsCopy := make(CellHeap, len(rankedCells))
	for i, cell := range rankedCells {
		cellCopy := new(RankedCell)
		*cellCopy = *cell
		cellsCopy[i] = cellCopy
	}
	return cellsCopy
}

// returns the best cell form a list of cells with potentials
// requires that the cells are a sorted heap
func updateBestCellHeap(cellHeap *CellHeap) (bool, *RankedCell) {
	bestCover := 0
	var bestCell *RankedCell

	for len(*cellHeap) > 0 && cellHeap.Peek().potential > bestCover {
		cell := cellHeap.Peek()
		cover := cell.recomputeCoverage()
		heap.Fix(cellHeap, cell.index)

		if cover > bestCover {
			bestCover = cover
			bestCell = cell
		}
	}

	return bestCover > 0 && len(*cellHeap) > 0, bestCell
}

// returns nil if no cell could be found that improves the formula
// requires cells to be a heap
func updateFormulaBestCellHeap(formulaCellHeap *CellHeap, formula *Formula) (bool, *RankedCell) {
	// the largest change that a cell can do
	bestCover := 0
	var bestCell *RankedCell

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
			// this means we can remove it because this cell will not be usable for this formula
			heap.Pop(formulaCellHeap)
			continue
		}

		if cellCover > bestCover {
			bestCover = cellCover
			bestCell = cell
		}

		heap.Fix(formulaCellHeap, cell.index)
	}

	return bestCover > 0 && len(*formulaCellHeap) > 0, bestCell
}

// Summarize summarizes
func (relation RelationIndex) Summarize(size int) SummaryResult {
	var formulaCover []int
	summaryCover := 0
	var summary Summary

	rankedCells := makeRankedCells(relation)
	heap.Init(&rankedCells)

	for len(summary) < size {
		// add new formula with best cell
		goodFormula, cell := updateBestCellHeap(&rankedCells)

		if !goodFormula {
			break
		}

		formula := NewFormula(*cell.cell)

		// make a copy of the ranked cells, we can use this now in the context of a formula and remove elements and reorder
		// note that CellHeap has pointers so we can safely modify the slice but not the cells it points to
		formulaRankedCells := copyRankedCells(rankedCells)

		heap.Remove(&formulaRankedCells, cell.index)

		// keep adding to formula
		for true {
			improved, cell := updateFormulaBestCellHeap(&formulaRankedCells, formula)

			if !improved {
				break
			}

			// add cell to formula
			formula.AddCell(*cell.cell)

			// remove the cell from the heap beacuse we used it in this formula
			heap.Remove(&formulaRankedCells, cell.index)

			// have to reset the potentials because we will reduce the set of tuples that the formula covers
			for i := range formulaRankedCells {
				formulaRankedCells[i].potential = formulaRankedCells[i].maxPotential
			}
			heap.Init(&formulaRankedCells)
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

	return SummaryResult{
		summary,
		formulaCover,
		summaryCover,
	}
}

// DebugPrint prints a summary
func (summary SummaryResult) DebugPrint() {
	fmt.Printf("Summary (covers %d attributes):\n", summary.SummaryCover)

	table := tablewriter.NewWriter(os.Stdout)

	// provides positions
	header := make(map[string]int)

	for _, cells := range summary.Summary {
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
	names = append(names, "cover")
	names = append(names, "# cells")

	table.SetHeader(names)
	table.SetColWidth(100)

	for i, name := range names[0:len(header)] {
		header[name] = i
	}

	for i, cells := range summary.Summary {
		values := make([]string, len(names))
		values[len(values)-2] = fmt.Sprintf("%d", summary.FormulaCover[i])
		for _, cell := range cells {
			key := fmt.Sprintf("%s (%s)", cell.attributeName, cell.attributeType)

			switch cell.attributeType {
			case set:
				prefix := ""
				if len(values[header[key]]) > 0 {
					prefix = ", "
				}
				values[header[key]] += prefix + cell.value
			case hierarchy:
				if len(values[header[key]]) < len(cell.value) {
					values[header[key]] = cell.value
				}
			case single:
				values[header[key]] = cell.value
			}

		}
		values[len(values)-1] = fmt.Sprintf("%d", len(cells))
		table.Append(values)
	}

	table.Render()
}
