package welp

import (
	"testing"

	"github.com/rtfb/welp/object"
	"github.com/rtfb/welp/parser"
	"github.com/stretchr/testify/assert"
)

func TestEvalMath(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"(+ 1 2 3 4)", 10},
		{"(* 2 (+ 3 7))", 20},
		{"(* 2 (+ 3 7) 5 9)", 900},
		{"(+ 123 321)", 444},
		{"(exp 2 3)", 8},
		{"(exp 2 3 2 4)", 512},
		{"(- 7 5)", 2},
	}
	for _, test := range tests {
		env := NewEnv()
		got := eval(env, parser.ParseString(test.input))
		assert.IsType(t, &object.Integer{}, got)
		intGot := got.(*object.Integer)
		assert.Equal(t, test.expected, intGot.Value, "eval(%q)", test.input)
	}
}

func TestEvalEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"", 0},
	}
	for _, test := range tests {
		env := NewEnv()
		got := eval(env, parser.ParseString(test.input))
		assert.IsType(t, &object.Null{}, got)
	}
}

func TestEvalEq(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"(eq 3 3)", true},
		{"(eq 3 4)", false},
	}
	for _, test := range tests {
		env := NewEnv()
		got := eval(env, parser.ParseString(test.input))
		assert.IsType(t, &object.Boolean{}, got)
		boolGot := got.(*object.Boolean)
		assert.Equal(t, test.expected, boolGot.Value, "eval(%q)", test.input)
	}
}
