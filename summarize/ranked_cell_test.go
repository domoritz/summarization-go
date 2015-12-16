package summarize

import (
	"container/heap"
	"testing"
)

func TestHeap(t *testing.T) {
	attr := Attribute{}

	zero := Cell{nil, &attr, "zero"}
	one := Cell{nil, &attr, "one"}
	two := Cell{nil, &attr, "two"}
	three := Cell{nil, &attr, "three"}
	five := Cell{nil, &attr, "five"}

	cells := CellHeap{RankedCell{&zero, 0, -1}, RankedCell{&one, 1, -1}, RankedCell{&three, 3, -1},
		RankedCell{&three, 3, -1}, RankedCell{&five, 5, -1}, RankedCell{&two, 2, -1}}

	heap.Init(&cells)

	previous := heap.Pop(&cells).(RankedCell)

	for len(cells) > 0 {
		currPeek := *cells.Peek()
		current := heap.Pop(&cells).(RankedCell)

		if currPeek.potential != current.potential {
			t.Error("Peek was not pop")
		}
		if previous.potential < current.potential {
			t.Error("Previous should have higher priority")
		}
	}
}
