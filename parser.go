package welp

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// Node defines a single node in the s-expression.
type Node struct {
	Tok  token
	L, R *Node
	Err  error
}

// Parser contains the state of the parser.
type Parser struct {
	tree   *Node
	head   *Node
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
		if tok.err != nil && tok.typ != tokEOF {
			panic(tok.err) // TODO: improve error handling
		}
		var treeNode *Node
		switch tok.typ {
		case tokOpenParen:
			p.depth++
			var branchPoint *Node
			if p.tree == nil {
				p.tree = &Node{}
				p.head = p.tree
			} else {
				branchPoint = p.head
				p.head.L = &Node{}
				p.head = p.head.L
			}
			p.rparse()
			if p.done {
				return
			}
			if branchPoint != nil {
				branchPoint.R = &Node{}
				p.head = branchPoint.R
			}
			continue
		case tokCloseParen:
			p.depth--
			if p.depth == 0 {
				p.done = true
			}
			return
		case tokIdentifier, tokNumber:
			treeNode = &Node{
				Tok: tok,
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
		p.head.L = treeNode
		p.head.R = &Node{}
		p.head = p.head.R
	}
}

// NewParser constructs a Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{
		tokzer: newTokenizer(r),
	}
}

// Start starts the concurrent part of the parser.
func (p *Parser) Start() {
	go p.tokzer.onStart()
}

// Parse parses source code into an expression tree.
func (p *Parser) Parse() (node *Node, n int) {
	p.Reset()
	p.rparse()
	return p.tree, p.tokzer.head
}

// Reset clears parser state.
func (p *Parser) Reset() {
	p.tree = nil
	p.done = false
}

// ParseString is a convenience func that parses a string.
func ParseString(input string) *Node {
	p := NewParser(strings.NewReader(input))
	n, _ := p.Parse()
	return n
}

// ParseStream reads and parses all expressions from a given stream and sends
// them down the channel.
func ParseStream(r io.Reader) <-chan *Node {
	p := NewParser(r)
	p.Start()
	n := 0
	ch := make(chan *Node)
	go func() {
		for {
			node, newN := p.Parse()
			if newN == n {
				ch <- &Node{
					Err: errors.New("newN == n, can't progress"),
				}
				break
			}
			if node == nil {
				break
			}
			n = newN
			ch <- node
		}
		close(ch)
	}()
	return ch
}
