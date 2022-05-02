package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

var file = flag.String("f", "problems.csv", "CSV file to parse for questions")
var timeout = flag.Int("t", 30, "Timeout for the quiz")

func main() {
	f, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	var points int
	for _, record := range records {
		var answer string
		printQuestion(record[0])
		requestAnswer(&answer)
		if compareAnswers(answer, record[1]) {
			points++
		}
	}
	fmt.Printf("RESULT: %d/%d\n", points, len(records))
}

func prepareString(s string) string {
	trimmed := strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsSpace(r)
	})

	return strings.ToLower(trimmed)
}

func compareAnswers(a, b string) bool {
	return prepareString(a) == prepareString(b)
}

func printQuestion(q string) {
	fmt.Printf("Q: %s\n", q)
}

func requestAnswer(a *string) {
	fmt.Print("A: ")
	_, err := fmt.Scanln(a)
	if err != nil {
		log.Fatal(err)
	}
}
