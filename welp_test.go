package welp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEval(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"(+ 1 2 3 4)", 10},
		{"(* 2 (+ 3 7))", 20},
		{"(* 2 (+ 3 7) 5 9)", 900},
		{"(+ 123 321)", 444},
		{"(exp 2 3)", 8},
		{"(exp 2 3 2 4)", 512},
	}
	for _, test := range tests {
		env := NewEnv()
		got := eval(env, ParseString(test.input))
		if got.numValue != test.expected {
			assert.Equal(t, got, test.expected, "eval(%q)", test.input)
		}
	}
}
