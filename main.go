package main

import (
	"encoding/base64"
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
	version = "0.1a"
	logo    = "CiAgIF9fX19fXyAgICAgIF9fICAgICAgICAgIF9fICAgIF9fICAgICAgICAgICAgIF9fX18gIF9fX19fX19fX18gIF9fIAogIC8gX19f" +
		"Xy9fX18gIC8gLyAgIF9fX18gIC8gL18gIC8gL18gIF9fICBfXyAgIC8gX18gXC8gX19fXy8gX18gXC8gLyAKIC8gLyBfXy8gX18gXC8gLy" +
		"AgIC8gX18gXC8gX18gXC8gX18gXC8gLyAvIC8gIC8gL18vIC8gX18vIC8gL18vIC8gLyAgCi8gL18vIC8gL18vIC8gL19fXy8gL18vIC8g" +
		"L18vIC8gL18vIC8gL18vIC8gIC8gXywgXy8gL19fXy8gX19fXy8gL19fXwpcX19fXy9cX19fXy9fX19fXy9cX19fXy9fLl9fXy9fLl9fXy9" +
		"cX18sIC8gIC9fLyB8Xy9fX19fXy9fLyAgIC9fX19fXy8KICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAvX19fXy8gICAgIC" +
		"AgICAgICAgICAgICAgICAgICAgICAg"
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
	err := currentInterpreter.Add(input)
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	fmt.Print(currentInterpreter.Eval())
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

	l, _ := base64.StdEncoding.DecodeString(logo)
	fmt.Println(string(l))
	fmt.Printf("GoLobby REPL v%s\n", version)

	p := prompt.New(handler, completer, prompt.OptionPrefix(">>> "))
	p.Run()
}
