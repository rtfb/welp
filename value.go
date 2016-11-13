package main

import (
	"strconv"
	"strings"
)

type valType int

func (v valType) String() string {
	switch v {
	case valNum:
		return "integer"
	case valUndef:
		return "undefined"
	default:
		panic("unknown value type")
	}
}

const (
	valUndef valType = iota
	valNum
	valFunc
)

type value struct {
	typ      valType
	numValue int
	funcName string
}

func (v *value) String() string {
	switch v.typ {
	case valNum:
		return strconv.Itoa(v.numValue)
	case valFunc:
		return strings.ToUpper(v.funcName)
	default:
		return "<UNK>"
	}
}
