package main

type token int

const (
	tokVoid = 0
	tokEOF  = 666
)

func getToken(input []byte, head int) (tok token, newHead int) {
	for head < len(input) && input[head] == ' ' {
		head++
	}
	if head < len(input) {
		return token(input[head]), head + 1
	}
	return tokEOF, head + 1
}
