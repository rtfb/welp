package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArray(t *testing.T) {
	a := Array{}
	assert.Equal(t, "[]", a.Inspect())
	a.Value = []Object{&Integer{Value: 7}}
	assert.Equal(t, "[7]", a.Inspect())
	a.Value = []Object{&Integer{Value: 7}, &Integer{Value: 12}}
	assert.Equal(t, "[7, 12]", a.Inspect())
}
