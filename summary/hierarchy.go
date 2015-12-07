package summary

// HierarchyAttribute is a hierarchical attribute
import (
	"bytes"
	"strings"
)

// Hierarchy is a hierarchy
type Hierarchy []string

// HierarchyAttribute is hierarchical
type HierarchyAttribute struct {
	hierarchy Hierarchy
	// TODO: add counts
}

// NewHierachy creates a new hierarchical attribute
func NewHierachy(hierarchy Hierarchy) HierarchyAttribute {
	return HierarchyAttribute{hierarchy}
}

// Prefix returns true if attribute is prefix of other attribute
func (a *HierarchyAttribute) Prefix(other *HierarchyAttribute) bool {
	if len(a.hierarchy) > len(other.hierarchy) {
		return false
	}

	for i, value := range a.hierarchy {
		if other.hierarchy[i] != value {
			return false
		}
	}
	return true
}

// DebugString prints the attribute name and value
func (a *HierarchyAttribute) DebugString() string {
	if a == nil {
		return "null"
	}

	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(strings.Join(a.hierarchy, " "))
	buffer.WriteString("]")
	return buffer.String()
}
