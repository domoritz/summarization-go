package summary

// SetAttribute has a set of values
import (
	"bytes"
	"fmt"
	"strings"
)

// SetAttribute is a set attribute
type SetAttribute map[string]*bool

// NewSet creates a new SetAttribute
func NewSet(elements []string) SetAttribute {
	set := make(SetAttribute)
	for _, element := range elements {
		covered := false
		set[element] = &covered
	}
	return set
}

// Subset returns true if attribute is subset of other attribute
func (attr SetAttribute) Subset(other SetAttribute) bool {
	for value := range attr {
		if _, has := other[value]; !has {
			return false
		}
	}
	return true
}

func (attr SetAttribute) getDebugValues() []string {
	values := make([]string, 0, len(attr))
	for value, covered := range attr {
		values = append(values, fmt.Sprintf("%s (%t)", value, *covered))
	}
	return values
}

// DebugString prints the attribute name and value
func (attr *SetAttribute) DebugString() string {
	if attr == nil {
		return "null"
	}

	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString(strings.Join(attr.getDebugValues(), " "))
	buffer.WriteString("}")
	return buffer.String()
}
