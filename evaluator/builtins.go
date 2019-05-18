package evaluator

import (
	"errors"
	"fmt"

	"github.com/rtfb/welp/lexer"
	"github.com/rtfb/welp/object"
	"github.com/rtfb/welp/parser"
)

type callable struct {
	name    string
	builtin bool

	// pointer to a built-in func if it's a builtin
	f func(env *Environ, expr *parser.Node) object.Object

	// params and body of a user-defined func if it's no a builtin
	params *parser.Node
	body   *parser.Node
}

func makeBuiltins() map[string]*callable {
	return map[string]*callable{
		"+":        &callable{name: "+", f: sum, builtin: true},
		"-":        &callable{name: "-", f: sub, builtin: true},
		"*":        &callable{name: "*", f: mul, builtin: true},
		"exp":      &callable{name: "exp", f: exp, builtin: true},
		"eval":     &callable{name: "eval", f: eval, builtin: true},
		"fn":       &callable{name: "fn", f: defun, builtin: true},
		"cond":     &callable{name: "cond", f: cond, builtin: true},
		"eq":       &callable{name: "eq", f: eq, builtin: true},
		"let":      &callable{name: "let", f: let, builtin: true},
		"mk-array": &callable{name: "mk-array", f: makeArray, builtin: true},
		"append":   &callable{name: "append", f: arrAppend, builtin: true},
		"nth":      &callable{name: "nth", f: nth, builtin: true},
	}
}

func sum(env *Environ, expr *parser.Node) object.Object {
	lval := eval(env, expr)
	intLval, ok := lval.(*object.Integer)
	if !ok {
		fmt.Printf("Type error: unexpected type %T for +\n", lval)
	}
	acc := intLval.Value
	for !nilNode(expr.R) {
		rval := eval(env, expr.R)
		intRval, ok := rval.(*object.Integer)
		if !ok {
			fmt.Printf("Type error: unexpected type %T for +\n", rval)
		}
		acc += intRval.Value
		expr = expr.R
	}
	return &object.Integer{Value: acc}
}

func sub(env *Environ, expr *parser.Node) object.Object {
	lval := eval(env, expr)
	intLval, ok := lval.(*object.Integer)
	if !ok {
		fmt.Printf("Type error: unexpected type %T for -\n", lval)
	}
	acc := intLval.Value
	for !nilNode(expr.R) {
		rval := eval(env, expr.R)
		intRval, ok := rval.(*object.Integer)
		if !ok {
			fmt.Printf("Type error: unexpected type %T for -\n", rval)
		}
		acc -= intRval.Value
		expr = expr.R
	}
	return &object.Integer{Value: acc}
}

func mul(env *Environ, expr *parser.Node) object.Object {
	acc := num(env, expr.L.Tok)
	for !nilNode(expr.R) {
		rval := eval(env, expr.R)
		intRval, ok := rval.(*object.Integer)
		if !ok {
			fmt.Printf("Type error: unexpected type %T for *\n", rval)
		}
		acc *= intRval.Value
		expr = expr.R
	}
	return &object.Integer{Value: acc}
}

// (exp base pow1 pow2 pow3) => base ^ (pow1 + pow2 + pow3)
func exp(env *Environ, expr *parser.Node) object.Object {
	base := num(env, expr.L.Tok)
	pow := int64(0)
	for !nilNode(expr.R) {
		rval := eval(env, expr.R)
		intRval, ok := rval.(*object.Integer)
		if !ok {
			fmt.Printf("Type error: unexpected type %T for exp\n", rval)
		}
		pow += intRval.Value
		expr = expr.R
	}
	result := base
	for pow > 1 {
		result *= base
		pow--
	}
	return &object.Integer{Value: result}
}

// Eval evals.
func Eval(env *Environ, expr *parser.Node) object.Object {
	if expr.Err != nil {
		return &object.Error{Err: expr.Err}
	}
	return eval(env, expr)
}

func eval(env *Environ, expr *parser.Node) object.Object {
	if expr == nil || expr.L == nil {
		return &object.Null{}
	}
	switch expr.L.Tok.Typ {
	case lexer.TokIdentifier:
		identName := string(expr.L.Tok.Value)
		if identName == "t" {
			return &object.Boolean{Value: true}
		}
		if identName == "nil" {
			return &object.Boolean{Value: true}
		}
		function, ok := env.funcs[identName]
		if ok {
			if function.builtin {
				return function.f(env, expr.R)
			}
			return callUserFunc(env, function, expr.R)
		}
		if v, ok := env.vars[identName]; ok {
			return v
		}
		fmt.Printf("No such symbol %q\n", expr.L.Tok.String())
	case lexer.TokNumber:
		return &object.Integer{Value: num(env, expr.L.Tok)}
	case lexer.TokVoid:
		return eval(env, expr.L)
	default:
		fmt.Printf("Unknown token type for %q\n", expr.L.Tok.String())
	}
	return &object.Error{Err: errors.New("huh?")}
}

// (eq 3 3) => T
// (eq 3 4) => NIL
func eq(env *Environ, expr *parser.Node) object.Object {
	left := num(env, expr.L.Tok)
	right := num(env, expr.R.L.Tok)
	return &object.Boolean{Value: left == right}
}

// (cond
//    ((eq x 1) 1)
//    ((eq x 2) 1)
//    (t (fib (- x 1))))
func cond(env *Environ, expr *parser.Node) object.Object {
	for expr.L != nil && expr.R.R != nil {
		conditional := eval(env, expr.L)
		boolCond, ok := conditional.(*object.Boolean)
		if !ok {
			fmt.Printf("Type error: cond clause evaluates to %T, not bool\n",
				conditional)
		}
		if boolCond.Value {
			return eval(env, expr.L.R)
		}
		expr = expr.R
	}
	return eval(env, expr.L.R)
}

// (fn add (a b) (+ a b)) => ADD
func defun(env *Environ, expr *parser.Node) object.Object {
	funcName := string(expr.L.Tok.Value)
	params := expr.R.L
	body := expr.R.R.L
	env.funcs[funcName] = &callable{
		name:    funcName,
		builtin: false,
		params:  params,
		body:    body,
	}
	return &object.Func{Name: funcName}
}

// (let identifier (+ 2 3)) => 5
// identifier => 5
func let(env *Environ, expr *parser.Node) object.Object {
	value := eval(env, expr.R)
	env.vars[string(expr.L.Tok.Value)] = value
	return value
}

// (let arr (mk-array))
// arr => []
// TODO: find a way to specify the array type upfront.
func makeArray(env *Environ, expr *parser.Node) object.Object {
	return &object.Array{Value: nil}
}

// (let arr (mk-array))
// arr => []
// (append arr 3 5)
// arr => [3, 5]
func arrAppend(env *Environ, expr *parser.Node) object.Object {
	arrObj := eval(env, expr)
	if arrObj.Type() != object.ArrayType {
		return &object.Error{Err: fmt.Errorf("expected array, got %v", arrObj.Type())}
	}
	arr := arrObj.(*object.Array)
	arg := expr.R
	for arg.L != nil {
		value := eval(env, arg)
		if len(arr.Value) == 0 {
			arr.ValueType = value.Type()
		} else if value.Type() != arr.ValueType {
			return &object.Error{Err: fmt.Errorf("type mismatch: %v and %v",
				value.Type(), arr.ValueType)}
		}
		arr.Value = append(arr.Value, value)
		arg = arg.R
	}
	return arr
}

// (append arr 1 2 3)
// (nth 1 arr)
// => 2
func nth(env *Environ, expr *parser.Node) object.Object {
	indexObj := eval(env, expr)
	if indexObj.Type() != object.IntegerType {
		return &object.Error{Err: fmt.Errorf("type mismatch: %v and %v",
			indexObj.Type(), object.IntegerType)}
	}
	index := (indexObj.(*object.Integer)).Value
	arrObj := eval(env, expr.R)
	arr, ok := arrObj.(*object.Array)
	if !ok {
		return &object.Error{Err: fmt.Errorf("expected array, got %v", arrObj.Type())}
	}
	if index >= int64(len(arr.Value)) {
		return &object.Error{Err: fmt.Errorf("out of bounds: %d >= %d",
			index, len(arr.Value))}
	}
	return arr.Value[index]
}
