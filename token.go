package main

import (
	"bytes"
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
	case tokEOF:
		return "tokEOF"
	default:
		panic("unknown tokType")
	}
}

type token struct {
	typ   tokType
	value []byte
}

func (t *token) String() string {
	return fmt.Sprintf("{%s, %q}", t.typ.String(), t.value)
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
	if len(input) == 0 {
		return token{typ: tokEOF}, 1
	}
	for head < len(input) && input[head] == ' ' {
		head++
	}
	start := head
	for head < len(input) && strings.IndexByte(" \t()", input[head]) == -1 {
		head++
	}
	if head >= len(input) && len(input[start:head]) == 0 {
		return token{typ: tokEOF}, -1
	}
	if len(input[start:head]) == 0 {
		switch input[head] {
		case '(':
			return token{typ: tokOpenParen, value: []byte{input[head]}}, head + 1
		case ')':
			return token{typ: tokCloseParen, value: []byte{input[head]}}, head + 1
		default:
			return token{typ: tokEOF}, -1
		}
	}
	allDigits := true
	for _, c := range input[start:head] {
		if strings.IndexByte("0123456789", c) == -1 {
			allDigits = false
		}
	}
	tType := tokIdentifier
	if allDigits {
		tType = tokNumber
	}
	return token{
		typ:   tType,
		value: input[start:head],
	}, head
}

func tokEqual(a, b token) bool {
	if a.typ != b.typ {
		return false
	}
	return bytes.Equal(a.value, b.value)
}
