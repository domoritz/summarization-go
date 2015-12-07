package summary

import "fmt"

// Counter counts the coverage
type Counter int

// SingleValueAttribute has a single value
type SingleValueAttribute struct {
	value   string
	covered Counter
}

// NewSingle creates a new SingleValueAttribute
func NewSingle(value string) SingleValueAttribute {
	return SingleValueAttribute{value, 0}
}

// Equal returns true if attribute satisfies the other attribute
func (attr *SingleValueAttribute) Equal(other *SingleValueAttribute) bool {
	return attr.value == other.value
}

// DebugString prints the attribute name and value
func (attr *SingleValueAttribute) DebugString() string {
	if attr == nil {
		return "null"
	}
	return fmt.Sprintf("%s (%d)", attr.value, attr.covered)
}
