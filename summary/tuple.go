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

// Coverage true if the tuple satisfies the formula (formula is subset) in addition to the size of the cover.
func (tuple *Tuple) Coverage(formula *Tuple) (bool, int) {
	coverage := 0

	for i, attr := range tuple.Single {
		if formula.Single[i] == nil {
			continue
		}
		if attr == nil {
			return false, -1
		}

		if !formula.Single[i].Equal(attr) {
			return false, -1
		}
		if !attr.covered {
			coverage++
		}
	}

	for i, attr := range tuple.Set {
		if formula.Set[i] == nil {
			continue
		}
		if attr == nil {
			return false, -1
		}

		if ok, cover := formula.Set[i].SubsetCover(attr); ok {
			coverage += cover
		} else {
			return false, -1
		}
	}

	for i, attr := range tuple.Hierarchy {
		if formula.Hierarchy[i] == nil {
			continue
		}
		if attr == nil {
			return false, -1
		}

		if !formula.Hierarchy[i].Prefix(attr) {
			return false, -1
		}
		coverage += len(formula.Hierarchy[i].hierarchy)
	}

	return true, coverage
}

// NewTupleFromCell creates a new tuple with only one cell
func NewTupleFromCell(cell Cell, sizes Sizes) Tuple {
	singles := make([]*SingleValueAttribute, sizes.single)
	sets := make([]*SetAttribute, sizes.set)
	hierarchies := make([]*HierarchyAttribute, sizes.hierarchy)
	tuple := Tuple{singles, sets, hierarchies}

	switch cell.Type {
	case single:
		a := NewSingle(cell.Value)
		tuple.Single[cell.Attribute] = &a
	case set:
		a := NewSet(Set{cell.Value: false})
		tuple.Set[cell.Attribute] = &a
	case hierarchy:
		// TODO
	}

	return tuple
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
				setValue[v] = false
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

// AddCells adds all cells from the tuple to the map
func (tuple *Tuple) AddCells(cells *map[CellKey]*Cell) {
	for i, attr := range tuple.Single {
		if attr != nil {
			key := CellKey{single, i, attr.value}
			cell, ok := (*cells)[key]
			if ok {
				// increase potential
				cell.Potential++
			} else {
				// add new cell
				(*cells)[key] = &Cell{key, 1}
			}
		}
	}

	for i, attr := range tuple.Set {
		if attr != nil {
			for value := range attr.values {
				key := CellKey{set, i, value}
				cell, ok := (*cells)[key]
				if ok {
					// increase potential
					cell.Potential++
				} else {
					// add new cell
					(*cells)[key] = &Cell{key, 1}
				}
			}
		}
	}

	// TODO: hierarchy
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
