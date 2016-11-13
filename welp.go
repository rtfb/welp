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

// (add 3 7) => 10
func callUserFunc(f *callable, ast *node) *value {
	// TODO: add checking. At least check if number of args is correct
	f.body.r = ast
	return eval(f.body)
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
	println(eval(ast).String())
	ast = parse([]byte("(fn add (a b) (+ a b))"))
	println(eval(ast).String())
	ast = parse([]byte("(add 3 7)"))
	println(eval(ast).String())
	ast = parse([]byte("(add 34 79)"))
	println(eval(ast).String())
}
