package summarize

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Type is the attribute type
type Type int

const (
	single Type = iota
	set
	hierarchy
)

func (t Type) String() string {
	switch t {
	case single:
		return "single"
	case set:
		return "set"
	case hierarchy:
		return "hierarchy"
	default:
		return "unknown"
	}
}

// TupleCover is a map from tuple index to whether it is covered or not
type TupleCover map[int]bool

// Attribute is an attribute
type Attribute struct {
	attributeType Type
	name          string
	tuples        map[string]TupleCover // TODO: make slice
}

// Cell is a cell
type Cell struct {
	attribute       *Attribute
	value           string
	coverage        int // coverage in the whole relation
	formulaCoverage int // coverage inside the tuple slice
}

type cellSlice []*Cell

// Len is part of sort.Interface.
func (cells cellSlice) Len() int {
	return len(cells)
}

// Swap is part of sort.Interface.
func (cells cellSlice) Swap(i, j int) {
	cells[i], cells[j] = cells[j], cells[i]
}

// Less is part of sort.Interface. Sort by Potential.
func (cells cellSlice) Less(i, j int) bool {
	return cells[i].coverage > cells[j].coverage
}

// RelationIndex is an inverted index
type RelationIndex struct {
	attrs     []Attribute
	numTuples int
}

type tupleCover []int

// Summary is a summary
type Summary [][]Cell

// Summarize summarizes
func (relation RelationIndex) Summarize(size int) Summary {
	var summary Summary

	rankedCells := make(cellSlice, 0)
	for i, attr := range relation.attrs {
		for value, cover := range attr.tuples {
			cell := Cell{&relation.attrs[i], value, len(cover), -1}
			rankedCells = append(rankedCells, &cell)
		}
	}

	sort.Sort(rankedCells)
	fmt.Println(rankedCells)

	// how much does the current formula contribute to the coverage
	tupleCover := make(tupleCover, relation.numTuples)

	for len(summary) < size {
		// add new formula with best cell
		cell := rankedCells[0]
		formula := []Cell{*cell}
		summary = append(summary, formula)

		for tuple, covered := range cell.attribute.tuples[cell.value] {
			if !covered {
				tupleCover[tuple]++
			}
		}

		fmt.Println(tupleCover)

		// keep adding to formula
		// for true {
		// 	var bestCell *Cell
		// 	for _, cell := range rankedCells {
		// 		if bestCell {
		//
		// 		}
		// 	}
		// }

		break
	}

	return summary
}

func (cells cellSlice) String() string {
	var buffer bytes.Buffer
	for _, cell := range cells {
		buffer.WriteString(fmt.Sprintf("%s\n", cell))
	}
	return buffer.String()
}

func (cover tupleCover) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("Cover:\n")
	for i, cover := range cover {
		buffer.WriteString(fmt.Sprintf("%d: %d\n", i, cover))
	}
	return buffer.String()
}

func (cell Cell) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Attr %s: %s (%d)", cell.attribute.name, cell.value, cell.coverage))
	return buffer.String()
}

func (relation RelationIndex) String() string {
	var buffer bytes.Buffer
	for _, attribute := range relation.attrs {
		buffer.WriteString(fmt.Sprintf("Attribute %s (%s):\n", attribute.name, attribute.attributeType))
		for value, cell := range attribute.tuples {
			buffer.WriteString(fmt.Sprintf("Value %s covers tuples: [", value))
			var tuples []string
			for tuple := range cell {
				tuples = append(tuples, fmt.Sprintf("%d", tuple))
			}

			buffer.WriteString(strings.Join(tuples, " "))

			buffer.WriteString("]\n")
		}
	}

	return buffer.String()
}
