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
	for tok.typ != tokEOF {
		tok, p.inputHead = getToken(p.input, p.inputHead)
		var treeNode *node
		switch tok.typ {
		case tokOpenParen:
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
		case tokCloseParen:
			return
		case tokIdentifier, tokNumber:
			treeNode = &node{
				tok: tok,
			}
		case tokEOF:
			p.done = true
			return
		default:
			fmt.Printf("Unknown token: %v\n", tok)
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

func parseS(input string) *node {
	return parse([]byte(input))
}
