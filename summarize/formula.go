package summarize

import (
	"bytes"
	"fmt"

	"golang.org/x/tools/container/intsets"
)

// TupleCovers gives us the value for each tuple
type TupleCovers map[int]int

// Formula is a map from attribute id to lists of cells
type Formula struct {
	cells                []Cell         // list of cells
	tupleCover           TupleCovers    // how much does a tuple contribute to the formula
	usedSingleAttributes intsets.Sparse // which single attributes are already used
}

// NewFormula creates a new formula from a cell
func NewFormula(cell Cell) *Formula {
	var formula Formula

	formula.tupleCover = make(TupleCovers)

	for tuple, covered := range cell.cover {
		if !covered {
			formula.tupleCover[tuple] = 1
		} else {
			formula.tupleCover[tuple] = 0
		}
	}

	formula.addCellNoUpdateValues(cell)

	return &formula
}

func (formula *Formula) addCellNoUpdateValues(cell Cell) {
	formula.cells = append(formula.cells, cell)

	// if the cell is a single, add attribute to exclude list
	if cell.attribute.attributeType == single {
		formula.usedSingleAttributes.Insert(cell.attribute.index)
	}
}

// AddCell adds a cell to the formula and updates internals
func (formula *Formula) AddCell(cell Cell) {
	formula.addCellNoUpdateValues(cell)

	// TODO: is other direction faster?
	for tuple := range formula.tupleCover {
		if _, has := cell.cover[tuple]; !has {
			delete(formula.tupleCover, tuple)
		} else {
			formula.tupleCover[tuple]++
		}
	}
}

// CoverIndex updates the cover so that in the next iteration the same tuples are not covered again
func (formula *Formula) CoverIndex(relation *RelationIndex) {
	// TODO: is other direction faster?
	for _, cell := range formula.cells {
		for tuple := range formula.tupleCover {
			if covered, has := cell.cover[tuple]; has && !covered {
				// set uncovered to covered
				cell.cover[tuple] = true
			}
		}
	}
}

func (values TupleCovers) String() string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "Tuple Value (%d):\n", len(values))
	for i, values := range values {
		fmt.Fprintf(&buffer, "%d: %d\n", i, values)
	}
	return buffer.String()
}
