package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"
)

var file = flag.String("f", "problems.csv", "CSV file to parse for questions")
var timeout = flag.Int("t", 30, "Timeout for the quiz")

func main() {
	flag.Parse()
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

	fmt.Print("Press ENTER to start the quiz")
	_, err = fmt.Scanln()
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.After(time.Duration(*timeout) * time.Second)
	onePointListener := make(chan struct{})
	quizDoneListener := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	var points int
	go doQuiz(ctx, onePointListener, quizDoneListener, records)

loop:
	for {
		select {
		case <-timeout:
			fmt.Println("\nTIMEOUT!")
			cancel()
			break loop
		case <-onePointListener:
			points++
		case <-quizDoneListener:
			break loop
		}
	}

	fmt.Printf("RESULT: %d/%d\n", points, len(records))
}

func doQuiz(ctx context.Context, onePointListener chan<- struct{}, quizDoneListener chan<- struct{}, records [][]string) {
	for _, record := range records {
		select {
		case <-ctx.Done():
			return
		default:
			var answer string
			printQuestion(record[0])
			requestAnswer(&answer)
			if compareAnswers(answer, record[1]) {
				onePointListener <- struct{}{}
			}
		}
	}
	quizDoneListener <- struct{}{}
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
