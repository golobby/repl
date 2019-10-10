package parser

import "strings"

type REPLCmd uint8

const (
	REPLCmdDoc = REPLCmd(iota)
	REPLCmdHelp
	REPLCmdTypeVal
	REPLCmdPop
	REPLCmdLog
)

func ParseCmd(code string) (REPLCmd, string) {
	if isGoDoc(code) {
		return REPLCmdDoc, strings.Split(code, " ")[1]
	} else if isHelp(code) {
		return REPLCmdHelp, ""
	} else if isTypeVal(code) {
		return REPLCmdTypeVal, strings.Split(code, " ")[1]
	} else if isPop(code) {
		if len(strings.Split(code, " ")) > 1 {
			return REPLCmdPop, strings.Split(code, " ")[1]
		}
		return REPLCmdPop, ""
	} else if isLog(code) {
		return REPLCmdLog, ""
	}
	return 0, ""
}

func isPop(code string) bool {
	if len(code) < len(":pop") {
		return false
	}
	return code[:4] == ":pop"
}
func isLog(code string) bool {
	if len(code) < len(":log") {
		return false
	}
	return code[:4] == ":log"
}
func isHelp(code string) bool {
	if len(code) < len(":help") {
		return false
	}
	seg := code[:5]
	return seg == ":help"
}
func isTypeVal(code string) bool {
	if len(code) < len(":e a") {
		return false
	}
	return code[:2] == ":e"
}

func isGoDoc(code string) bool {
	if len(code) < len(":doc a") {
		return false
	}
	seg := code[:4]
	return seg == ":doc"
}
