package evaluator

import (
	"testing"

	"github.com/rtfb/welp/object"
	"github.com/rtfb/welp/parser"
	"github.com/stretchr/testify/assert"
)

var testEvaluator = Evaluator{stdlibEnv: newEmptyEnv()}

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
		env := testEvaluator.NewEnv()
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
		env := testEvaluator.NewEnv()
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
		env := testEvaluator.NewEnv()
		got := eval(env, parser.ParseString(test.input))
		assert.IsType(t, &object.Boolean{}, got)
		boolGot := got.(*object.Boolean)
		assert.Equal(t, test.expected, boolGot.Value, "eval(%q)", test.input)
	}
}

func TestLet(t *testing.T) {
	env := testEvaluator.NewEnv()
	// assign something to x
	got := eval(env, parser.ParseString("(let x 5)"))
	assert.IsType(t, &object.Integer{}, got)
	intGot := got.(*object.Integer)
	assert.Equal(t, int64(5), intGot.Value)
	// now look it up: as an expr:
	got2 := eval(env, parser.ParseString("(x)"))
	assert.IsType(t, &object.Integer{}, got2)
	intGot2 := got2.(*object.Integer)
	assert.Equal(t, int64(5), intGot2.Value)
	// and as a standalone identifier:
	got3 := eval(env, parser.ParseString("x"))
	assert.IsType(t, &object.Integer{}, got3)
	intGot3 := got3.(*object.Integer)
	assert.Equal(t, int64(5), intGot3.Value)
}

func TestArrays(t *testing.T) {
	tests := []struct {
		input    string
		expected object.Object
	}{
		{"(mk-array)", &object.Array{}},
		{"(append (mk-array) 3 4)",
			&object.Array{
				ValueType: object.IntegerType,
				Value: []object.Object{
					&object.Integer{Value: 3},
					&object.Integer{Value: 4},
				},
			},
		},
		{"(nth 2 (append (mk-array) 3 (+ 2 3) 7))", &object.Integer{Value: 7}},
		{"(len (append (mk-array) 7 9))", &object.Integer{Value: 2}},
	}
	for _, test := range tests {
		env := testEvaluator.NewEnv()
		got := eval(env, parser.ParseString(test.input))
		assert.Equal(t, test.expected, got)
	}
}
