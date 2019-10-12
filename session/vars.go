package session

import (
	"fmt"
	"regexp"
	"strings"
)

type Var struct {
	Name  string
	Type  string
	Value string
}
type Vars map[string]Var

func (v Var) String() string {
	if v.Value != "" {
		return fmt.Sprintf("%s %s = %s", v.Name, v.Type, v.Value)
	}
	return fmt.Sprintf("%s %s", v.Name, v.Type)
}

func (s *Session) addVar(v Var) {
	s.vars[strings.TrimSpace(v.Name)] = v
}

func (vs Vars) String() string {
	var sets []string
	for _, v := range vs {
		sets = append(sets, v.String())
	}
	return strings.Join(sets, "\n\t")
}

func NewVar(code string) Var {
	if strings.Contains(code, "=") {
		regex := regexp.MustCompile(`(var)?\s*(?P<varnames>[a-zA-Z0-9_,\s]+)\s*(?P<type>[a-zA-Z0-9_]+)?\s*:?=(?P<value>.+)`)
		matched, err := reSubMatchMap(regex, code)
		if err != nil {
			return Var{}
		}
		varname, _ := matched["varnames"]
		value, _ := matched["value"]
		typ, _ := matched["type"]
		return Var{
			varname, typ, value,
		}
	}
	regex := regexp.MustCompile(`(var)?\s*(?P<varname>[a-zA-Z0-9_]+)\s*(?P<type>.+)`)
	matched, err := reSubMatchMap(regex, code)
	if err != nil {
		return Var{}
	}
	varname, _ := matched["varname"]
	typ, _ := matched["type"]
	return Var{varname, typ, ""}
}

func IsVarDecl(code string) bool {
	return (strings.Contains(code, "=") && !strings.Contains(code, "==") && !strings.Contains(strings.Split(code, "=")[0], ".")) ||
		strings.Contains(code, "var")
}
