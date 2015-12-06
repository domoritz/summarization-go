package summary

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// Tuple is a map of attributes
type Tuple []Attribute

// NewTupleFromString creates a tuple from a string
func NewTupleFromString(description string, numAttr int) (Tuple, error) {
	tuple := make(Tuple, numAttr)

	values := strings.Split(description, ",")
	if len(values) != numAttr {
		err := fmt.Sprintf("Wrong number of attributes. Expected %d but got %d.", numAttr, len(values))
		return nil, errors.New(err)
	}

	for i, value := range values {
		value = strings.TrimSpace(value)
		if strings.HasPrefix(value, "{") {
			value = strings.TrimPrefix(value, "{")
			value = strings.TrimSuffix(value, "}")
			setValues := strings.Split(value, " ")
			setValue := make(map[string]bool)
			for _, v := range setValues {
				setValue[v] = true
			}
			a := NewSet(setValue)
			tuple[i] = a
		} else if strings.HasPrefix(value, "[") {
			// TODO
		} else {
			a := NewSingle(value)
			tuple[i] = a
		}
	}

	return tuple, nil
}

// DebugString prints a tuple without attribute names
func (tuple *Tuple) DebugString() string {
	var buffer bytes.Buffer

	for i := 0; i < len(*tuple); i++ {
		buffer.WriteString((*tuple)[i].DebugString())
		buffer.WriteString(" ")
	}

	return buffer.String()
}
