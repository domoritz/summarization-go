package summary

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// Type is the attribute type
type Type int

// Sizes gives us the size of the attribute lists
type Sizes struct {
	single    int
	set       int
	hierarchy int
}

// Relation is a slice of tuples
type Relation struct {
	Tuples         []Tuple
	AttributeNames []string
	AttributeTypes []Type
	attributeSizes Sizes
}

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

// GetSizes gets the attribute sizes of the relation used to initialize tuples
func (relation *Relation) GetSizes() Sizes {
	return relation.attributeSizes
}

// NewRelationFromString creates a relation from a string
func NewRelationFromString(description string) (*Relation, error) {
	var tuples []Tuple

	lines := strings.Split(description, "\n")
	typeNames := strings.Split(lines[0], ",")

	types := make([]Type, len(typeNames))
	sizes := Sizes{}

	for i, name := range typeNames {
		name = strings.TrimSpace(name)
		switch name {
		case single.String():
			types[i] = single
			sizes.single++
		case set.String():
			types[i] = set
			sizes.set++
		case hierarchy.String():
			types[i] = hierarchy
			sizes.hierarchy++
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

	relation := Relation{tuples, names, types, sizes}

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
		names[i] = fmt.Sprintf("%s (%s)", relation.AttributeNames[i], relation.AttributeTypes[i])
	}

	table.SetHeader(names)

	for _, tuple := range relation.Tuples {
		values := tuple.GetDebugStrings(relation.AttributeTypes)
		table.Append(values)
	}

	table.Render()
}
