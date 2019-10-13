package interpreter

import (
	"fmt"
	"strings"
)

const dumpTEMPLATE = `imports => [
	%s
],
types => [
	%s
],
funcs => [
	%s
],
vars => [
	%s
],
main => [
	%s
]`

func StringOf(s map[string]string) string {
	var types []string
	for k, v := range s {
		types = append(types, fmt.Sprintf("%s => %s", k, v))
	}
	return strings.Join(types, "\n\t")
}

func (s *Interpreter) dump() string {
	return fmt.Sprintf(dumpTEMPLATE, s.imports.AsDump(), StringOf(s.types), s.funcsAsString(), s.vars.String(), strings.Join(s.code, "\n"))
}
