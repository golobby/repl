package session

import (
	"fmt"
	"regexp"
	"strings"
)

func (s *Session) addVar(name string, value string) {
	s.vars[name] = value
}

func (s *Session) varsForSource() string {
	var sets []string
	for k, v := range s.vars {
		sets = append(sets, fmt.Sprintf("%s = %s", k, v))
	}
	return "var (\n" + strings.Join(sets, "\n") + ")"
}

func (s *Session) varsString() string {
	var sets []string
	for k, v := range s.vars {
		sets = append(sets, fmt.Sprintf("%s => %s", k, v))
	}
	return strings.Join(sets, "\n")
}

func ExtractNameAndValueFromVarInit(code string) (string, string) {
	regex := regexp.MustCompile(`(var)?\s*(?P<varname>[a-zA-Z0-9_]+)\s*.*\s*:?=(?P<value>.+)`)
	matched, err := reSubMatchMap(regex, code)
	if err != nil {
		return "", ""
	}
	varname, _ := matched["varname"]
	value, _ := matched["value"]
	return varname, value
}

func IsVarDecl(code string) bool {
	regex := regexp.MustCompile(`(var)?\s*(?P<varname>[a-zA-Z0-9_]+)\s*.*\s*:?=(?P<value>.+)`)
	matched, err := reSubMatchMap(regex, code)
	if err != nil {
		return false
	}
	return !(len(matched) == 0)
}
