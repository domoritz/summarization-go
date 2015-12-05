package summary

import (
	"bytes"
	"errors"
	"sort"
	"strings"
)

// Tuple is a map of attributes
type Tuple map[string]Attribute

// NewTupleFromString creates a tuple from a string
func NewTupleFromString(description string, attributes []string) (Tuple, error) {
	tuple := make(Tuple)

	for i, value := range strings.Split(description, ",") {
		if i >= len(attributes) {
			return nil, errors.New("Not enough attributes")
		}
		name := attributes[i]
		value = strings.TrimSpace(value)
		if strings.HasPrefix(value, "{") {
			value = strings.TrimPrefix(value, "{")
			value = strings.TrimSuffix(value, "}")
			values := strings.Split(value, " ")
			setValue := make(map[string]bool)
			for _, v := range values {
				setValue[v] = true
			}
			a := NewSet(name, setValue)
			tuple[name] = a
		} else if strings.HasPrefix(value, "[") {
			// TODO
		} else {
			a := NewSingle(name, value)
			tuple[name] = a
		}
	}

	return tuple, nil
}

// DebugString prints a tuple
func (tuple *Tuple) DebugString() string {
	var buffer bytes.Buffer

	names := make([]string, 0, len(*tuple))
	for n := range *tuple {
		names = append(names, n)
	}
	sort.Strings(names)

	for _, name := range names {
		buffer.WriteString((*tuple)[name].DebugString())
		buffer.WriteString(" ")
	}

	return buffer.String()
}
