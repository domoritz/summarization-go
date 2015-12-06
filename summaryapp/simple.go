package main

import (
	"log"

	"github.com/domoritz/summarization-go/summary"
)

func main() {
	r, err := summary.NewRelationFromString("single,single,set\nx,y,z\na, b, c d\n b, c, d e f\n a, b, c e")
	if err != nil {
		log.Fatal(err)
	}
	r.PrintDebugString()
}
