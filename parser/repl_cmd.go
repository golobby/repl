package parser

import "strings"

type REPLCmd uint8

const (
	REPLCmdDoc = REPLCmd(iota)
)

func ParseCmd(code string) (REPLCmd, string) {
	if isGoDoc(code) {
		return REPLCmdDoc, strings.Split(code, " ")[1]
	}
	return 0, ""
}

func isGoDoc(code string) bool {
	if len(code) == 0 {
		return false
	}
	return code[:4] == ":doc"
}
