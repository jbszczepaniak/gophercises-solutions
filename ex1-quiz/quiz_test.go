package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestParseProblems(t *testing.T) {
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
