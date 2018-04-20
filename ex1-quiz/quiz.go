package main

import (
	"encoding/csv"
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

func main() {
	file, err := os.Open("problems.csv")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	parseProblems(file)
}
