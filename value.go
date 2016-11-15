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
	case valFunc:
		return "function"
	case valBool:
		return "bool"
	default:
		panic("unknown value type")
	}
}

const (
	valUndef valType = iota
	valNum
	valFunc
	valBool
)

type value struct {
	typ       valType
	numValue  int
	boolValue bool
	funcName  string
}

func (v *value) String() string {
	switch v.typ {
	case valNum:
		return strconv.Itoa(v.numValue)
	case valFunc:
		return strings.ToUpper(v.funcName)
	case valBool:
		if v.boolValue {
			return "T"
		}
		return "NIL"
	default:
		return "<UNK>"
	}
}
