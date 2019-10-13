package main

import "github.com/golobby/repl/interpreter"

type Session struct {
	*interpreter.Interpreter
}

func newSession(wd string) *Session {
	i, err := interpreter.NewInterpreter(wd)
	if err != nil {
		panic(err)
	}
	return &Session{i}
}
