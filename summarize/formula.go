package summarize

import (
	"bytes"
	"fmt"

	"golang.org/x/tools/container/intsets"
)

// TupleCovers gives us the value for each tuple
type TupleCovers map[int]float64

// Formula is a map from attribute id to lists of cells
type Formula struct {
	cells                []Cell         // list of cells
	usedSingleAttributes intsets.Sparse // which single attributes are already used
	tupleCover           TupleCovers    // how much does a tuple contribute to the formula
	cover                float64        // how much does this formula cover, sum of valid tupleCover
}

// NewFormula creates a new formula from a cell
func NewFormula(cell Cell) *Formula {
	var formula Formula

	formula.tupleCover = make(TupleCovers)
	formula.cover = 0

	for tuple, cover := range cell.covers {
		if !cover.covered {
			formula.tupleCover[tuple] = cover.weight
			formula.cover += cover.weight
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
		if cover, has := cell.covers[tuple]; has {
			formula.cover += cover.weight
			formula.tupleCover[tuple] += cover.weight
		} else {
			formula.cover -= formula.tupleCover[tuple]
			delete(formula.tupleCover, tuple)
		}
	}
}

// CoverIndex updates the cover so that in the next iteration the same tuples are not covered again
func (formula *Formula) CoverIndex(relation *RelationIndex) {
	// TODO: is other direction faster?
	for _, cell := range formula.cells {
		for tuple := range formula.tupleCover {
			if cover, has := cell.covers[tuple]; has && !cover.covered {
				// set uncovered to covered
				cell.covers[tuple].covered = true
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
