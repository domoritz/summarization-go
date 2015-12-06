package summary

import "testing"

func TestCoverage(t *testing.T) {
	a, _ := NewTupleFromString("a,a,a b", []Type{single, set, hierarchy})
	b, _ := NewTupleFromString("a,a b,a b c", []Type{single, set, hierarchy})

	if ok, size := a.Coverage(&a); !ok || size != 4 {
		t.Error("Should satisfy")
	}

	if ok, size := b.Coverage(&a); !ok || size != 4 {
		t.Error("Should satisfy")
	}

	if ok, _ := a.Coverage(&b); ok {
		t.Error("Should not satisfy")
	}
}

func TestCoverageNulls(t *testing.T) {
	validSingle, _ := NewTupleFromString("a", []Type{single})
	nullSingle, _ := NewTupleFromString("", []Type{single})

	validSet, _ := NewTupleFromString("a", []Type{set})
	nullSet, _ := NewTupleFromString("", []Type{set})

	validHierarchy, _ := NewTupleFromString("a", []Type{hierarchy})
	nullHierarchy, _ := NewTupleFromString("", []Type{hierarchy})

	// both null should be satisfied
	if ok, size := nullSingle.Coverage(&nullSingle); !ok || size != 0 {
		t.Error("Should satisfy")
	}
	if ok, size := nullSet.Coverage(&nullSet); !ok || size != 0 {
		t.Error("Should satisfy")
	}
	if ok, size := nullHierarchy.Coverage(&nullHierarchy); !ok || size != 0 {
		t.Error("Should satisfy")
	}

	// tuple should always satisfy empty formula
	if ok, size := validSingle.Coverage(&nullSingle); !ok || size != 0 {
		t.Error("Should satisfy")
	}
	if ok, size := validSet.Coverage(&nullSet); !ok || size != 0 {
		t.Error("Should satisfy")
	}
	if ok, size := validHierarchy.Coverage(&nullHierarchy); !ok || size != 0 {
		t.Error("Should satisfy")
	}

	// null tuple should never satisfy formula with value
	if ok, _ := nullSingle.Coverage(&validSingle); ok {
		t.Error("Should not satisfy")
	}
	if ok, _ := nullSet.Coverage(&validSet); ok {
		t.Error("Should not satisfy")
	}
	if ok, _ := nullHierarchy.Coverage(&validHierarchy); ok {
		t.Error("Should not satisfy")
	}
}
