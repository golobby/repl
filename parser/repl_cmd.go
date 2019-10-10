package parser

import "strings"

type REPLCmd uint8

const (
	REPLCmdDoc = REPLCmd(iota)
	REPLCmdHelp
)

func ParseCmd(code string) (REPLCmd, string) {
	if isGoDoc(code) {
		return REPLCmdDoc, strings.Split(code, " ")[1]
	} else if isHelp(code) {
		return REPLCmdHelp, ""
	}
	return 0, ""
}

func isHelp(code string) bool {
	if len(code) == 0 {
		return false
	}
	seg := code[:5]
	return seg == ":help"
}

func isGoDoc(code string) bool {
	if len(code) == 0 {
		return false
	}
	seg := code[:4]
	return seg == ":doc"
}
