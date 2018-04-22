package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

type Problem struct {
	question string
	answer   string
}

var csvNewReader = csv.NewReader

func parseProblems(file io.Reader) ([]Problem, error) {
	r := csvNewReader(file)
	problems := make([]Problem, 0)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		problems = append(problems, Problem{question: row[0], answer: row[1]})
	}
	return problems, nil
}

func ask(problems []Problem, sdtin io.Reader, stdout io.Writer) int {
	correct := 0
	scanner := bufio.NewScanner(sdtin)
	for _, p := range problems {
		fmt.Fprintf(stdout, "%s\n", p.question)
		scanner.Scan()
		if scanner.Text() == p.answer {
			correct++
		}
	}
	return correct
}

var problemsFile string

func init() {
	flag.StringVar(&problemsFile, "file", "problems.csv", "file with containing questions/answers")
	flag.Parse()
}

func main() {
	file, err := os.Open(problemsFile)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	problems, err := parseProblems(file)
	if err != nil {
		panic(err)
	}
	correctAnswers := ask(problems, os.Stdin, os.Stdout)
	fmt.Printf("You answered correctly for %d/%d problems\n", correctAnswers, len(problems))

}
