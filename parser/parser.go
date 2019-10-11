package parser

import (
	"regexp"
	"strings"
)

type StmtType uint8

const (
	StmtTypeFunctionCall = StmtType(iota)
	StmtTypeImport
	StmtTypePrint
	StmtTypeComment
	StmtTypeExpr
	StmtTypeTypeDecl
	StmtTypeFuncDecl
	StmtUnknown
	StmtEmpty
	StmtShell
	StmtVarDecl
)

func Parse(code string) (StmtType, error) {
	if isEmpty(code) {
		return StmtEmpty, nil
	} else if isShellCommand(code) {
		return StmtShell, nil
	} else if isComment(code) {
		return StmtTypeComment, nil
	} else if isImport(code) {
		return StmtTypeImport, nil
	} else if IsFunc(code) {
		return StmtTypeFuncDecl, nil
	} else if isTypeDecl(code) {
		return StmtTypeTypeDecl, nil
	} else if isPrint(code) {
		return StmtTypePrint, nil
	} else if IsVarDecl(code) {
		return StmtVarDecl, nil
	} else {
		return StmtUnknown, nil
	}
}
func ShouldContinue(code string) (int, bool) {
	var stillOpenChars int

	for _, c := range code {
		if c == '{' || c == '(' {
			stillOpenChars++
			continue
		}

		if c == '}' || c == ')' {
			stillOpenChars--
		}
	}

	return stillOpenChars, stillOpenChars > 0
}
func isEmpty(code string) bool {
	return len(code) == 0
}
func hasOutput(code string) bool {
	return strings.Contains(code, "=") && !strings.Contains(code, "==") && len(strings.Split(code, "=")) > 1
}
func isComment(code string) bool {
	if len(code) < 2 {
		return false
	}
	if code[:2] == "//" || code[:2] == "/*" {
		return true
	}
	return false
}

func isShellCommand(code string) bool {
	if len(code) == 0 {
		return false
	}
	return code[0] == ':'
}

func isTypeDecl(code string) bool {
	matched, err := regexp.Match("type .+", []byte(code))
	if err != nil {
		return false
	}
	return matched
}
func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
	match := r.FindStringSubmatch(str)
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap
}
func ExtractVarName(code string) string {
	regex := regexp.MustCompile(`(var)?\s+(?P<varname>[a-zA-Z0-9_]+)\s*.*(:?=.+)?`)
	matched := reSubMatchMap(regex, code)
	if name, ok := matched["varname"]; ok {
		return name
	}
	return ""
}
func ExtractFuncName(code string) string {
	matched := reSubMatchMap(regexp.MustCompile(`func\s+(\(.*\))(?P<funcname>.+)\(.*\).*`), code)
	if name, ok := matched["funcname"]; ok {
		return name
	}
	return ""
}
func IsVarDecl(code string) bool {
	matched1, err := regexp.Match(`(var)?\s*[a-zA-Z0-9]+\s+:?=\s*[a-zA-Z0-9_.-]+\(.*\)`, []byte(code))
	if err != nil {
		return false
	}
	matched2, err := regexp.Match(`var\s+[a-zA-Z0-9]+\s+.+`, []byte(code))
	if err != nil {
		return false
	}
	matched3, err := regexp.Match(`(var)?\s*[a-zA-Z0-9]+\s?:?=\s*.+`, []byte(code))
	if err != nil {
		return false
	}
	return matched1 || matched2 || matched3
}
func isFunctionCall(code string) bool {
	m, err := regexp.Match("^[a-zA-Z0-9_.-]+\\(.*\\)", []byte(code))
	if err != nil {
		return false
	}
	return m && strings.ContainsAny(code, "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm")
}

func isExpr(code string) bool {
	if IsVarDecl(code) || (isFunctionCall(code) && hasOutput(code)) {
		return false
	}
	return true
}
func IsFunc(code string) bool {
	matched, err := regexp.Match("^func.+", []byte(code))
	if err != nil {
		return false
	}
	return matched
}
func isImport(im string) bool {
	matched, err := regexp.Match("import .+", []byte(im))
	if err != nil {
		panic(err)
	}
	return matched
}
func isPrint(code string) bool {

	matched1, err := regexp.Match(`(fmt)\.Print.*\(\s*.*\s*\)`, []byte(code))
	if err != nil {
		panic(err)
	}
	matched2, err := regexp.Match(`^print(ln|f)\(\s*.*\s*\)`, []byte(code))
	if err != nil {
		panic(err)
	}
	return matched1 || matched2
}
