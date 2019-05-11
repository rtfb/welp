package welp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rtfb/welp/lexer"
	"github.com/rtfb/welp/object"
	"github.com/rtfb/welp/parser"
)

// Environ represents the execution environment.
type Environ struct {
	vars map[string]object.Object
}

// NewEnv creates an environment.
func NewEnv() *Environ {
	return &Environ{
		vars: make(map[string]object.Object),
	}
}

func copyEnv(src *Environ) *Environ {
	dst := NewEnv()
	for k, v := range src.vars {
		dst.vars[k] = v
	}
	return dst
}

func num(env *Environ, tok lexer.Token) int64 {
	switch tok.Typ {
	case lexer.TokNumber:
		n, err := strconv.Atoi(string(tok.Value))
		if err != nil {
			panic(err)
		}
		return int64(n)
	case lexer.TokIdentifier:
		val := env.vars[string(tok.Value)]
		intVal, ok := val.(*object.Integer)
		if !ok {
			fmt.Printf("wrong type %T, expected %s", val, object.IntegerType)
			return -1
		}
		return intVal.Value
	default:
		fmt.Printf("error looking up num for %s\n", tok.String())
		panic("err")
	}
}

// TODO fix parser to produce an actual nil instead of this node
func nilNode(n *parser.Node) bool {
	return n.Tok.Typ == lexer.TokVoid && n.L == nil && n.R == nil
}

// (add 3 7) => 10
func callUserFunc(env *Environ, f *callable, expr *parser.Node) object.Object {
	newFrame := copyEnv(env)
	param := f.params
	arg := expr
	// TODO: add checking. At least check if number of args is correct
	for param.R != nil && arg.R != nil {
		val := eval(env, arg)
		nval, ok := val.(*object.Integer)
		if !ok {
			fmt.Printf("Type error: unexpected type %T for %s\n", val, f.name)
		}
		newFrame.vars[string(param.L.Tok.Value)] = &object.Integer{
			Value: nval.Value,
		}
		param = param.R
		arg = arg.R
	}
	return eval(newFrame, f.body)
}

var indent int

func dump(expr *parser.Node) {
	indent += 4
	prefix := strings.Repeat(" ", indent)
	if expr.Tok.Typ == lexer.TokVoid {
		fmt.Printf("%s-\n", prefix)
	} else {
		fmt.Printf("%stok: %s\n", prefix, expr.Tok.String())
	}
	if expr.L != nil {
		dump(expr.L)
	}
	if expr.R != nil {
		dump(expr.R)
	}
	indent -= 4
}
