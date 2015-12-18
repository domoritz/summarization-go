package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/domoritz/summarization-go/summarize"
)

func main() {
	fmt.Println("numTuples, size, actual size, trial, elapsed nanoseconds")

	// test various sizes
	sizes := []int{1, 10, 50, 100, 200, 300}

	// test different # of tuples
	for numTuples := 10000; numTuples < 1000000; numTuples += 20000 {

		types := []string{"single", "single", "single", "set", "set"}
		names := []string{"s0", "s1", "s2", "set0", "set1"}

		relation, err := summarize.NewIndex(types, names, numTuples)
		if err != nil {
			log.Fatal(err)
		}

		attrs := relation.Attrs()

		for i := 0; i < numTuples; i++ {
			(*attrs)[0].AddCell(randomdata.FirstName(randomdata.Female), i)
			(*attrs)[1].AddCell(randomdata.LastName(), i)
			(*attrs)[2].AddCell(randomdata.FullName(randomdata.RandomGender), i)

			for j := 0; j < 3; j++ {
				(*attrs)[3].AddCell(randomdata.City(), i)
			}
			for j := 0; j < 6; j++ {
				(*attrs)[4].AddCell(randomdata.State(randomdata.Large), i)
			}
		}

		for _, size := range sizes {

			// run each experiment multiple times
			for trial := 0; trial < 6; trial++ {
				start := time.Now()
				summary := relation.Summarize(size)
				elapsed := time.Since(start)
				relation.Reset()
				fmt.Printf("%d, %d, %d, %d, %v\n", numTuples, size, len(summary.Summary), trial, elapsed.Nanoseconds())
			}
		}
	}
}
