package welp

import (
	"fmt"
	"strings"
)

type tokType int

func (t tokType) String() string {
	switch t {
	case tokVoid:
		return "tokVoid"
	case tokOpenParen:
		return "tokOpenParen"
	case tokCloseParen:
		return "tokCloseParen"
	case tokNumber:
		return "tokNumber"
	case tokIdentifier:
		return "tokIdentifier"
	case tokString:
		return "tokString"
	case tokEOF:
		return "tokEOF"
	default:
		panic("unknown tokType")
	}
}

type token struct {
	typ   tokType
	value []byte
	pos   int
}

func (t *token) String() string {
	return fmt.Sprintf("{%s, %q (P%d)}", t.typ.String(), t.value, t.pos)
}

const (
	tokVoid tokType = iota
	tokOpenParen
	tokCloseParen
	tokNumber
	tokIdentifier
	tokString
	tokEOF
)

type tokenizer struct {
	input []byte
	head  int
	tok   chan token
}

func newTokenizer(input []byte) *tokenizer {
	return &tokenizer{
		input: input,
		tok:   make(chan token),
	}
}

func (t *tokenizer) onStart() {
	for t.head < len(t.input) {
		h := t.input[t.head]
		switch {
		case h >= '0' && h <= '9':
			t.onNumber()
		case h == ' ' || h == '\t' || h == '\n':
			t.head++
		case h == '(':
			t.onOpenParen()
		case h == ')':
			t.onCloseParen()
		case h == '"':
			t.onDoublequote()
		default:
			t.onChar()
		}
	}
	t.tok <- token{typ: tokEOF, pos: len(t.input)}
}

func (t *tokenizer) onNumber() {
	var tail int
	var h byte
	for tail, h = range t.input[t.head:] {
		if h < '0' || h > '9' {
			if h != '.' {
				break
			}
		}
	}
	tail = t.head + tail
	t.tok <- token{typ: tokNumber, value: t.input[t.head:tail], pos: t.head}
	t.head = tail
}

func (t *tokenizer) onChar() {
	var tail int
	var h byte
	for tail, h = range t.input[t.head:] {
		if strings.IndexByte(" \n\t()\"", h) != -1 {
			break
		}
	}
	tail = t.head + tail
	t.tok <- token{typ: tokIdentifier, value: t.input[t.head:tail], pos: t.head}
	t.head = tail
}

func (t *tokenizer) onOpenParen() {
	t.tok <- token{typ: tokOpenParen, value: []byte{'('}, pos: t.head}
	t.head++
}

func (t *tokenizer) onCloseParen() {
	t.tok <- token{typ: tokCloseParen, value: []byte{')'}, pos: t.head}
	t.head++
}

func (t *tokenizer) onDoublequote() {
	var tail int
	var h byte
	t.head++
	for tail, h = range t.input[t.head:] {
		if h == '"' {
			if tail > 1 && t.input[tail-1] == '\\' {
				continue
			}
			break
		}
	}
	tail = t.head + tail
	t.tok <- token{typ: tokString, value: t.input[t.head:tail], pos: t.head}
	t.head = tail + 1
}
