package summarize

import (
	"bytes"
	"fmt"
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

type cell map[int]bool

type attribute struct {
	attributeType Type
	name          string
	coveredTuples map[string]cell // TODO: make slice
}

type rankedAttribute struct {
	attribute
	coverage int
}

// RelationIndex is an inverted index
type RelationIndex []attribute

// Summary is a summary
type Summary [][]attribute

// Summarize summarizes
func (relation RelationIndex) Summarize(size int) Summary {

	var summary Summary

	// how much does a tuple cover
	// var tupleCovers []int

	return summary
}

func (relation RelationIndex) String() string {
	var buffer bytes.Buffer
	for _, attribute := range relation {
		buffer.WriteString(fmt.Sprintf("Attribute %s (%s):\n", attribute.name, attribute.attributeType))
		for value, cell := range attribute.coveredTuples {
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
