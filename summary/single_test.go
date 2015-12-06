package summary

import "testing"

func TestEqual(t *testing.T) {
	a := NewSingle("a")
	b := NewSingle("b")

	if !a.Equal(&a) {
		t.Error("Should satisfy")
	}

	if a.Equal(&b) {
		t.Error("Should not satisfy")
	}
}
