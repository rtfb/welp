package main

import (
	"fmt"
	"io"
	"os"

	"github.com/chzyer/readline"
	"github.com/rtfb/welp"
)

func doFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	env := welp.NewEnv()
	ch := welp.ParseStream(f)
	for expr := range ch {
		fmt.Println(welp.Eval(env, expr))
	}
	return nil
}

type repl struct {
	rl     *readline.Instance
	ch     chan *welp.Node
	prompt string
	env    *welp.Environ
	r      *io.PipeReader
	w      *io.PipeWriter
	p      *welp.Parser
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
		ch:     make(chan *welp.Node),
		prompt: "welp> ",
		env:    welp.NewEnv(),
		r:      r,
		w:      w,
		p:      welp.NewParser(r),
	}, nil
}

func (r *repl) epl() {
	expr, _ := r.p.Parse()
	if expr == nil {
		r.rl.SetPrompt("> ")
		return
	}
	r.rl.SetPrompt("welp> ")
	fmt.Println(welp.Eval(r.env, expr))
	r.p.Reset()
}

func (r *repl) Run() {
	for {
		line, err := r.rl.Readline()
		if err != nil || line == "(q)" {
			break
		}
		r.w.Write([]byte(line))
		r.epl()
	}
	fmt.Println("Quitting")
}

func main() {
	if len(os.Args) > 1 {
		if err := doFile(os.Args[1]); err != nil {
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
