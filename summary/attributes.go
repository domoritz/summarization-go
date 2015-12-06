package summary

import (
	"bytes"
	"strings"
)

//
// Single
//

// SingleValueAttribute has a single value
type SingleValueAttribute struct {
	value string
}

// NewSingle creates a new SingleValueAttribute
func NewSingle(value string) SingleValueAttribute {
	return SingleValueAttribute{value}
}

// Satisfies returns true if attribute satisfies the other attribute
func (a *SingleValueAttribute) Satisfies(other *SingleValueAttribute) bool {
	return a.value == other.value
}

// DebugString prints the attribute name and value
func (a *SingleValueAttribute) DebugString() string {
	return a.value
}

//
// Set
//

// SetAttribute has a set of values
type SetAttribute struct {
	values map[string]bool
}

// NewSet creates a new SetAttribute
func NewSet(value map[string]bool) SetAttribute {
	return SetAttribute{value}
}

// Satisfies returns true if attribute satisfies the other attribute
func (a *SetAttribute) Satisfies(other *SetAttribute) bool {
	for value := range other.values {
		if _, has := a.values[value]; !has {
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
	var buffer bytes.Buffer
	buffer.WriteString("{")
	buffer.WriteString(strings.Join(a.getValues(), " "))
	buffer.WriteString("}")
	return buffer.String()
}

//
// Hierarchy
//

// HierarchyAttribute is a hierarchical attribute
type HierarchyAttribute struct {
	hierarchy []string
}

// DebugString prints the attribute name and value
func (a *HierarchyAttribute) DebugString() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(strings.Join(a.hierarchy, " "))
	buffer.WriteString("]")
	return buffer.String()
}
