package main

import (
	"fmt"
	"log"

	"github.com/domoritz/summarization-go/summarize"
)

func main() {
	relation, err := summarize.NewIndexFromString("single,single,set\nx,y,z\na,b,c d\na,b,c\na,b,c\nb,,d 4 5\na,b,c e\na,a,")
	if err != nil {
		log.Fatal(err)
	}

	summary := relation.Summarize(3)

	fmt.Println("Summary:")
	fmt.Println(summary)

	fmt.Println("🐙")
}
