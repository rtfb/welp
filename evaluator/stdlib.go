package evaluator

import (
	"os"

	"github.com/rtfb/welp/object"
	"github.com/rtfb/welp/parser"
)

var stdlibEnv *Environ

func init() {
	stdlibEnv = initStdlib()
}

func initStdlib() *Environ {
	bootstrapEnv := &Environ{
		vars:  make(map[string]object.Object),
		funcs: makeBuiltins(),
	}
	err := EvalFile(bootstrapEnv, "stdlib/stdlib.lisp")
	if err != nil {
		panic(err)
	}
	return bootstrapEnv
}

// EvalFile reads a file and evaluates its entire content.
func EvalFile(env *Environ, name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	ch := parser.ParseStream(f)
	for expr := range ch {
		Eval(env, expr)
	}
	return nil
}
