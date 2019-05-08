package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer2(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"", []string{""}},
		{" ", []string{""}},
		{"(", []string{"("}},
		{"(* 2 (+ 3 7) 5 9)", []string{"(", "*", "2", "(", "+", "3", "7", ")", "5", "9", ")"}},
		{`(print "foo bar")`, []string{"(", "print", "foo bar", ")"}},
		{"(+ 123 321)", []string{"(", "+", "123", "321", ")"}},
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
			assert.Equal(t, test.expected[i], string(tok.Value))
			i++
		}
	}
}
