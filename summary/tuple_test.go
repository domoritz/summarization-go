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

func TestSatisfiesNulls(t *testing.T) {
	validSingle, _ := NewTupleFromString("a", []Type{single})
	nullSingle, _ := NewTupleFromString("", []Type{single})

	validSet, _ := NewTupleFromString("a", []Type{set})
	nullSet, _ := NewTupleFromString("", []Type{set})

	validHierarchy, _ := NewTupleFromString("a", []Type{hierarchy})
	nullHierarchy, _ := NewTupleFromString("", []Type{hierarchy})

	// both null should be satisfied
	if !nullSingle.Satisfies(&nullSingle) {
		t.Error("Should satisfy")
	}
	if !nullSet.Satisfies(&nullSet) {
		t.Error("Should satisfy")
	}
	if !nullHierarchy.Satisfies(&nullHierarchy) {
		t.Error("Should satisfy")
	}

	// tuple should always satisfy empty formula
	if !validSingle.Satisfies(&nullSingle) {
		t.Error("Should satisfy")
	}
	if !validSet.Satisfies(&nullSet) {
		t.Error("Should satisfy")
	}
	if !validHierarchy.Satisfies(&nullHierarchy) {
		t.Error("Should satisfy")
	}

	// null tuple should never satisfy formula with value
	if nullSingle.Satisfies(&validSingle) {
		t.Error("Should not satisfy")
	}
	if nullSet.Satisfies(&validSet) {
		t.Error("Should not satisfy")
	}
	if nullHierarchy.Satisfies(&validHierarchy) {
		t.Error("Should not satisfy")
	}
}

func TestSize(t *testing.T) {
	a, _ := NewTupleFromString("a,a b c,a b", []Type{single, set, hierarchy})

	if a.Size() != 6 {
		t.Error("Size should be 6")
	}
}
