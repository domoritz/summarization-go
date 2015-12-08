package summarize

import (
	"errors"
	"fmt"
	"strings"
)

// NewIndexFromString creates a relation index from a string
func NewIndexFromString(description string) (*RelationIndex, error) {
	lines := strings.Split(description, "\n")

	typeNames := strings.Split(lines[0], ",")
	names := strings.Split(lines[1], ",")

	if len(names) != len(typeNames) {
		err := fmt.Sprintf("Mismatching number of names and types. %d != %d", len(names), len(typeNames))
		return nil, errors.New(err)
	}

	numAttr := len(typeNames)

	index := make([]Attribute, numAttr)

	for i, typeName := range typeNames {
		typeName = strings.TrimSpace(typeName)
		switch typeName {
		case single.String():
			index[i].attributeType = single
		case set.String():
			index[i].attributeType = set
		case hierarchy.String():
			index[i].attributeType = hierarchy
		}
		index[i].name = names[i]
		index[i].tuples = make(map[string]TupleCover)
	}

	numTuples := len(lines[2:])

	for tuple, line := range lines[2:] {
		values := strings.Split(line, ",")
		if len(values) != numAttr {
			err := fmt.Sprintf("Wrong number of attributes. Expected %d but got %d.", numAttr, len(values))
			return nil, errors.New(err)
		}

		for i, value := range values {
			attr := index[i]
			value = strings.TrimSpace(value)

			if len(value) == 0 {
				// null
				continue
			}

			switch attr.attributeType {
			case single:
				if tc, has := attr.tuples[value]; !has {
					c := make(TupleCover)
					c[tuple] = false
					attr.tuples[value] = c
				} else {
					tc[tuple] = false
				}
			case set:
				setValues := strings.Split(value, " ")
				for _, setValue := range setValues {
					if tc, has := attr.tuples[setValue]; !has {
						c := make(TupleCover)
						c[tuple] = false
						attr.tuples[setValue] = c
					} else {
						tc[tuple] = false
					}
				}
			case hierarchy:
				prefix := ""
				hValues := strings.Split(value, " ")
				for _, hValue := range hValues {
					prefix += hValue
					if tc, has := attr.tuples[prefix]; !has {
						c := make(TupleCover)
						c[tuple] = false
						attr.tuples[prefix] = c
					} else {
						tc[tuple] = false
					}
				}
			}
		}
	}

	return &RelationIndex{index, numTuples}, nil
}
