package main

import "strconv"

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
)

type value struct {
	typ      valType
	numValue int
}

func (v *value) String() string {
	if v.typ == valNum {
		return strconv.Itoa(v.numValue)
	}
	return "<UNK>"
}
