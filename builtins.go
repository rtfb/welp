package main

import "fmt"

type callable struct {
	name    string
	builtin bool

	// pointer to a built-in func if it's a builtin
	f func(env *environ, ast *node) *value

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
	funcTbl = append(funcTbl, &callable{name: "t", f: t, builtin: true})
}

func sum(env *environ, ast *node) *value {
	lval := eval(env, ast)
	if lval.typ != valNum {
		fmt.Printf("Type error: unexpected type %s for +\n", lval.typ.String())
	}
	acc := lval.numValue
	for !nilNode(ast.r) {
		rval := eval(env, ast.r)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for +\n", rval.typ.String())
		}
		acc += rval.numValue
		ast = ast.r
	}
	return &value{
		typ:      valNum,
		numValue: acc,
	}
}

func sub(env *environ, ast *node) *value {
	lval := eval(env, ast)
	if lval.typ != valNum {
		fmt.Printf("Type error: unexpected type %s for -\n", lval.typ.String())
	}
	acc := lval.numValue
	for !nilNode(ast.r) {
		rval := eval(env, ast.r)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for -\n", rval.typ.String())
		}
		acc -= rval.numValue
		ast = ast.r
	}
	return &value{
		typ:      valNum,
		numValue: acc,
	}
}

func mul(env *environ, ast *node) *value {
	acc := num(env, ast.l.tok)
	for !nilNode(ast.r) {
		rval := eval(env, ast.r)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for *\n", rval.typ.String())
		}
		acc *= rval.numValue
		ast = ast.r
	}
	return &value{
		typ:      valNum,
		numValue: acc,
	}
}

// (exp base pow1 pow2 pow3) => base ^ (pow1 + pow2 + pow3)
func exp(env *environ, ast *node) *value {
	base := num(env, ast.l.tok)
	pow := 0
	for !nilNode(ast.r) {
		rval := eval(env, ast.r)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for exp\n", rval.typ.String())
		}
		pow += rval.numValue
		ast = ast.r
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

func eval(env *environ, ast *node) *value {
	if ast == nil || ast.l == nil {
		return &value{}
	}
	switch ast.l.tok.typ {
	case tokIdentifier:
		identName := string(ast.l.tok.value)
		for _, f := range funcTbl {
			if identName == f.name {
				if f.builtin {
					return f.f(env, ast.r)
				}
				return callUserFunc(env, f, ast.r)
			}
		}
		if v, ok := env.vars[identName]; ok {
			return v
		}
		fmt.Printf("No such symbol %q\n", ast.l.tok)
	case tokNumber:
		return &value{
			typ:      valNum,
			numValue: num(env, ast.l.tok),
		}
	case tokVoid:
		return eval(env, ast.l)
	default:
		fmt.Printf("Unknown token type for %q\n", ast.l.tok)
	}
	return &value{}
}

// (eq 3 3) => T
// (eq 3 4) => NIL
func eq(env *environ, ast *node) *value {
	left := num(env, ast.l.tok)
	right := num(env, ast.r.l.tok)
	return &value{
		typ:       valBool,
		boolValue: left == right,
	}
}

// (t (+ 4 6)) => 10
func t(env *environ, ast *node) *value {
	return eval(env, ast)
}

// (cond
//    ((eq x 1) 1)
//    ((eq x 2) 1)
//    (t (fib (- x 1))))
func cond(env *environ, ast *node) *value {
	for ast.l != nil && ast.r.r != nil {
		conditional := eval(env, ast.l)
		if conditional.typ != valBool {
			fmt.Printf("Type error: cond clause evaluates to %s, not bool\n",
				conditional.typ.String())
		}
		if conditional.boolValue {
			return eval(env, ast.l.r)
		}
		ast = ast.r
	}
	return eval(env, ast)
}

// (fn add (a b) (+ a b)) => ADD
func defun(env *environ, ast *node) *value {
	funcName := string(ast.l.tok.value)
	params := ast.r.l
	body := ast.r.r.l
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
