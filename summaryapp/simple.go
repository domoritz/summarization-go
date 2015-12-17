package main

import (
	"fmt"
	"log"

	"github.com/domoritz/summarization-go/summarize"
)

func main() {
	relation, err := summarize.NewIndexFromString("single,single,set,hierarchy\nw,x,y,z\na,b,c d f,a b c\na,b,c,a b\na,b,c,a b c\nb,,d e f,a b\na,b,c e,\na,a,,a")
	//relation, err := summarize.NewIndexFromString("hierarchy\nx\na\na b c e\na b c e\na b c\na b e f")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(relation)

	summary := relation.Summarize(4)
	summary.DebugPrint()
}
