package main

import (
	"fmt"
	"strconv"
	"strings"
)

func num(tok token) int {
	n, err := strconv.Atoi(string(tok.value))
	if err != nil {
		panic(err)
	}
	return n
}

// TODO fix parser to produce an actual nil instead of this node
func nilNode(n *node) bool {
	return n.tok.typ == tokVoid && n.l == nil && n.r == nil
}

func sum(ast *node) int {
	acc := num(ast.l.tok)
	for !nilNode(ast.r) {
		acc += eval(ast.r)
		ast = ast.r
	}
	return acc
}

func mul(ast *node) int {
	acc := num(ast.l.tok)
	for !nilNode(ast.r) {
		acc *= eval(ast.r)
		ast = ast.r
	}
	return acc
}

// (exp base pow1 pow2 pow3) => base ^ (pow1 + pow2 + pow3)
func exp(ast *node) int {
	base := num(ast.l.tok)
	pow := 0
	for !nilNode(ast.r) {
		pow += eval(ast.r)
		ast = ast.r
	}
	result := base
	for pow > 1 {
		result *= base
		pow--
	}
	return result
}

func eval(ast *node) int {
	if ast == nil || ast.l == nil {
		return 0
	}
	switch ast.l.tok.typ {
	case tokIdentifier:
		switch string(ast.l.tok.value) {
		case "+":
			return sum(ast.r)
		case "*":
			return mul(ast.r)
		case "exp":
			return exp(ast.r)
		}
	case tokNumber:
		return num(ast.l.tok)
	case tokVoid:
		return eval(ast.l)
	default:
		fmt.Printf("No such func %q\n", ast.l.tok)
	}
	return -1
}

var indent int

func dump(ast *node) {
	indent += 4
	prefix := strings.Repeat(" ", indent)
	if ast.tok.typ == tokVoid {
		fmt.Printf("%s-\n", prefix)
	} else {
		fmt.Printf("%stok: %s\n", prefix, ast.tok.value)
	}
	if ast.l != nil {
		dump(ast.l)
	}
	if ast.r != nil {
		dump(ast.r)
	}
	indent -= 4
}

func main() {
	ast := parse([]byte("(* 2 (+ 3 7) 5 9)"))
	dump(ast)
	println(eval(ast))
}
