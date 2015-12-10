package summarize

import (
	"bytes"
	"fmt"
)

// TupleCover is a map from tuple index to how much the tuple contributes
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
	buffer.WriteString(fmt.Sprintf("Attr %s: %s", cell.attribute.attributeName, cell.value))
	return buffer.String()
}
