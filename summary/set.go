package summary

// SetAttribute has a set of values
import (
	"bytes"
	"fmt"
	"strings"
)

// Set is a set of values. The integer counts how often the value has been covered.
type Set map[string]Counter

// SetAttribute is a set attribute
type SetAttribute struct {
	values Set
}

// NewSet creates a new SetAttribute
func NewSet(value Set) SetAttribute {
	return SetAttribute{value}
}

// Subset returns true if attribute is subset of other attribute
func (a *SetAttribute) Subset(other *SetAttribute) bool {
	for value := range a.values {
		if _, has := other.values[value]; !has {
			return false
		}
	}
	return true
}

func (a *SetAttribute) getDebugValues() []string {
	values := make([]string, 0, len(a.values))
	for value, count := range a.values {
		values = append(values, fmt.Sprintf("%s (%d)", value, *count))
	}
	return values
}

// DebugString prints the attribute name and value
func (a *SetAttribute) DebugString() string {
	if a == nil {
		return "null"
	}

	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString(strings.Join(a.getDebugValues(), " "))
	buffer.WriteString("}")
	return buffer.String()
}
