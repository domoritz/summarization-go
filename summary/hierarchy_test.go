package summary

import "testing"

func TestPrefix(t *testing.T) {
	a := NewHierachy(Hierarchy{"a", "b"})
	b := NewHierachy(Hierarchy{"a", "b", "c"})

	if !a.Prefix(a) {
		t.Error("Should satisfy")
	}

	if !a.Prefix(b) {
		t.Error("Should satisfy")
	}

	if b.Prefix(a) {
		t.Error("Should not satisfy")
	}
}
