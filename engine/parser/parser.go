package parser

import (
	"regexp"
	"strings"
)

type StmtType uint8

const (
	StmtTypeFunctionCall = iota
	StmtTypeImport
	StmtTypePrint
	StmtTypeComment
	StmtTypeExpr
	StmtTypeTypeDecl
	StmtTypeVarDecl
	StmtTypeFuncDecl
	StmtUnknown
	StmtEmpty
)

func Parse(code string) (StmtType, bool, error) {

	if len(code) < 1 {
		return StmtEmpty, false, nil
	}
	if isComment(code) {
		return StmtTypeComment, ShouldContinue(code), nil
	} else if isImport(code) {
		return StmtTypeImport, ShouldContinue(code), nil
	} else if isFunc(code) {
		return StmtTypeFuncDecl, ShouldContinue(code), nil
	} else if isTypeDecl(code) {
		return StmtTypeTypeDecl, ShouldContinue(code), nil
	} else if isPrint(code) {
		return StmtTypePrint, ShouldContinue(code), nil
	} else if isComment(code) {
		return StmtTypeComment, ShouldContinue(code), nil
	} else if isExpr(code) {
		return StmtTypeExpr, ShouldContinue(code), nil
	} else {
		return StmtUnknown, ShouldContinue(code), nil
	}
}
func ShouldContinue(code string) bool {
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
	return stillOpenChars > 0
}
func isComment(code string) bool {
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
func isFunctionCall(code string) bool {
	m, err := regexp.Match("^[a-zA-Z0-9_.-]+\\(.*\\)", []byte(code))
	if err != nil {
		return false
	}
	return m && strings.ContainsAny(code, "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm")
}

func isExpr(code string) bool {
	if (strings.Contains(code, "=") && !strings.Contains(code, "==")) || strings.Contains(code, "var") || isFunctionCall(code) {
		return false
	}
	return true
}
func isFunc(code string) bool {
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
	matched1, err := regexp.Match("^fmt.Print.*\\(.*\\)", []byte(code))
	if err != nil {
		panic(err)
	}
	matched2, err := regexp.Match("^print(ln|f).*", []byte(code))
	if err != nil {
		panic(err)
	}
	return matched1 || matched2
}
