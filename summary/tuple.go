package summary

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Tuple is a map of attributes
type Tuple struct {
	Single    []*SingleValueAttribute
	Set       []*SetAttribute
	Hierarchy []*HierarchyAttribute
}

// Satisfies is true if the tuple satisfies the formula (formula is subset)
func (tuple *Tuple) Satisfies(formula *Tuple) bool {
	for i, attr := range tuple.Single {
		if formula.Single[i] == nil {
			return true
		}
		if attr == nil {
			return false
		}

		if !formula.Single[i].Equal(attr) {
			return false
		}
	}

	for i, attr := range tuple.Set {
		if formula.Set[i] == nil {
			return true
		}
		if attr == nil {
			return false
		}

		if !formula.Set[i].Subset(attr) {
			return false
		}
	}

	for i, attr := range tuple.Hierarchy {
		if formula.Hierarchy[i] == nil {
			return true
		}
		if attr == nil {
			return false
		}

		if !formula.Hierarchy[i].Prefix(attr) {
			return false
		}
	}

	return true
}

// Size returns the size of the tuple
func (tuple *Tuple) Size() int {
	length := len(tuple.Single)

	for _, attr := range tuple.Set {
		length += len(attr.values)
	}

	for _, attr := range tuple.Hierarchy {
		length += len(attr.hierarchy)
	}

	return length
}

// NewTupleFromString creates a tuple from a string
func NewTupleFromString(description string, types []Type) (Tuple, error) {
	tuple := Tuple{}

	values := strings.Split(description, ",")
	if len(values) != len(types) {
		err := fmt.Sprintf("Wrong number of attributes. Expected %d but got %d.", len(types), len(values))
		return tuple, errors.New(err)
	}

	for i, value := range values {
		value = strings.TrimSpace(value)
		switch types[i] {
		case set:
			if len(value) == 0 {
				// insert null
				tuple.Set = append(tuple.Set, nil)
				break
			}

			setValues := strings.Split(value, " ")

			setValue := make(map[string]bool)
			for _, v := range setValues {
				setValue[v] = true
			}
			a := NewSet(setValue)
			tuple.Set = append(tuple.Set, &a)
		case single:
			if len(value) == 0 {
				// insert null
				tuple.Single = append(tuple.Single, nil)
				break
			}
			a := NewSingle(value)
			tuple.Single = append(tuple.Single, &a)
		case hierarchy:
			if len(value) == 0 {
				// insert null
				tuple.Hierarchy = append(tuple.Hierarchy, nil)
				break
			}
			a := NewHierachy(strings.Split(value, " "))
			tuple.Hierarchy = append(tuple.Hierarchy, &a)
		}
	}

	return tuple, nil
}

// GetDebugStrings returns a list of values in order of the types
func (tuple *Tuple) GetDebugStrings(types []Type) []string {
	values := make([]string, len(types))

	singleIndex := 0
	setIndex := 0
	hierarchyIndex := 0

	for i, t := range types {
		switch t {
		case single:
			values[i] = tuple.Single[singleIndex].DebugString()
			singleIndex++
		case set:
			values[i] = tuple.Set[setIndex].DebugString()
			setIndex++
		case hierarchy:
			values[i] = tuple.Hierarchy[hierarchyIndex].DebugString()
			hierarchyIndex++
		}
	}

	return values
}

// DebugString prints a tuple without attribute names
func (tuple Tuple) DebugString() string {
	var buffer bytes.Buffer

	for i := 0; i < len(tuple.Single); i++ {
		buffer.WriteString(tuple.Single[i].DebugString())
		buffer.WriteString(" ")
	}

	for i := 0; i < len(tuple.Set); i++ {
		buffer.WriteString(tuple.Set[i].DebugString())
		buffer.WriteString(" ")
	}

	for i := 0; i < len(tuple.Hierarchy); i++ {
		buffer.WriteString(tuple.Hierarchy[i].DebugString())
		buffer.WriteString(" ")
	}

	return buffer.String()
}
