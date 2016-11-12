package main

import (
	"fmt"
	"strings"
)

func num(tok token) int {
	return int(tok - '0')
}

// TODO fix parser to produce an actual nil instead of this node
func nilNode(n *node) bool {
	return n.tok == tokVoid && n.l == nil && n.r == nil
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

func eval(ast *node) int {
	switch ast.l.tok {
	case '+':
		return sum(ast.r)
	case '*':
		return mul(ast.r)
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
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
	if ast.tok == tokVoid {
		fmt.Printf("%s-\n", prefix)
	} else {
		fmt.Printf("%stok: %c\n", prefix, rune(ast.tok))
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
