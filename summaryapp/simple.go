package main

import (
	"fmt"
	"log"

	"github.com/domoritz/summarization-go/summarize"
)

func main() {
	relation, err := summarize.NewIndexFromString("single,single,set\nx,y,z\n0,1,2 3\n0,1,2\n0,1,2\n1,,3 4 5\n0,1,2 4\n0,0,")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Relation Index:")
	fmt.Println(relation)

	summary := relation.Summarize(3)

	fmt.Println("Summary:")
	fmt.Println(summary)

	fmt.Println("🐙")
}
