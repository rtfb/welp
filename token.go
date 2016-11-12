package main

import "bytes"

type tokType int

type token struct {
	typ   tokType
	value []byte
}

const (
	tokVoid tokType = iota
	tokOpenParen
	tokCloseParen
	tokNumber
	tokIdentifier
	tokEOF
)

func getToken(input []byte, head int) (tok token, newHead int) {
	for head < len(input) && input[head] == ' ' {
		head++
	}
	if head < len(input) {
		var tType tokType
		switch input[head] {
		case '(':
			tType = tokOpenParen
		case ')':
			tType = tokCloseParen
		case '+', '*', '-', '/':
			tType = tokIdentifier
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			tType = tokNumber
		}
		return token{
			typ:   tType,
			value: []byte{input[head]},
		}, head + 1
	}
	return token{
		typ: tokEOF,
	}, head + 1
}

func tokEqual(a, b token) bool {
	if a.typ != b.typ {
		return false
	}
	return bytes.Equal(a.value, b.value)
}
