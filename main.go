package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/golobby/repl/interpreter"

	"github.com/c-bata/go-prompt"
)

var (
	currentInterpreter *interpreter.Interpreter
	DEBUG              bool
)

const (
	version = "0.1.2"
)

func completer(d prompt.Document) []prompt.Suggest {
	var s []prompt.Suggest
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func handler(input string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Panic: %v\n%s", err, debug.Stack())
		}
	}()
	var start time.Time
	if DEBUG {
		start = time.Now()
	}
	out, err := currentInterpreter.Eval(input)
	if err != nil {
		fmt.Print(err.Error())
		return
	}
	fmt.Print(out)
	if DEBUG {
		fmt.Printf(":::::: D => %v\n", time.Since(start))
	}
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	debug := flag.Bool("timestamp", false, "turns timestamp mode on")
	flag.Parse()
	DEBUG = *debug

	currentInterpreter, err = interpreter.NewSession(wd)
	if err != nil {
		panic(err)
	}
	_, err = currentInterpreter.Eval(":e 1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("repl v%s\n", version)

	p := prompt.New(handler, completer, prompt.OptionPrefix("repl> "))
	p.Run()
}
