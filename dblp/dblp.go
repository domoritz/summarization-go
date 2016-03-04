package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/domoritz/summarization-go/summarize"
	_ "github.com/mattn/go-sqlite3"
)

var database = flag.String("db", "./dblp.sqlite", "the sqlite database")

func main() {
	flag.Parse()
	db, err := sql.Open("sqlite3", *database)
	checkErr(err)

	for true {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter Query (e.g. 'database'): ")

		query, err := reader.ReadString('\n')
		checkErr(err)

		stmt, err := db.Prepare("select count(*) from data where data match ?")
		checkErr(err)
		defer stmt.Close()
		var numTuples int
		err = stmt.QueryRow(query).Scan(&numTuples)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("# of results:", numTuples)

		types := []string{"set", "single", "single", "single", "single", "single", "single"}
		names := []string{"author", "school", "journal", "publisher", "year", "organization", "institution"}
		weights := []float64{1, 0.7, 0.6, 0.3, 0.1, 0.7, 0.7}

		assessor := summarize.MakeExponentialAssessor(weights)
		assessor.NumTuples = numTuples

		relation, err := summarize.NewIndex(types, names, numTuples)
		if err != nil {
			log.Fatal(err)
		}

		attrs := relation.Attrs()

		rows, err := db.Query("select author, school, journal, publisher, year, organization, institution from data where data match ?", query)
		checkErr(err)
		defer rows.Close()

		i := 0
		for rows.Next() {
			var authors string
			var school string
			var journal string
			var publisher string
			var year string
			var organization string
			var institution string
			err = rows.Scan(&authors, &school, &journal, &publisher, &year, &organization, &institution)
			checkErr(err)

			for _, author := range strings.Split(authors, ",") {
				author = strings.TrimSpace(author)
				if len(author) > 0 {
					(*attrs)[0].AddCell(author, i, assessor)
				}
			}

			if len(school) > 0 {
				(*attrs)[1].AddCell(school, i, assessor)
			}
			if len(journal) > 0 {
				(*attrs)[2].AddCell(journal, i, assessor)
			}

			if len(publisher) > 0 {
				(*attrs)[3].AddCell(publisher, i, assessor)
			}

			if len(year) > 0 {
				(*attrs)[4].AddCell(year, i, assessor)
			}

			if len(organization) > 0 {
				(*attrs)[5].AddCell(organization, i, assessor)
			}

			if len(institution) > 0 {
				(*attrs)[6].AddCell(institution, i, assessor)
			}

			i++
		}

		start := time.Now()
		summary := relation.Summarize(16)
		elapsed := time.Since(start)
		log.Printf("Summarization took %s\n", elapsed)

		summary.DebugPrint()
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
