package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

var (
	problemsFile string
	timeout      time.Duration
	doShuffling  bool
	seed         int64
)

func init() {
	seed = time.Now().Unix()
	flag.StringVar(&problemsFile, "file", "problems.csv", "file with containing questions/answers")
	flag.DurationVar(&timeout, "timeout", time.Second*30, "time after which quiz ends")
	flag.BoolVar(&doShuffling, "shuffle", false, "if true, questions will be shuffled")
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

	if doShuffling {
		problems = shuffle(problems)
	}

	correct := make(chan bool)
	total := 0

	fmt.Println("If you are ready press enter")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	timer := time.NewTimer(timeout)
	go ask(problems, os.Stdin, os.Stdout, correct)

	for n := 1; n > 0; {
		select {
		case v, ok := <-correct:
			if !ok {
				n--
			}
			if v {
				total++
			}
		case <-timer.C:
			n--
		}
	}
	fmt.Printf("You answered correctly %d/%d problems\n", total, len(problems))
}

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

func ask(problems []Problem, sdtin io.Reader, stdout io.Writer, correct chan bool) {
	scanner := bufio.NewScanner(sdtin)
	for _, p := range problems {
		fmt.Fprintf(stdout, "%s\n", p.question)
		scanner.Scan()
		if scanner.Text() == p.answer {
			correct <- true
		}
	}
	close(correct)
}

func shuffle(slice []Problem) []Problem {
	r := rand.New(rand.NewSource(seed))
	ret := make([]Problem, len(slice))
	perm := r.Perm(len(slice))
	for i, randIndex := range perm {
		ret[i] = slice[randIndex]
	}
	return ret
}
