package parser

import "strings"

type REPLCmd uint8

const (
	REPLCmdDoc = REPLCmd(iota)
	REPLCmdHelp
	REPLCmdTypeVal
)

func ParseCmd(code string) (REPLCmd, string) {
	if isGoDoc(code) {
		return REPLCmdDoc, strings.Split(code, " ")[1]
	} else if isHelp(code) {
		return REPLCmdHelp, ""
	} else if isTypeVal(code) {
		return REPLCmdTypeVal, strings.Split(code, " ")[1]
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
func isTypeVal(code string) bool {
	if len(code) == 0 {
		return false
	}
	//:typeval
	return code[:8] == ":typeval"
}

func isGoDoc(code string) bool {
	if len(code) == 0 {
		return false
	}
	seg := code[:4]
	return seg == ":doc"
}
