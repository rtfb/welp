package parser

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/rtfb/welp/lexer"
)

// Node defines a single node in the s-expression.
type Node struct {
	Tok  lexer.Token
	L, R *Node
	Err  error
}

// Parser contains the state of the parser.
type Parser struct {
	tree   *Node
	head   *Node
	depth  int
	done   bool
	tokzer *lexer.Tokenizer
	debug  bool
}

// New constructs a Parser.
func New(r io.Reader) *Parser {
	return &Parser{
		tokzer: lexer.NewTokenizer(r),
		debug:  false,
	}
}

func (p *Parser) rparse() {
	if p.done {
		return
	}
	var tok lexer.Token
	for tok.Typ != lexer.TokEOF {
		tok = <-p.tokzer.Tok
		if p.debug {
			fmt.Println(&tok)
		}
		if tok.Err != nil && tok.Typ != lexer.TokEOF {
			panic(tok.Err) // TODO: improve error handling
		}
		var treeNode *Node
		switch tok.Typ {
		case lexer.TokOpenParen:
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
		case lexer.TokCloseParen:
			p.depth--
			if p.depth == 0 {
				p.done = true
			}
			return
		case lexer.TokIdentifier, lexer.TokNumber:
			treeNode = &Node{
				Tok: tok,
			}
			if p.depth == 0 {
				// this is a special case for a standalone identifier in REPL:
				// welp> (let x 7)
				// welp> x
				// 7
				p.tree = &Node{L: treeNode}
				p.done = true
				return
			}
		case lexer.TokEOF:
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

// Start starts the concurrent part of the parser.
func (p *Parser) Start() {
	go p.tokzer.OnStart()
}

// Parse parses source code into an expression tree.
func (p *Parser) Parse() (node *Node, n int) {
	p.Reset()
	p.rparse()
	return p.tree, p.tokzer.Head
}

// Reset clears parser state.
func (p *Parser) Reset() {
	p.tree = nil
	p.head = nil
	p.done = false
}

// ParseString is a convenience func that parses a string.
func ParseString(input string) *Node {
	p := New(strings.NewReader(input))
	p.Start()
	n, _ := p.Parse()
	return n
}

// ParseStream reads and parses all expressions from a given stream and sends
// them down the channel.
func ParseStream(r io.Reader) <-chan *Node {
	p := New(r)
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
