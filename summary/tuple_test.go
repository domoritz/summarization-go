package summary

import "testing"

func TestSatisfies(t *testing.T) {
	a, _ := NewTupleFromString("a,a,a b", []Type{single, set, hierarchy})
	b, _ := NewTupleFromString("a,a b,a b c", []Type{single, set, hierarchy})

	if !a.Satisfies(&a) {
		t.Error("Should satisfy")
	}

	if !b.Satisfies(&a) {
		t.Error("Should satisfy")
	}

	if a.Satisfies(&b) {
		t.Error("Should not satisfy")
	}
}
