package welp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rtfb/welp/lexer"
	"github.com/rtfb/welp/parser"
)

type Environ struct {
	vars map[string]*value
}

// NewEnv creates an environment.
func NewEnv() *Environ {
	return &Environ{
		vars: make(map[string]*value),
	}
}

func copyEnv(src *Environ) *Environ {
	dst := NewEnv()
	for k, v := range src.vars {
		dst.vars[k] = v
	}
	return dst
}

func num(env *Environ, tok lexer.Token) int {
	switch tok.Typ {
	case lexer.TokNumber:
		n, err := strconv.Atoi(string(tok.Value))
		if err != nil {
			panic(err)
		}
		return n
	case lexer.TokIdentifier:
		val := env.vars[string(tok.Value)]
		if val.typ != valNum {
			fmt.Printf("wrong type %s, expected %s", val.typ, valNum)
			return -1
		}
		return val.numValue
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
func callUserFunc(env *Environ, f *callable, expr *parser.Node) *value {
	newFrame := copyEnv(env)
	param := f.params
	arg := expr
	// TODO: add checking. At least check if number of args is correct
	for param.R != nil && arg.R != nil {
		nval := eval(env, arg)
		if nval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for %s\n", nval.typ.String(), f.name)
		}
		newFrame.vars[string(param.L.Tok.Value)] = &value{
			typ:      valNum,
			numValue: nval.numValue,
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

/*
func main() {
	env := NewEnv()
	expr := parser.ParseString(`
(fn fib (n)
  (cond
    ((eq n 1) 1)
    ((eq n 2) 1)
    (t (+ (fib (- n 1)) (fib (- n 2))))))`)
	println(eval(env, expr).String())
	expr = parser.ParseString("(fib 7)") // => 13
	println(eval(env, expr).String())
}
*/
