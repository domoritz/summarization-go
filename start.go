package main

import (
	"fmt"
	"log"

	"github.com/domoritz/summarization-go/summary"
)

func main() {
	a := summary.NewSingle("x", "foo")
	fmt.Println(a.DebugString())

	attributeNames := []string{"x", "y", "z"}
	t, err := summary.NewTupleFromString("a, b, {c d}", attributeNames)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t.DebugString())

	r, err := summary.NewRelationFromString("x,y,z\na, b, {c d}\n b, c, {e f g}")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(r.DebugString())
}
