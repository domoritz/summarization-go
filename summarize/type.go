package summarize

// Type is the attribute type
type Type int

const (
	single Type = iota
	set
	hierarchy
)

func (t Type) String() string {
	switch t {
	case single:
		return "single"
	case set:
		return "set"
	case hierarchy:
		return "hierarchy"
	default:
		return "unknown"
	}
}
