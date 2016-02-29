package summarize

import "testing"

func TestCreate(t *testing.T) {
	y := Cover{true, 1}
	n := Cover{false, 1}
	attribute := Attribute{0, set, "x", nil, nil}
	cell := Cell{TupleCover{0: &y, 1: &n}, &attribute, "a", true}

	formula := NewFormula(cell)

	if formula.tupleCover[0] != 0 {
		t.Error("Should have cover")
	}
	if formula.tupleCover[1] != 1 {
		t.Error("Should not have cover")
	}

	attribute2 := Attribute{0, set, "x", nil, nil}
	cell2 := Cell{TupleCover{1: &n, 2: &y}, &attribute2, "a", true}
	formula.AddCell(cell2)

	if _, has := formula.tupleCover[0]; has {
		t.Error("Should not have cover")
	}
	if formula.tupleCover[1] != 2 {
		t.Error("Should have cover")
	}
	if _, has := formula.tupleCover[2]; has {
		t.Error("This should not be affected because we only shrink")
	}
}
