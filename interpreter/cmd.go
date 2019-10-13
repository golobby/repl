package interpreter

import "strings"

type REPLCmd uint8

const (
	REPLCmdDoc = REPLCmd(iota)
	REPLCmdHelp
	REPLCmdTypeVal
	REPLCmdPop
	REPLCmdDump
	REPLCmdFile
	REPLCmdVars
	REPLCmdImports
	REPLCmdTypes
	REPLCmdFuncs
)

func isShellCommand(code string) bool {
	if len(code) == 0 {
		return false
	}
	return code[0] == ':'
}

func (s *Interpreter) handleShellCommands(code string) error {
	typ, data := ParseCmd(code)
	switch typ {
	case REPLCmdDoc:
		output, err := goDoc(data)
		if err != nil {
			return err
		}
		s.shellCmdOutput = string(output)
		return nil
	case REPLCmdHelp:
		s.shellCmdOutput = helpText
		return nil
	case REPLCmdTypeVal:
		return s.Add(wrapInPrint(data))
	case REPLCmdPop:
		s.code = s.code[:len(s.code)-1]
		return nil
	case REPLCmdDump:
		s.shellCmdOutput = s.dump()
		return nil
	case REPLCmdFile:
		s.shellCmdOutput = s.String()
		return nil
	case REPLCmdImports:
		s.shellCmdOutput = s.imports.AsDump()
		return nil
	case REPLCmdVars:
		s.shellCmdOutput = s.vars.String()
		return nil
	case REPLCmdFuncs:
		s.shellCmdOutput = s.funcsAsString()
		return nil
	case REPLCmdTypes:
		s.shellCmdOutput = s.typesForSource()
		return nil
	default:
		return nil
	}
	return nil
}

func ParseCmd(code string) (REPLCmd, string) {
	if isGoDoc(code) {
		return REPLCmdDoc, strings.Split(code, " ")[1]
	} else if isHelp(code) {
		return REPLCmdHelp, ""
	} else if isTypeVal(code) {
		return REPLCmdTypeVal, code[strings.Index(code, " ")+1:]
	} else if isPop(code) {
		if len(strings.Split(code, " ")) > 1 {
			return REPLCmdPop, strings.Split(code, " ")[1]
		}
		return REPLCmdPop, ""
	} else if isFile(code) {
		return REPLCmdFile, ""
	} else if isDump(code) {
		return REPLCmdDump, ""
	} else if isTypes(code) {
		return REPLCmdTypes, ""
	} else if isVars(code) {
		return REPLCmdVars, ""
	} else if isImports(code) {
		return REPLCmdImports, ""
	} else if isFuncs(code) {
		return REPLCmdFuncs, ""
	}
	return 0, ""
}

func isPop(code string) bool {
	if len(code) < len(":pop") {
		return false
	}
	return code[:4] == ":pop"
}
func isDump(code string) bool {
	if len(code) < len(":dump") {
		return false
	}
	return code[:5] == ":dump"
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

func isFile(code string) bool {
	if len(code) < len(":file") {
		return false
	}
	return code[:5] == ":file"
}

func isTypes(code string) bool {
	if len(code) < len(":types") {
		return false
	}
	return code[:6] == ":types"
}
func isVars(code string) bool {
	if len(code) < len(":vars") {
		return false
	}
	return code[:5] == ":vars"
}
func isFuncs(code string) bool {
	if len(code) < len(":funcs") {
		return false
	}
	return code[:6] == ":funcs"
}
func isImports(code string) bool {
	if len(code) < len(":imports") {
		return false
	}
	return code[:7] == ":imports"
}
