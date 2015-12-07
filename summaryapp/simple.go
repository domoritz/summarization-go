package main

import (
	"fmt"
	"log"

	"github.com/domoritz/summarization-go/summary"
)

func main() {
	relation, err := summary.NewRelationFromString("single,single,set\nx,y,z\na,b,c d\nb,,d e f\na,b,c e\na,a,")
	if err != nil {
		log.Fatal(err)
	}

	summary := relation.Summarize(1)

	fmt.Println("Summary:")
	summary.PrintDebugString()
}
