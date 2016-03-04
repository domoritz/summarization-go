package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/domoritz/summarization-go/summarize"
)

var cpuprofile = flag.String("cpuprofile", "cpu.prof", "write cpu profile to file")
var memprofile = flag.String("memprofile", "mem.prof", "write memory profile to this file")
var weightfunc = flag.String("weightfunc", "equal", "weight function (equal or exponential)")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	rand.Seed(42)

	types := []string{"single", "single", "single", "set", "set"}
	names := []string{"s0", "s1", "s2", "set0", "set1"}
	numTuples := 100000

	assessor := summarize.MakeEqualWeightAssessor()

	if *weightfunc == "exponential" {
		assessor = summarize.MakeExponentialAssessor([]float64{0.5, 0.5, 0.5, 0.5, 0.5})
		assessor.NumTuples = numTuples
	}

	relation, err := summarize.NewIndex(types, names, numTuples)
	if err != nil {
		log.Fatal(err)
	}

	attrs := relation.Attrs()

	for i := 0; i < numTuples; i++ {
		(*attrs)[0].AddCell(randomdata.FirstName(randomdata.Female), i, assessor)
		(*attrs)[1].AddCell(randomdata.LastName(), i, assessor)
		(*attrs)[2].AddCell(randomdata.FullName(randomdata.RandomGender), i, assessor)

		for j := 0; j < 3; j++ {
			(*attrs)[3].AddCell(randomdata.City(), i, assessor)
		}
		for j := 0; j < 6; j++ {
			(*attrs)[4].AddCell(randomdata.State(randomdata.Large), i, assessor)
		}
	}

	start := time.Now()
	summary := relation.Summarize(200)
	elapsed := time.Since(start)
	log.Printf("Summarization took %s\n", elapsed)

	summary.DebugPrint()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}
