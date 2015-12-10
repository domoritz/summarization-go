package summarize

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Attribute is an attribute
type Attribute struct {
	index         int            // id for this attribute, used to see what attributes have been used in formula
	attributeType Type           // attribute type
	attributeName string         // attribute name
	valueIndex    map[string]int // index for attribute values
	cells         []Cell         // values and what tuples are covered
}

// RelationIndex is an inverted index
type RelationIndex struct {
	attrs     []Attribute
	numTuples int
	numValues int
}

func addCell(attr *Attribute, value string, tuple int) bool {
	added := false

	idx, has := attr.valueIndex[value]
	if !has {
		c := MakeCell(attr, value)
		c.cover[tuple] = false
		attr.valueIndex[value] = len(attr.cells)
		attr.cells = append(attr.cells, c)
		added = true
	} else {
		attr.cells[idx].cover[tuple] = false
	}
	return added
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

		index[i].attributeName = names[i]
		index[i].index = i
		index[i].valueIndex = make(map[string]int)
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
			value = strings.TrimSpace(value)

			if len(value) == 0 {
				// null
				continue
			}

			switch index[i].attributeType {
			case single:
				if addCell(&index[i], value, tuple) {
					numValues++
				}
			case set:
				setValues := strings.Split(value, " ")
				for _, setValue := range setValues {
					if addCell(&index[i], setValue, tuple) {
						numValues++
					}
				}
			case hierarchy:
				prefix := ""
				hValues := strings.Split(value, " ")
				for _, hValue := range hValues {
					prefix += hValue
					if addCell(&index[i], prefix, tuple) {
						numValues++
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
		buffer.WriteString(fmt.Sprintf("Attribute %s (%s):\n", attribute.attributeName, attribute.attributeType))
		for _, cell := range attribute.cells {
			buffer.WriteString(fmt.Sprintf("Value %s covers: [", cell.value))
			var tuples []string
			for tuple, covered := range cell.cover {
				tuples = append(tuples, fmt.Sprintf("%d:%s", tuple, bString(covered)))
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
