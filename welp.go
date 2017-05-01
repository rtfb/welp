package welp

import (
	"fmt"
	"strconv"
	"strings"
)

type environ struct {
	vars map[string]*value
}

// NewEnv creates an environment.
func NewEnv() *environ {
	return &environ{
		vars: make(map[string]*value),
	}
}

func copyEnv(src *environ) *environ {
	dst := NewEnv()
	for k, v := range src.vars {
		dst.vars[k] = v
	}
	return dst
}

func num(env *environ, tok token) int {
	switch tok.typ {
	case tokNumber:
		n, err := strconv.Atoi(string(tok.value))
		if err != nil {
			panic(err)
		}
		return n
	case tokIdentifier:
		val := env.vars[string(tok.value)]
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
func nilNode(n *Node) bool {
	return n.Tok.typ == tokVoid && n.L == nil && n.R == nil
}

// (add 3 7) => 10
func callUserFunc(env *environ, f *callable, expr *Node) *value {
	newFrame := copyEnv(env)
	param := f.params
	arg := expr
	// TODO: add checking. At least check if number of args is correct
	for param.R != nil && arg.R != nil {
		nval := eval(env, arg)
		if nval.typ != valNum {
			fmt.Printf("Type error: unexpected type %s for %s\n", nval.typ.String(), f.name)
		}
		newFrame.vars[string(param.L.Tok.value)] = &value{
			typ:      valNum,
			numValue: nval.numValue,
		}
		param = param.R
		arg = arg.R
	}
	return eval(newFrame, f.body)
}

var indent int

func dump(expr *Node) {
	indent += 4
	prefix := strings.Repeat(" ", indent)
	if expr.Tok.typ == tokVoid {
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

func main() {
	env := NewEnv()
	expr := ParseString(`
(fn fib (n)
  (cond
    ((eq n 1) 1)
    ((eq n 2) 1)
    (t (+ (fib (- n 1)) (fib (- n 2))))))`)
	println(eval(env, expr).String())
	expr = ParseString("(fib 7)") // => 13
	println(eval(env, expr).String())
}
