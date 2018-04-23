package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseProblems(t *testing.T) {
	t.Parallel()

	cases := []struct {
		content  string
		expected []Problem
	}{
		{
			"",
			[]Problem{},
		},
		{
			"1+1,2\n",
			[]Problem{{"1+1", "2"}},
		},
		{
			"1+1,2\n2+3,5\n",
			[]Problem{{"1+1", "2"}, {"2+3", "5"}},
		},
	}
	for _, c := range cases {
		problems, err := parseProblems(bytes.NewBufferString(c.content))
		if err != nil {
			t.Fatalf("did not expect to fail, but it did with %s", err)
		}
		if !reflect.DeepEqual(problems, c.expected) {
			t.Fatalf("want %v, got %v", c.expected, problems)
		}
	}
}

func TestAsk(t *testing.T) {
	t.Parallel()

	cases := []struct {
		problems []Problem
		answers  string
		correct  int
	}{
		{
			[]Problem{{"1+1", "2"}, {"1+2", "3"}},
			"2\n3\n",
			2,
		},
		{
			[]Problem{{"1+1", "2"}, {"1+2", "3"}},
			"",
			0,
		},
		{
			[]Problem{{"1+1", "2"}, {"1+2", "3"}},
			"2\n",
			1,
		},
		{
			[]Problem{{"niechęć?", "tak"}, {"łaskotać?", "nie"}},
			"tak\nnie\n",
			2,
		},
	}

	for _, c := range cases {
		out := bytes.Buffer{}
		in := bytes.NewBufferString(c.answers)
		correctAnswer := make(chan bool)

		go ask(c.problems, in, &out, correctAnswer)
		totalCorrect := 0
		for range correctAnswer {
			totalCorrect++
		}

		if totalCorrect != c.correct {
			t.Errorf("want %d correct answers, have %d", c.correct, totalCorrect)
		}
	}
}

func TestShuffleChangesOrder(t *testing.T) {
	seed = 2213123 // Some seed value that gives changed order of slice elements
	problems := []Problem{
		{question: "?", answer: "!"},
		{question: "yes?", answer: "no!"},
		{question: "2+2", answer: "4"},
	}
	isShuffled := false
	shuffled := shuffle(problems)

	for i := 0; i < len(shuffled); i++ {
		if problems[i] != shuffled[i] {
			isShuffled = true
			return
		}
	}
	if !isShuffled {
		t.Errorf("expected slice to be shuffled but it didn't")
	}
}
