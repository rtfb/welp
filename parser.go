package welp

import (
	"fmt"
	"io"
	"io/ioutil"
)

type node struct {
	tok  token
	l, r *node
}

// Parser contains the state of the parser.
type Parser struct {
	tree   *node
	head   *node
	depth  int
	done   bool
	tokzer *tokenizer
}

func (p *Parser) rparse() {
	if p.done {
		return
	}
	var tok token
	for tok.typ != tokEOF {
		tok = <-p.tokzer.tok
		var treeNode *node
		switch tok.typ {
		case tokOpenParen:
			p.depth++
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
			p.depth--
			if p.depth == 0 {
				p.done = true
			}
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

// NewParser constructs a Parser.
func NewParser(input []byte) *Parser {
	return &Parser{
		tokzer: newTokenizer(input),
	}
}

// Start starts the concurrent part of the parser.
func (p *Parser) Start() {
	go p.tokzer.onStart()
}

// Parse parses source code into an expression tree.
func (p *Parser) Parse() (node *node, n int) {
	p.tree = nil
	p.done = false
	p.rparse()
	return p.tree, p.tokzer.head
}

// ParseString is a convenience func that parses a string.
func ParseString(input string) *node {
	p := NewParser([]byte(input))
	n, _ := p.Parse()
	return n
}

// ParseStream reads and parses all expressions from a given stream and sends
// them down the channel.
func ParseStream(r io.Reader) (<-chan *node, error) {
	allBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	p := NewParser(allBytes)
	p.Start()
	n := 0
	ch := make(chan *node)
	go func() {
		for n < len(allBytes) {
			node, newN := p.Parse()
			if newN == n {
				// TODO: fix error handling
				// return ch, errors.New("newN == n, can't progress")
				panic("newN == n, can't progress")
			}
			if node == nil {
				break
			}
			n = newN
			ch <- node
		}
		close(ch)
	}()
	return ch, nil
}
