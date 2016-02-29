package summarize

import (
	"container/heap"
	"testing"

	"golang.org/x/tools/container/intsets"
)

var y = Cover{true, 1}
var n = Cover{false, 1}

func TestHeap(t *testing.T) {
	attr := Attribute{}

	zero := Cell{nil, &attr, "zero", true}
	one := Cell{nil, &attr, "one", true}
	two := Cell{nil, &attr, "two", true}
	three := Cell{nil, &attr, "three", true}
	five := Cell{nil, &attr, "five", true}

	cells := CellHeap{&RankedCell{&zero, 0, -1, 0}, &RankedCell{&one, 1, -1, 1}, &RankedCell{&three, 3, -1, 2},
		&RankedCell{&three, 3, -1, 3}, &RankedCell{&five, 5, -1, 4}, &RankedCell{&two, 2, -1, 5}}

	heap.Init(&cells)

	if !cells.Valid(0) {
		t.Error("Invalid heap")
	}

	previous := heap.Pop(&cells).(*RankedCell)

	if !cells.Valid(0) {
		t.Error("Invalid heap")
	}

	for len(cells) > 0 {
		currPeek := *cells.Peek()
		current := heap.Pop(&cells).(*RankedCell)

		if !cells.Valid(0) {
			t.Error("Invalid heap")
		}

		if currPeek.potential != current.potential {
			t.Error("Peek was not pop")
		}
		if previous.potential < current.potential {
			t.Error("Previous should have higher priority")
		}
	}
}

func TestRecomputeCoverage(t *testing.T) {
	cover := make(TupleCover)
	cover[12] = &n
	cover[17] = &y
	cover[42] = &n
	cell := Cell{cover, nil, "x", true}
	rankedCell := RankedCell{&cell, 10, -1, 0}

	result := rankedCell.recomputeCoverage()

	if result != 2 {
		t.Error("Wrong cover")
	}

	if result != rankedCell.potential {
		t.Error("Should be equal")
	}
}

func TestRecomputeCoverageFormula(t *testing.T) {
	cover := make(TupleCover)
	cover[12] = &n
	cover[17] = &y
	cover[42] = &n
	cover[99] = &n
	cover[123] = &n
	cell := Cell{cover, nil, "x", true}
	rankedCell := RankedCell{&cell, 10, -1, 0}

	covers := make(TupleCovers)
	covers[17] = 2
	covers[42] = 1
	covers[123] = 3
	covers[255] = 2
	var set intsets.Sparse
	formula := Formula{nil, set, covers, 5}

	formulaPotential := rankedCell.recomputeFormulaCoverage(&formula)

	// (2 + 1 + 3) + (2) - 5 = 3
	if formulaPotential != 3 {
		t.Error("Wrong cover", formulaPotential)
	}

	if formulaPotential != rankedCell.potential {
		t.Error("Should be equal")
	}

	// 2
	if rankedCell.maxPotential != 2 {
		t.Error("Wrong cover")
	}
}
