package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/Pallinder/go-randomdata"
	"github.com/domoritz/summarization-go/summarize"
)

var cpuprofile = flag.String("cpuprofile", "cpu.prof", "write cpu profile to file")
var memprofile = flag.String("memprofile", "mem.prof", "write memory profile to this file")

func main() {
	types := []string{"single", "single", "single", "set", "set"}
	names := []string{"s0", "s1", "s2", "set0", "set1"}
	numTuples := 1000

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
		for j := 0; j < 5; j++ {
			(*attrs)[4].AddCell(randomdata.State(randomdata.Large), i)
		}
	}

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	summary := relation.Summarize(100)

	fmt.Println("Summary:")
	summary.DebugPrint()
}
