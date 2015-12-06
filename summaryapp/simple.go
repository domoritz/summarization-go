package main

import (
	"log"

	"github.com/domoritz/summarization-go/summary"
)

func main() {
	r, err := summary.NewRelationFromString("single,single,set\nx,y,z\na,b,c d\nb,,d e f\na,b,c e\na,a,")
	if err != nil {
		log.Fatal(err)
	}
	r.PrintDebugString()
}