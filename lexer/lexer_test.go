package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"", []string{""}},
		{" ", []string{""}},
		{"(", []string{"("}},
		{"(* 2 (+ 3 7) 5 9)", []string{"(", "*", "2", "(", "+", "3", "7", ")", "5", "9", ")"}},
		{`(print "foo bar")`, []string{"(", "print", "foo bar", ")"}},
		{`(print "foo \"bar\"")`, []string{"(", "print", `foo "bar"`, ")"}},
		{"(+ 123 321)", []string{"(", "+", "123", "321", ")"}},
		{`(print "foo\\")`, []string{"(", "print", "foo\\", ")"}},
		{`(print "foo\n")`, []string{"(", "print", "foo\n", ")"}},
		{"x", []string{"x"}},
		{"9", []string{"9"}},
		{"19", []string{"19"}},
	}
	for _, test := range tests {
		tokzer := NewTokenizer(strings.NewReader(test.input))
		go tokzer.OnStart()
		i := 0
		for {
			tok := <-tokzer.Tok
			if tok.Typ == TokEOF {
				break
			}
			assert.NoError(t, tok.Err)
			assert.Equal(t, test.want[i], string(tok.Value))
			i++
		}
	}
}

func TestInvalidEscapeSequences2(t *testing.T) {
	tests := []struct {
		input   string
		wantErr string
	}{
		{input: `"foo\d"`, wantErr: "unrecognized escape sequence: \\d"},
		{input: `"foo\"`, wantErr: "unclosed string"},
		{input: `"foo\`, wantErr: "unclosed string"},
		{input: `"`, wantErr: "unclosed string"},
	}
	for _, test := range tests {
		tokzer := NewTokenizer(strings.NewReader(test.input))
		go tokzer.OnStart()
		tok := <-tokzer.Tok
		assert.EqualError(t, tok.Err, test.wantErr)
	}
}
