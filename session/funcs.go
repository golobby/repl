package session

import (
	"fmt"
	"strings"
)

func (s *Session) addFunc(name string, value string) {
	s.funcs[name] = value
}

func (s *Session) funcsForSource() string {
	var fs []string
	for _, v := range s.funcs {
		fs = append(fs, v)
	}
	return strings.Join(fs, "\n")
}

func (s *Session) funcsAsString() string {
	var funcs []string
	for k, v := range s.funcs {
		funcs = append(funcs, fmt.Sprintf("%s => %s", k, v))
	}
	return strings.Join(funcs, "\n")
}
