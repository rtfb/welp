package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/rtfb/welp"
)

func doFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	env := welp.NewEnv()
	ch := welp.ParseStream(f)
	for e := range ch {
		if e.Err != nil {
			println("Error: ", e.Err.Error())
		} else {
			println(welp.Eval(env, e).String())
		}
	}
	return nil
}

func main() {
	if len(os.Args) > 1 {
		if err := doFile(os.Args[1]); err != nil {
			panic(err)
		}
		return
	}
	fmt.Println("welp")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		fmt.Print("> ")
		line := scanner.Text()
		if line == "(q)" {
			break
		}
	}
}
