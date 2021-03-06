package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/rtfb/welp/evaluator"
	"github.com/rtfb/welp/parser"
)

type repl struct {
	rl     *readline.Instance
	ch     chan *parser.Node
	prompt string
	env    *evaluator.Environ
	r      *io.PipeReader
	w      *io.PipeWriter
	p      *parser.Parser
}

func newREPL() (*repl, error) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 "welp> ",
		HistoryFile:            "/tmp/readline-multiline",
		DisableAutoSaveHistory: true,
	})
	if err != nil {
		return nil, err
	}
	r, w := io.Pipe()
	return &repl{
		rl:     rl,
		ch:     make(chan *parser.Node),
		prompt: "welp> ",
		env:    evaluator.New().NewEnv(),
		r:      r,
		w:      w,
		p:      parser.New(r),
	}, nil
}

func (r *repl) epl() {
	expr, _ := r.p.Parse()
	if expr == nil {
		r.rl.SetPrompt("> ")
		return
	}
	r.rl.SetPrompt("welp> ")
	fmt.Println(evaluator.Eval(r.env, expr).Inspect())
	r.p.Reset()
}

func (r *repl) Run() {
	for {
		line, err := r.rl.Readline()
		if err != nil || line == "(q)" {
			break
		}
		r.w.Write([]byte(line + "\n"))
		r.epl()
	}
	fmt.Println("Quitting")
}

func main() {
	if len(os.Args) > 1 {
		env := evaluator.New().NewEnv()
		if err := evaluator.EvalFile(env, os.Args[1]); err != nil {
			panic(err)
		}
		return
	}
	repl, err := newREPL()
	if err != nil {
		panic(err)
	}
	defer repl.rl.Close()
	defer repl.r.Close()
	defer repl.w.Close()
	repl.p.Start()
	repl.Run()
}
