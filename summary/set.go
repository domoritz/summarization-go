package summary

// SetAttribute has a set of values
import (
	"bytes"
	"strings"
)

// Set is a set of values. The boolean value indicates whether the particular value is covered.
type Set map[string]bool

// SetAttribute is a set attribute
type SetAttribute struct {
	values Set
}

// NewSet creates a new SetAttribute
func NewSet(value Set) SetAttribute {
	return SetAttribute{value}
}

// SubsetCover returns true if attribute is subset of other attribute. If a is a subset of other, also returns the count of new covers.
func (a *SetAttribute) SubsetCover(other *SetAttribute) (bool, int) {
	cover := 0
	for value := range a.values {
		if covered, has := other.values[value]; has {
			if !covered {
				cover++
			}
		} else {
			return false, -1
		}

	}
	return true, cover
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
