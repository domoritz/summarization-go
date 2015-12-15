package summarize

import (
	"bytes"
	"fmt"
)

// TupleCover is a map from tuple index to whether the tuple covers it
type TupleCover map[int]bool

// Cell is a cell
type Cell struct {
	cover     TupleCover // what the attribute covers
	attribute *Attribute // attribute
	value     string     // attribute value
}

// MakeCell makes a new cell
func MakeCell(attr *Attribute, value string) Cell {
	cover := make(TupleCover)
	cell := Cell{cover, attr, value}
	return cell
}

func (cell Cell) String() string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "Attr %s: %s", cell.attribute.attributeName, cell.value)
	return buffer.String()
}
