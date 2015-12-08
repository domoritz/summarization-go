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
	covered   *bool
}

// NewHierachy creates a new hierarchical attribute
func NewHierachy(hierarchy Hierarchy) HierarchyAttribute {
	covered := false
	return HierarchyAttribute{hierarchy, &covered}
}

// Prefix returns true if attribute is prefix of other attribute
func (attr HierarchyAttribute) Prefix(other HierarchyAttribute) bool {
	if len(attr.hierarchy) > len(other.hierarchy) {
		return false
	}

	for i, value := range attr.hierarchy {
		if other.hierarchy[i] != value {
			return false
		}
	}
	return true
}

// DebugString prints the attribute name and value
func (attr *HierarchyAttribute) DebugString() string {
	if attr == nil {
		return "null"
	}

	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(strings.Join(attr.hierarchy, " "))
	buffer.WriteString("]")
	return buffer.String()
}
