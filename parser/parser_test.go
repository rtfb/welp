package parser

import (
	"testing"

	"github.com/rtfb/welp/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParseString(t *testing.T) {
	node := ParseString("(+ 1 2)")
	assert.NoError(t, node.Err)
	assert.NotNil(t, node)
	assert.Equal(t, lexer.TokIdentifier, node.L.Tok.Typ)
	assert.Equal(t, lexer.TokNumber, node.R.L.Tok.Typ)
}
