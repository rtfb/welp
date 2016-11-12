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
		{"1+2", 1, token{typ: tokIdentifier, value: []byte("+")}},
		{"1+2", 2, token{typ: tokNumber, value: []byte("2")}},
		{"1+2", 3, token{typ: tokEOF}},
	}
	for _, test := range tests {
		got, _ := getToken([]byte(test.input), test.head)
		if !tokEqual(got, test.expected) {
			assert.Equal(t, got, test.expected, "getToken(%q, %d)",
				test.input, test.head)
		}
	}
}
