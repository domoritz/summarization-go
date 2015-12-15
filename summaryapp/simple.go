package main

import (
	"fmt"
	"log"

	"github.com/domoritz/summarization-go/summarize"
)

func main() {
	relation, err := summarize.NewIndexFromString("single,single,set\nx,y,z\na,b,c d f\na,b,c\na,b,c\nb,,d e f\na,b,c e\na,a,")
	if err != nil {
		log.Fatal(err)
	}

	summary := relation.Summarize(3)

	fmt.Println("Summary:")
	summary.DebugPrint()
}
