package evaluator

import "github.com/rtfb/welp/object"

// Environ represents the execution environment.
type Environ struct {
	vars  map[string]object.Object
	funcs map[string]*callable
}

// newEnv creates an environment.
func newEnv(evtor *Evaluator) *Environ {
	return newEmptyEnv().extend(evtor.stdlibEnv)
}

func newEmptyEnv() *Environ {
	return &Environ{
		vars:  make(map[string]object.Object),
		funcs: makeBuiltins(),
	}
}

func (e *Environ) deepCopy() *Environ {
	fresh := newEmptyEnv()
	return fresh.extend(e)
}

func (e *Environ) extend(src *Environ) *Environ {
	for k, v := range src.vars {
		e.vars[k] = v
	}
	for k, v := range src.funcs {
		e.funcs[k] = v
	}
	return e
}
