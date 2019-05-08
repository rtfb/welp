package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// TokType defines the type of a token.
type TokType int

// These are the individual token type constants.
const (
	TokVoid TokType = iota
	TokOpenParen
	TokCloseParen
	TokNumber
	TokIdentifier
	TokString
	TokEOF
)

// String implements Stringer.
func (t TokType) String() string {
	switch t {
	case TokVoid:
		return "tokVoid"
	case TokOpenParen:
		return "TokOpenParen"
	case TokCloseParen:
		return "TokCloseParen"
	case TokNumber:
		return "TokNumber"
	case TokIdentifier:
		return "TokIdentifier"
	case TokString:
		return "TokString"
	case TokEOF:
		return "TokEOF"
	default:
		panic("unknown TokType")
	}
}

// Token represents a specific token discovered by the lexer.
type Token struct {
	Typ   TokType
	Value []byte
	Pos   int
	Err   error
}

// String implements Stringer.
func (t *Token) String() string {
	return fmt.Sprintf("{%s, %q (P%d)}", t.Typ.String(), t.Value, t.Pos)
}

// Tokenizer is our lexer.
type Tokenizer struct {
	r    *bufio.Reader
	Head int
	Tok  chan Token
}

// NewTokenizer creates a lexer with a reader to read data from.
func NewTokenizer(r io.Reader) *Tokenizer {
	return &Tokenizer{
		r:   bufio.NewReader(r),
		Tok: make(chan Token),
	}
}

// OnStart runs the actual lexing process, meant to be launched in a separate
// goroutine.
func (t *Tokenizer) OnStart() {
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
			t.Head++
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
	t.Tok <- Token{Typ: TokEOF, Pos: t.Head, Err: err}
}

func (t *Tokenizer) onNumber() {
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
		buf.WriteByte(b)
	}
	t.Tok <- Token{Typ: TokNumber, Value: buf.Bytes(), Pos: t.Head, Err: err}
	t.Head += buf.Len()
}

func (t *Tokenizer) onChar() {
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
		buf.WriteByte(b)
	}
	t.Tok <- Token{Typ: TokIdentifier, Value: buf.Bytes(), Pos: t.Head, Err: err}
	t.Head += buf.Len()
}

func (t *Tokenizer) onOpenParen() {
	t.Tok <- Token{Typ: TokOpenParen, Value: []byte{'('}, Pos: t.Head}
	t.Head++
}

func (t *Tokenizer) onCloseParen() {
	t.Tok <- Token{Typ: TokCloseParen, Value: []byte{')'}, Pos: t.Head}
	t.Head++
}

func (t *Tokenizer) onDoublequote() {
	var buf bytes.Buffer
	var b byte
	var err error
	foundClosingDoublequote := false
loop:
	for {
		b, err = t.r.ReadByte()
		if err != nil {
			break
		}
		if b == '\\' {
			var bAfter byte
			bAfter, err = t.r.ReadByte()
			if err != nil {
				break
			}
			switch bAfter {
			case '"':
				buf.WriteByte(bAfter)
				continue
			case '\\':
				buf.WriteByte(bAfter)
				continue
			case 'n':
				buf.WriteByte('\n')
				continue
			default:
				t.r.UnreadByte()
				err = fmt.Errorf("unrecognized escape sequence: \\%c", bAfter)
				break loop
			}
		}
		if b == '"' {
			foundClosingDoublequote = true
			break
		}
		buf.WriteByte(b)
	}
	if (err == nil || err == io.EOF) && !foundClosingDoublequote {
		err = fmt.Errorf("unclosed string")
	}
	t.Tok <- Token{Typ: TokString, Value: buf.Bytes(), Pos: t.Head, Err: err}
	t.Head += buf.Len()
}
