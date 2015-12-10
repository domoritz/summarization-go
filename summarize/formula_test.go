package summarize

import "testing"

func TestCreate(t *testing.T) {
	attribute := Attribute{0, set, "x", nil, nil}
	cell := Cell{TupleCover{0: true, 1: false}, &attribute, "a"}

	formula := NewFormula(cell)

	if formula.tupleValue[0] != 0 {
		t.Error("Should have value")
	}
	if formula.tupleValue[1] != 1 {
		t.Error("Should not have value")
	}

	attribute2 := Attribute{0, set, "x", nil, nil}
	cell2 := Cell{TupleCover{1: false, 2: true}, &attribute2, "a"}
	formula.AddCell(cell2)

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
