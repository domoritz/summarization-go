package summarize

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Attribute is an attribute
type Attribute struct {
	attributeType Type                   // attribute type
	name          string                 // attribute name
	index         int                    // a number for this attribute
	tuples        map[string]*TupleCover // TODO: make slice
}

// RelationIndex is an inverted index
type RelationIndex struct {
	attrs     []Attribute
	numTuples int
	numValues int
}

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
		index[i].index = i
		index[i].tuples = make(map[string]*TupleCover)
	}

	numTuples := len(lines[2:])
	numValues := 0

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
					attr.tuples[value] = &c
					numValues++
				} else {
					(*tc)[tuple] = false
				}
			case set:
				setValues := strings.Split(value, " ")
				for _, setValue := range setValues {
					if tc, has := attr.tuples[setValue]; !has {
						c := make(TupleCover)
						c[tuple] = false
						attr.tuples[setValue] = &c
						numValues++
					} else {
						(*tc)[tuple] = false
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
						attr.tuples[prefix] = &c
						numValues++
					} else {
						(*tc)[tuple] = false
					}
				}
			}
		}
	}

	return &RelationIndex{index, numTuples, numValues}, nil
}

func (relation RelationIndex) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Relation Index (%d attributes, %d tuples, %d values):\n", len(relation.attrs), relation.numTuples, relation.numValues))
	for _, attribute := range relation.attrs {
		buffer.WriteString(fmt.Sprintf("Attribute %s (%s):\n", attribute.name, attribute.attributeType))
		for value, cover := range attribute.tuples {
			buffer.WriteString(fmt.Sprintf("Value %s covers: [", value))
			var tuples []string
			for tuple, covered := range *cover {
				tuples = append(tuples, fmt.Sprintf("%d: %s", tuple, bString(covered)))
			}

			buffer.WriteString(strings.Join(tuples, " "))

			buffer.WriteString("]\n")
		}
	}

	return buffer.String()
}

func bString(b bool) string {
	if b {
		return "y"
	}
	return "n"
}
