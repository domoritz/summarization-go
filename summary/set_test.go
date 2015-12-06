package summary

import "testing"

func TestSubset(t *testing.T) {
	a := NewSet(Set{"a": true})
	b := NewSet(Set{"a": true, "b": true})

	if !a.Subset(&a) {
		t.Error("Should satisfy")
	}

	if !a.Subset(&b) {
		t.Error("Should satisfy")
	}

	if b.Subset(&a) {
		t.Error("Should not satisfy")
	}
}
