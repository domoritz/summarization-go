package summarize

import (
	"bytes"
	"fmt"

	"golang.org/x/tools/container/intsets"
)

// TupleValues gives us the value for each tuple
type TupleValues map[int]int

// Formula is a map from attribute id to lists of cells
type Formula struct {
	cells                CellPointers // list of cells
	tupleValue           TupleValues  // how much does a tuple contribute to the formula
	usedSingleAttributes Set          // which single attributes are already used
	skipTheseCells       intsets.Sparse
}

// NewFormula creates a new formula from a cell
func NewFormula(cell *Cell) *Formula {
	var formula Formula

	formula.usedSingleAttributes = make(Set)
	formula.tupleValue = make(TupleValues)

	for tuple, covered := range *cell.cover {
		if !covered {
			formula.tupleValue[tuple] = 1
		} else {
			formula.tupleValue[tuple] = 0
		}
	}

	formula.addCellNoUpdateValues(cell)

	return &formula
}

func (formula *Formula) addCellNoUpdateValues(cell *Cell) {
	formula.cells = append(formula.cells, cell)

	formula.skipTheseCells.Insert(cell.uid)

	// if the cell is a single, add attribute to exclude list
	if cell.attribute.attributeType == single {
		formula.usedSingleAttributes.Add(cell.attribute.index)
	}
}

// AddCell adds a cell to the formula and updates internals
func (formula *Formula) AddCell(cell *Cell) {
	formula.addCellNoUpdateValues(cell)

	// TODO: is other direction faster?
	for tuple := range formula.tupleValue {
		if _, has := (*cell.cover)[tuple]; !has {
			delete(formula.tupleValue, tuple)
		} else {
			formula.tupleValue[tuple]++
		}
	}
}

// CoverIndex updates the cover so that in the next iteration the same tuples are not covered again
func (formula *Formula) CoverIndex(relation *RelationIndex) {
	// TODO: is other direction faster?
	for _, cell := range formula.cells {
		for tuple := range formula.tupleValue {
			if covered, has := (*cell.cover)[tuple]; has && !covered {
				// set uncovered to covered
				(*cell.cover)[tuple] = true
			}
		}
	}
}

func (values TupleValues) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Tuple Value (%d):\n", len(values)))
	for i, values := range values {
		buffer.WriteString(fmt.Sprintf("%d: %d\n", i, values))
	}
	return buffer.String()
}
