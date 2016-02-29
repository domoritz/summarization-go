package summarize

import (
	"bytes"
	"fmt"
)

type Cover struct {
	covered bool    // whether the cell covers the cell in this tuple
	weight  float64 // the cover weight
}

// TupleCover is a map from tuple index to whether the tuple covers it
type TupleCover map[int]*Cover

// Cell is an attribute value that covers tuples
type Cell struct {
	covers       TupleCover // what cells the attribute covers
	attribute    *Attribute // attribute
	value        string     // attribute value
	equalWeights bool       // the cover weights of this cell
}

// MakeCell makes a new cell
func MakeCell(attr *Attribute, value string, equalWeights bool) Cell {
	covers := make(TupleCover)
	cell := Cell{covers, attr, value, equalWeights}
	return cell
}

// Weight computes the sum of weights for all covered cells
func (cell *Cell) SumWeights() float64 {
	// shortcut since all weights are unit 1
	if cell.equalWeights {
		return float64(len(cell.covers))
	}

	// sum up the weights in covers
	s := 0.0
	for _, cover := range cell.covers {
		s += cover.weight
	}
	return s
}

func (cell Cell) String() string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "Attr %s: %s", cell.attribute.attributeName, cell.value)
	return buffer.String()
}
