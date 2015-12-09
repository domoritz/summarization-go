package summarize

import "testing"

func TestCreate(t *testing.T) {
	cover := TupleCover{0: true, 1: false}
	attribute := Attribute{0, set, "x", nil}
	cell := Cell{0, &attribute, "a", &cover, 0, 0}

	formula := NewFormula(&cell)

	if formula.tupleValue[0] != 0 {
		t.Error("Should have value")
	}
	if formula.tupleValue[1] != 1 {
		t.Error("Should not have value")
	}

	cover2 := TupleCover{1: false, 2: true}
	attribute2 := Attribute{0, set, "y", nil}
	cell2 := Cell{0, &attribute2, "a", &cover2, 0, 0}
	formula.AddCell(&cell2)

	if _, has := formula.tupleValue[0]; has {
		t.Error("Should not have value")
	}
	if formula.tupleValue[1] != 2 {
		t.Error("Should have value")
	}
	if _, has := formula.tupleValue[2]; has {
		t.Error("This should not be affected because we only shrink")
	}
}
