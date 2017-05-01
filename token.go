package welp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
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
	err   error
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
	r    *bufio.Reader
	head int
	tok  chan token
}

func newTokenizer(r io.Reader) *tokenizer {
	return &tokenizer{
		r:   bufio.NewReader(r),
		tok: make(chan token),
	}
}

func (t *tokenizer) onStart() {
	var err error
	var b byte
	for {
		b, err = t.r.ReadByte()
		if err != nil {
			break
		}
		switch {
		case b >= '0' && b <= '9':
			t.r.UnreadByte()
			t.onNumber()
		case b == ' ' || b == '\t' || b == '\n' || b == '\r':
			t.head++
		case b == '(':
			t.onOpenParen()
		case b == ')':
			t.onCloseParen()
		case b == '"':
			t.onDoublequote()
		default:
			t.r.UnreadByte()
			t.onChar()
		}
	}
	t.tok <- token{typ: tokEOF, pos: t.head, err: err}
}

func (t *tokenizer) onNumber() {
	var buf bytes.Buffer
	var b byte
	var err error
	for {
		b, err = t.r.ReadByte()
		if err != nil {
			break
		}
		if (b < '0' || b > '9') && b != '.' {
			t.r.UnreadByte()
			break
		}
		err = buf.WriteByte(b)
	}
	t.tok <- token{typ: tokNumber, value: buf.Bytes(), pos: t.head, err: err}
	t.head += buf.Len()
}

func (t *tokenizer) onChar() {
	var buf bytes.Buffer
	var b byte
	var err error
	for {
		b, err = t.r.ReadByte()
		if err != nil {
			break
		}
		if strings.IndexByte(" \n\t()\"", b) != -1 {
			t.r.UnreadByte()
			break
		}
		err = buf.WriteByte(b)
	}
	t.tok <- token{typ: tokIdentifier, value: buf.Bytes(), pos: t.head, err: err}
	t.head += buf.Len()
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
	var buf bytes.Buffer
	var b, lastB byte
	var err error
	for {
		b, err = t.r.ReadByte()
		if err != nil {
			break
		}
		if b == '"' {
			if lastB == '\\' {
				continue
			}
			t.r.UnreadByte()
			break
		}
		lastB = b
		err = buf.WriteByte(b)
	}
	t.tok <- token{typ: tokString, value: buf.Bytes(), pos: t.head, err: err}
	t.head += buf.Len()
}
