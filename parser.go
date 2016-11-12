package main

import "fmt"

type node struct {
	tok  token
	l, r *node
}

type parser struct {
	tree      *node
	head      *node
	done      bool
	input     []byte
	inputHead int
}

func (p *parser) rparse() {
	if p.done {
		return
	}
	var tok token
	for tok != tokEOF {
		tok, p.inputHead = getToken(p.input, p.inputHead)
		var treeNode *node
		switch tok {
		case '(':
			var branchPoint *node
			if p.tree == nil {
				p.tree = &node{}
				p.head = p.tree
			} else {
				branchPoint = p.head
				p.head.l = &node{}
				p.head = p.head.l
			}
			p.rparse()
			if p.done {
				return
			}
			if branchPoint != nil {
				branchPoint.r = &node{}
				p.head = branchPoint.r
			}
			continue
		case ')':
			return
		case '+', '-', '*', '/', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			treeNode = &node{
				tok: tok,
			}
		case tokEOF:
			p.done = true
			return
		default:
			fmt.Printf("Unknown token: %c\n", rune(tok))
		}
		if p.head == nil {
			println("err1")
		}
		p.head.l = treeNode
		p.head.r = &node{}
		p.head = p.head.r
	}
}

func parse(input []byte) *node {
	p := &parser{
		input: input,
	}
	p.rparse()
	return p.tree
}
