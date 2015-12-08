package summary

import "fmt"

// SingleValueAttribute has a single value
type SingleValueAttribute struct {
	value   string
	covered *bool
}

// NewSingle creates a new SingleValueAttribute
func NewSingle(value string) SingleValueAttribute {
	covered := false
	return SingleValueAttribute{value, &covered}
}

// Equal returns true if attribute satisfies the other attribute
func (attr SingleValueAttribute) Equal(other SingleValueAttribute) bool {
	return attr.value == other.value
}

// DebugString prints the attribute name and value
func (attr *SingleValueAttribute) DebugString() string {
	if attr == nil {
		return "null"
	}
	return fmt.Sprintf("%s (%t)", attr.value, *attr.covered)
}
