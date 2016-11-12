package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		// {"", 0},
		{"(+ 1 2 3 4)", 10},
		{"(* 2 (+ 3 7))", 20},
		{"(* 2 (+ 3 7) 5 9)", 900},
	}
	for _, test := range tests {
		got := eval(parse([]byte(test.input)))
		if got != test.expected {
			assert.Equal(t, got, test.expected, "eval(%q)", test.input)
		}
	}
}
