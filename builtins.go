package welp

import "fmt"

type callable struct {
	name    string
	builtin bool

	// pointer to a built-in func if it's a builtin
	f func(env *environ, expr *node) *value

	// params and body of a user-defined func if it's no a builtin
	params *node
	body   *node
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

func sum(env *environ, expr *node) *value {
	lval := eval(env, expr)
	if lval.typ != valNum {
		fmt.Printf("Type error: unexpected type %s for +\n", lval.typ.String())
	}
	acc := lval.numValue
	for !nilNode(expr.r) {
		rval := eval(env, expr.r)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for +\n", rval.typ.String())
		}
		acc += rval.numValue
		expr = expr.r
	}
	return &value{
		typ:      valNum,
		numValue: acc,
	}
}

func sub(env *environ, expr *node) *value {
	lval := eval(env, expr)
	if lval.typ != valNum {
		fmt.Printf("Type error: unexpected type %s for -\n", lval.typ.String())
	}
	acc := lval.numValue
	for !nilNode(expr.r) {
		rval := eval(env, expr.r)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for -\n", rval.typ.String())
		}
		acc -= rval.numValue
		expr = expr.r
	}
	return &value{
		typ:      valNum,
		numValue: acc,
	}
}

func mul(env *environ, expr *node) *value {
	acc := num(env, expr.l.tok)
	for !nilNode(expr.r) {
		rval := eval(env, expr.r)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for *\n", rval.typ.String())
		}
		acc *= rval.numValue
		expr = expr.r
	}
	return &value{
		typ:      valNum,
		numValue: acc,
	}
}

// (exp base pow1 pow2 pow3) => base ^ (pow1 + pow2 + pow3)
func exp(env *environ, expr *node) *value {
	base := num(env, expr.l.tok)
	pow := 0
	for !nilNode(expr.r) {
		rval := eval(env, expr.r)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for exp\n", rval.typ.String())
		}
		pow += rval.numValue
		expr = expr.r
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
func Eval(env *environ, expr *node) *value {
	return eval(env, expr)
}

func eval(env *environ, expr *node) *value {
	if expr == nil || expr.l == nil {
		return &value{}
	}
	switch expr.l.tok.typ {
	case tokIdentifier:
		identName := string(expr.l.tok.value)
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
					return f.f(env, expr.r)
				}
				return callUserFunc(env, f, expr.r)
			}
		}
		if v, ok := env.vars[identName]; ok {
			return v
		}
		fmt.Printf("No such symbol %q\n", expr.l.tok)
	case tokNumber:
		return &value{
			typ:      valNum,
			numValue: num(env, expr.l.tok),
		}
	case tokVoid:
		return eval(env, expr.l)
	default:
		fmt.Printf("Unknown token type for %q\n", expr.l.tok)
	}
	return &value{}
}

// (eq 3 3) => T
// (eq 3 4) => NIL
func eq(env *environ, expr *node) *value {
	left := num(env, expr.l.tok)
	right := num(env, expr.r.l.tok)
	return &value{
		typ:       valBool,
		boolValue: left == right,
	}
}

// (cond
//    ((eq x 1) 1)
//    ((eq x 2) 1)
//    (t (fib (- x 1))))
func cond(env *environ, expr *node) *value {
	for expr.l != nil && expr.r.r != nil {
		conditional := eval(env, expr.l)
		if conditional.typ != valBool {
			fmt.Printf("Type error: cond clause evaluates to %s, not bool\n",
				conditional.typ.String())
		}
		if conditional.boolValue {
			return eval(env, expr.l.r)
		}
		expr = expr.r
	}
	return eval(env, expr.l.r)
}

// (fn add (a b) (+ a b)) => ADD
func defun(env *environ, expr *node) *value {
	funcName := string(expr.l.tok.value)
	params := expr.r.l
	body := expr.r.r.l
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
