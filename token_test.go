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
		{"", 0, tokEOF},
		{"", 1, tokEOF},
		{" ", 0, tokEOF},
		{"(", 0, '('},
		{"1+2", 1, '+'},
		{"1+2", 2, '2'},
		{"1+2", 3, tokEOF},
	}
	for _, test := range tests {
		got, _ := getToken([]byte(test.input), test.head)
		if got != test.expected {
			assert.Equal(t, got, test.expected, "getToken(%q, %d)",
				test.input, test.head)
		}
	}
}
