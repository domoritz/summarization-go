package main

import (
	"fmt"
	"log"

	"github.com/domoritz/summarization-go/summary"
)

func main() {
	a := summary.NewSingle("x")
	fmt.Println(a.DebugString())

	t, err := summary.NewTupleFromString("a, b, {c d}", 3)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(t.DebugString())

	r, err := summary.NewRelationFromString("x,y,z\na, b, {c d}\n b, c, {d e f}\n a, b, {c e}")
	if err != nil {
		log.Fatal(err)
	}
	r.PrintDebugString()
}
