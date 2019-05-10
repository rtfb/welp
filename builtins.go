package welp

import (
	"fmt"

	"github.com/rtfb/welp/lexer"
	"github.com/rtfb/welp/parser"
)

type callable struct {
	name    string
	builtin bool

	// pointer to a built-in func if it's a builtin
	f func(env *Environ, expr *parser.Node) *value

	// params and body of a user-defined func if it's no a builtin
	params *parser.Node
	body   *parser.Node
}

var funcTbl []*callable

func init() {
	funcTbl = append(funcTbl, &callable{name: "+", f: sum, builtin: true})
	funcTbl = append(funcTbl, &callable{name: "-", f: sub, builtin: true})
	funcTbl = append(funcTbl, &callable{name: "*", f: mul, builtin: true})
	funcTbl = append(funcTbl, &callable{name: "exp", f: exp, builtin: true})
	funcTbl = append(funcTbl, &callable{name: "eval", f: eval, builtin: true})
	funcTbl = append(funcTbl, &callable{name: "fn", f: defun, builtin: true})
	funcTbl = append(funcTbl, &callable{name: "cond", f: cond, builtin: true})
	funcTbl = append(funcTbl, &callable{name: "eq", f: eq, builtin: true})
}

func sum(env *Environ, expr *parser.Node) *value {
	lval := eval(env, expr)
	if lval.typ != valNum {
		fmt.Printf("Type error: unexpected type %s for +\n", lval.typ.String())
	}
	acc := lval.numValue
	for !nilNode(expr.R) {
		rval := eval(env, expr.R)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for +\n", rval.typ.String())
		}
		acc += rval.numValue
		expr = expr.R
	}
	return &value{
		typ:      valNum,
		numValue: acc,
	}
}

func sub(env *Environ, expr *parser.Node) *value {
	lval := eval(env, expr)
	if lval.typ != valNum {
		fmt.Printf("Type error: unexpected type %s for -\n", lval.typ.String())
	}
	acc := lval.numValue
	for !nilNode(expr.R) {
		rval := eval(env, expr.R)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for -\n", rval.typ.String())
		}
		acc -= rval.numValue
		expr = expr.R
	}
	return &value{
		typ:      valNum,
		numValue: acc,
	}
}

func mul(env *Environ, expr *parser.Node) *value {
	acc := num(env, expr.L.Tok)
	for !nilNode(expr.R) {
		rval := eval(env, expr.R)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for *\n", rval.typ.String())
		}
		acc *= rval.numValue
		expr = expr.R
	}
	return &value{
		typ:      valNum,
		numValue: acc,
	}
}

// (exp base pow1 pow2 pow3) => base ^ (pow1 + pow2 + pow3)
func exp(env *Environ, expr *parser.Node) *value {
	base := num(env, expr.L.Tok)
	pow := 0
	for !nilNode(expr.R) {
		rval := eval(env, expr.R)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for exp\n", rval.typ.String())
		}
		pow += rval.numValue
		expr = expr.R
	}
	result := base
	for pow > 1 {
		result *= base
		pow--
	}
	return &value{
		typ:      valNum,
		numValue: result,
	}
}

// Eval evals.
func Eval(env *Environ, expr *parser.Node) *value {
	if expr.Err != nil {
		return newErrorValue(expr.Err)
	}
	return eval(env, expr)
}

func eval(env *Environ, expr *parser.Node) *value {
	if expr == nil || expr.L == nil {
		return &value{}
	}
	switch expr.L.Tok.Typ {
	case lexer.TokIdentifier:
		identName := string(expr.L.Tok.Value)
		if identName == "t" {
			return &value{
				typ:       valBool,
				boolValue: true,
			}
		}
		if identName == "nil" {
			return &value{
				typ:       valBool,
				boolValue: false,
			}
		}
		for _, f := range funcTbl {
			if identName == f.name {
				if f.builtin {
					return f.f(env, expr.R)
				}
				return callUserFunc(env, f, expr.R)
			}
		}
		if v, ok := env.vars[identName]; ok {
			return v
		}
		fmt.Printf("No such symbol %q\n", expr.L.Tok.String())
	case lexer.TokNumber:
		return &value{
			typ:      valNum,
			numValue: num(env, expr.L.Tok),
		}
	case lexer.TokVoid:
		return eval(env, expr.L)
	default:
		fmt.Printf("Unknown token type for %q\n", expr.L.Tok.String())
	}
	return &value{}
}

// (eq 3 3) => T
// (eq 3 4) => NIL
func eq(env *Environ, expr *parser.Node) *value {
	left := num(env, expr.L.Tok)
	right := num(env, expr.R.L.Tok)
	return &value{
		typ:       valBool,
		boolValue: left == right,
	}
}

// (cond
//    ((eq x 1) 1)
//    ((eq x 2) 1)
//    (t (fib (- x 1))))
func cond(env *Environ, expr *parser.Node) *value {
	for expr.L != nil && expr.R.R != nil {
		conditional := eval(env, expr.L)
		if conditional.typ != valBool {
			fmt.Printf("Type error: cond clause evaluates to %s, not bool\n",
				conditional.typ.String())
		}
		if conditional.boolValue {
			return eval(env, expr.L.R)
		}
		expr = expr.R
	}
	return eval(env, expr.L.R)
}

// (fn add (a b) (+ a b)) => ADD
func defun(env *Environ, expr *parser.Node) *value {
	funcName := string(expr.L.Tok.Value)
	params := expr.R.L
	body := expr.R.R.L
	funcTbl = append(funcTbl, &callable{
		name:    funcName,
		builtin: false,
		params:  params,
		body:    body,
	})
	return &value{
		typ:      valFunc,
		funcName: funcName,
	}
}
