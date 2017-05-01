package welp

import (
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
		tokzer := newTokenizer([]byte(test.input))
		go tokzer.onStart()
		i := 0
		for {
			tok := <-tokzer.tok
			if tok.typ == tokEOF {
				break
			}
			assert.Equal(t, test.expected[i], string(tok.value))
			i++
		}
	}
}
