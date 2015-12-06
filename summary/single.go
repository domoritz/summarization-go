package summary

// SingleValueAttribute has a single value
type SingleValueAttribute struct {
	value string
}

// NewSingle creates a new SingleValueAttribute
func NewSingle(value string) SingleValueAttribute {
	return SingleValueAttribute{value}
}

// Equal returns true if attribute satisfies the other attribute
func (a *SingleValueAttribute) Equal(other *SingleValueAttribute) bool {
	return a.value == other.value
}

// DebugString prints the attribute name and value
func (a *SingleValueAttribute) DebugString() string {
	if a == nil {
		return "null"
	}
	return a.value
}
