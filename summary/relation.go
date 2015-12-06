package summary

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Type is the attribute type
type Type int

// Relation is a slice of tuples
type Relation struct {
	Tuples         []Tuple
	AttributeNames []string
	AttributeTypes []Type
}

const (
	single Type = iota
	set
	hierarchy
)

// Name gets the type name
func (t Type) Name() string {
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

// NewRelationFromString creates a relation from a string
func NewRelationFromString(description string) (*Relation, error) {
	var tuples []Tuple

	lines := strings.Split(description, "\n")
	typeNames := strings.Split(lines[0], ",")

	types := make([]Type, len(typeNames))
	for i, name := range typeNames {
		name = strings.TrimSpace(name)
		switch name {
		case single.Name():
			types[i] = single
		case set.Name():
			types[i] = set
		case hierarchy.Name():
			types[i] = hierarchy
		}
	}

	names := strings.Split(lines[1], ",")

	for _, line := range lines[2:] {
		tuple, err := NewTupleFromString(line, types)
		if err != nil {
			return nil, err
		}
		tuples = append(tuples, tuple)
	}

	relation := Relation{tuples, names, types}

	return &relation, nil
}

// NumAttributes returns the number of attributes
func (relation *Relation) NumAttributes() int {
	return len(relation.AttributeTypes)
}

// PrintDebugString prints a relation
func (relation *Relation) PrintDebugString() {
	table := tablewriter.NewWriter(os.Stdout)

	names := make([]string, len(relation.AttributeNames))
	for i := 0; i < relation.NumAttributes(); i++ {
		names[i] = fmt.Sprintf("%s (%s)", relation.AttributeNames[i], relation.AttributeTypes[i].Name())
	}

	table.SetHeader(names)

	for _, tuple := range relation.Tuples {
		values := tuple.GetValues(relation.AttributeTypes)
		table.Append(values)
	}

	table.Render()
}
