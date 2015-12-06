package summary

import "testing"

func TestSubset(t *testing.T) {
	a := NewSet(Set{"a": false})
	b := NewSet(Set{"a": false, "b": false})

	if ok, size := a.SubsetCover(&a); !ok || size != 1 {
		t.Error("Should satisfy")
	}

	if ok, size := a.SubsetCover(&b); !ok || size != 1 {
		t.Error("Should satisfy")
	}

	if ok, _ := b.SubsetCover(&a); ok {
		t.Error("Should not satisfy")
	}

	c := NewSet(Set{"a": false, "b": true})
	d := NewSet(Set{"a": true, "b": false})
	e := NewSet(Set{"a": true, "b": true})

	if ok, size := a.SubsetCover(&c); !ok || size != 1 {
		t.Error("Should satisfy")
	}

	if ok, size := a.SubsetCover(&d); !ok || size != 0 {
		t.Error("Should satisfy")
	}

	if ok, size := a.SubsetCover(&e); !ok || size != 0 {
		t.Error("Should satisfy")
	}
}
