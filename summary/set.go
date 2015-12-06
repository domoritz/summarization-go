package summary

// SetAttribute has a set of values
import (
	"bytes"
	"strings"
)

// Set is a set
type Set map[string]bool

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

func (a *SetAttribute) getValues() []string {
	keys := make([]string, 0, len(a.values))
	for k := range a.values {
		keys = append(keys, k)
	}
	return keys
}

// DebugString prints the attribute name and value
func (a *SetAttribute) DebugString() string {
	if a == nil {
		return "null"
	}

	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString(strings.Join(a.getValues(), " "))
	buffer.WriteString("}")
	return buffer.String()
}
