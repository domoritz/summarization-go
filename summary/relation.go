package summary

import (
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Relation is a slice of tuples
type Relation struct {
	Tuples         []Tuple
	AttributeNames []string
}

// NewRelationFromString creates a relation from a string
func NewRelationFromString(description string) (*Relation, error) {
	var tuples []Tuple

	lines := strings.Split(description, "\n")

	attributes := strings.Split(lines[0], ",")

	for _, line := range lines[1:] {
		tuple, err := NewTupleFromString(line, len(attributes))
		if err != nil {
			return nil, err
		}
		tuples = append(tuples, tuple)
	}

	relation := Relation{tuples, attributes}

	return &relation, nil
}

// PrintDebugString prints a relation
func (relation *Relation) PrintDebugString() {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(relation.AttributeNames)

	for _, tuple := range relation.Tuples {
		values := make([]string, len(relation.Tuples))
		for i, attr := range tuple {
			values[i] = attr.DebugString()
		}
		table.Append(values)
	}

	table.Render()
}
