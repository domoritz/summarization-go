package summary

import (
	"bytes"
	"strings"
)

// Attribute is an attribute
type Attribute interface {
	DebugString() string
}

type baseAttribute struct {
	name string
}

// SingleValueAttribute has a single value
type SingleValueAttribute struct {
	baseAttribute
	value string
}

// DebugString prints the attribute name and value
func (a *SingleValueAttribute) DebugString() string {
	return a.name + "=" + a.value
}

// SetAttribute has a set of values
type SetAttribute struct {
	baseAttribute
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
	buffer.WriteString(a.name)
	buffer.WriteString("={")
	buffer.WriteString(strings.Join(a.getValues(), " "))
	buffer.WriteString("}")
	return buffer.String()
}

// HierarchyAttribute is a hierarchical attribute
type HierarchyAttribute struct {
	baseAttribute
	hierarchy []string
}

// NewSingle creates a new SingleValueAttribute
func NewSingle(name string, value string) *SingleValueAttribute {
	return &SingleValueAttribute{baseAttribute{name}, value}
}

// NewSet creates a new SetAttribute
func NewSet(name string, value map[string]bool) *SetAttribute {
	return &SetAttribute{baseAttribute{name}, value}
}
