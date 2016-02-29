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
	attrs     []Attribute // the attributes
	numTuples int         // not really needed
}

// Attrs returns the attributes
func (relation RelationIndex) Attrs() *[]Attribute {
	return &relation.attrs
}

// AddCell adds a cell to an attribute
func (attr *Attribute) AddCell(value string, tuple int, assessor Assessor) bool {
	added := false

	cover := Cover{false, assessor.Weight(attr, tuple)}

	idx, has := attr.valueIndex[value]
	if !has {
		c := MakeCell(attr, value, assessor.function == Equal)
		c.covers[tuple] = &cover
		attr.valueIndex[value] = len(attr.cells)
		attr.cells = append(attr.cells, c)
		added = true
	} else {
		attr.cells[idx].covers[tuple] = &cover
	}
	return added
}

// NewIndex creates a new index
func NewIndex(typeNames []string, names []string, numTuples int) (*RelationIndex, error) {
	if len(names) != len(typeNames) {
		err := fmt.Sprintf("Mismatching number of names and types. %d != %d", len(names), len(typeNames))
		return nil, errors.New(err)
	}

	index := make([]Attribute, len(typeNames))

	for i, attributeType := range typeNames {
		attr := &index[i]

		switch attributeType {
		case single.String():
			attr.attributeType = single
		case set.String():
			attr.attributeType = set
		case hierarchy.String():
			attr.attributeType = hierarchy
		}

		attr.attributeName = names[i]
		attr.index = i
		attr.valueIndex = make(map[string]int)
	}

	return &RelationIndex{index, numTuples}, nil
}

// NewIndexFromString creates a relation index from a string
func NewIndexFromString(description string, assessor Assessor) (*RelationIndex, error) {
	lines := strings.Split(description, "\n")

	typeNames := strings.Split(lines[0], ",")
	names := strings.Split(lines[1], ",")

	relation, err := NewIndex(typeNames, names, len(lines[2:]))
	if err != nil {
		return nil, err
	}

	assessor.NumTuples = relation.numTuples

	index := relation.attrs
	numAttr := len(index)

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
				index[i].AddCell(value, tuple, assessor)
			case set:
				setValues := strings.Split(value, " ")
				for _, setValue := range setValues {
					index[i].AddCell(setValue, tuple, assessor)
				}
			case hierarchy:
				prefix := ""
				hValues := strings.Split(value, " ")
				for _, hValue := range hValues {
					p := ""
					if len(prefix) > 0 {
						p = "/"
					}
					prefix += p + hValue
					index[i].AddCell(prefix, tuple, assessor)
				}
			}
		}
	}

	return relation, nil
}

// Reset resets coverage
func (relation *RelationIndex) Reset() {
	for ia := range relation.attrs {
		attr := &relation.attrs[ia]
		for ic := range attr.cells {
			cell := &attr.cells[ic]
			for tc := range cell.covers {
				cell.covers[tc].covered = false
			}
		}
	}
}

func (relation RelationIndex) String() string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "Relation Index (%d attributes, %d tuples):\n", len(relation.attrs), relation.numTuples)
	for _, attribute := range relation.attrs {
		fmt.Fprintf(&buffer, "Attribute %s (%s) of length %d:\n", attribute.attributeName, attribute.attributeType, len(attribute.cells))
		for _, cell := range attribute.cells {
			fmt.Fprintf(&buffer, " Value '%s' covers: [", cell.value)
			var tuples []string
			for tuple, covered := range cell.covers {
				tuples = append(tuples, fmt.Sprintf("%d:(%s %.3g)", tuple, bString(covered.covered), covered.weight))
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
