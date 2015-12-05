package summary

import (
	"bytes"
	"strings"
)

// Relation is a slice of tuples
type Relation []Tuple

// NewRelationFromString creates a relation from a string
func NewRelationFromString(description string) (Relation, error) {
	relation := make(Relation, 0)

	lines := strings.Split(description, "\n")

	attributes := strings.Split(lines[0], ",")

	for _, line := range lines[1:] {
		tuple, err := NewTupleFromString(line, attributes)
		if err != nil {
			return nil, err
		}
		relation = append(relation, tuple)
	}

	return relation, nil
}

// DebugString prints a relation
func (relation *Relation) DebugString() string {
	var buffer bytes.Buffer

	for _, t := range *relation {
		buffer.WriteString(t.DebugString())
		buffer.WriteString("\n")
	}

	return buffer.String()
}
