package session

import (
	"fmt"
	"regexp"
	"strings"
)

func (s *Session) addFunc(name string, value string) {
	s.funcs[name] = value
}

func IsFuncDecl(code string) bool {
	matched, err := regexp.Match("^func.+", []byte(code))
	if err != nil {
		return false
	}
	return matched
}

func ExtractFuncName(code string) string {
	matched, err := reSubMatchMap(regexp.MustCompile(`func\s+(\(.*\))?\s*(?P<funcname>[a-zA-Z0-9]+)\(.*\).*`), code)
	if err != nil {
		return ""
	}
	if name, ok := matched["funcname"]; ok {
		return name
	}
	return ""
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
