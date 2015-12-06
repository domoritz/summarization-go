package summary

import (
	"bytes"
	"strings"
)

// Attribute is an attribute
type Attribute interface {
	DebugString() string
}

// SingleValueAttribute has a single value
type SingleValueAttribute struct {
	value string
}

// DebugString prints the attribute name and value
func (a *SingleValueAttribute) DebugString() string {
	return a.value
}

// SetAttribute has a set of values
type SetAttribute struct {
	values map[string]bool
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

// HierarchyAttribute is a hierarchical attribute
type HierarchyAttribute struct {
	hierarchy []string
}

// NewSingle creates a new SingleValueAttribute
func NewSingle(value string) *SingleValueAttribute {
	return &SingleValueAttribute{value}
}

// NewSet creates a new SetAttribute
func NewSet(value map[string]bool) *SetAttribute {
	return &SetAttribute{value}
}
