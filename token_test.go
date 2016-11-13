package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		input    string
		head     int
		expected token
	}{
		{"", 0, token{typ: tokEOF}},
		{"", 1, token{typ: tokEOF}},
		{" ", 0, token{typ: tokEOF}},
		{"(", 0, token{typ: tokOpenParen, value: []byte("(")}},
		{"1 + 2", 2, token{typ: tokIdentifier, value: []byte("+")}},
		{"1 + 2", 4, token{typ: tokNumber, value: []byte("2")}},
		{"1 + 2", 5, token{typ: tokEOF}},
	}
	for _, test := range tests {
		got, _ := getToken([]byte(test.input), test.head)
		if !tokEqual(got, test.expected) {
			assert.Equal(t, test.expected, got, "getToken(%q, %d)",
				test.input, test.head)
		}
	}
}

func TestTokenizer2(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"(* 2 (+ 3 7) 5 9)", []string{"(", "*", "2", "(", "+", "3", "7", ")", "5", "9", ")"}},
	}
	for _, test := range tests {
		head := 0
		var tok token
		i := 0
		for {
			tok, head = getToken([]byte(test.input), head)
			if tok.typ == tokEOF {
				break
			}
			assert.Equal(t, test.expected[i], string(tok.value))
			i++
		}
	}
}
