package main

import "fmt"

type callable struct {
	name    string
	builtin bool

	// pointer to a built-in func if it's a builtin
	f func(ast *node) *value

	// params and body of a user-defined func if it's no a builtin
	params *node
	body   *node
}

var funcTbl []*callable

func init() {
	funcTbl = append(funcTbl, &callable{name: "+", f: sum, builtin: true})
	funcTbl = append(funcTbl, &callable{name: "*", f: mul, builtin: true})
	funcTbl = append(funcTbl, &callable{name: "exp", f: exp, builtin: true})
	funcTbl = append(funcTbl, &callable{name: "eval", f: eval, builtin: true})
	funcTbl = append(funcTbl, &callable{name: "fn", f: defun, builtin: true})
}

func sum(ast *node) *value {
	acc := num(ast.l.tok)
	for !nilNode(ast.r) {
		rval := eval(ast.r)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for +", rval.typ.String())
		}
		acc += rval.numValue
		ast = ast.r
	}
	return &value{
		typ:      valNum,
		numValue: acc,
	}
}

func mul(ast *node) *value {
	acc := num(ast.l.tok)
	for !nilNode(ast.r) {
		rval := eval(ast.r)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for *", rval.typ.String())
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
func exp(ast *node) *value {
	base := num(ast.l.tok)
	pow := 0
	for !nilNode(ast.r) {
		rval := eval(ast.r)
		if rval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for exp", rval.typ.String())
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

func eval(ast *node) *value {
	if ast == nil || ast.l == nil {
		return &value{}
	}
	switch ast.l.tok.typ {
	case tokIdentifier:
		for _, f := range funcTbl {
			if string(ast.l.tok.value) == f.name {
				if f.builtin {
					return f.f(ast.r)
				}
				return callUserFunc(f, ast.r)
			}
		}
		fmt.Printf("No such func %q\n", ast.l.tok)
	case tokNumber:
		return &value{
			typ:      valNum,
			numValue: num(ast.l.tok),
		}
	case tokVoid:
		return eval(ast.l)
	default:
		fmt.Printf("No such func %q\n", ast.l.tok)
	}
	return &value{}
}

// (fn add (a b) (+ a b)) => ADD
func defun(ast *node) *value {
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
