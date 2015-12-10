package summarize

import (
	"container/heap"
	"testing"
)

func TestHeap(t *testing.T) {
	cells := CellHeap{RankedCell{nil, 1}, RankedCell{nil, 3},
		RankedCell{nil, 3}, RankedCell{nil, 6}, RankedCell{nil, 2}}

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
