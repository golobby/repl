package session

import (
	"errors"
	"regexp"
	"strings"
)

type StmtType uint8

const (
	StmtTypeImport = StmtType(iota)
	StmtTypePrint
	StmtTypeTypeDecl
	StmtTypeFuncDecl
	StmtUnknown
	StmtEmpty
	StmtShell
	StmtVarDecl
	StmtFunctionCall
)

func Parse(code string) (StmtType, error) {
	if isEmpty(code) {
		return StmtEmpty, nil
	} else if isShellCommand(code) {
		return StmtShell, nil
	} else if isImport(code) {
		return StmtTypeImport, nil
	} else if IsFunc(code) {
		return StmtTypeFuncDecl, nil
	} else if isTypeDecl(code) {
		return StmtTypeTypeDecl, nil
	} else if isPrint(code) {
		return StmtTypePrint, nil
	} else if isFunctionCall(code) {
		return StmtFunctionCall, nil
	} else if IsVarDecl(code) {
		return StmtVarDecl, nil
	} else {
		return StmtUnknown, nil
	}
}
func isFunctionCall(code string) bool {
	m, err := regexp.Match("^[a-zA-Z0-9_.-]+\\(.*\\)", []byte(code))
	if err != nil {
		return false
	}
	return m && strings.ContainsAny(code, "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm")

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

func reSubMatchMap(r *regexp.Regexp, str string) (map[string]string, error) {
	match := r.FindStringSubmatch(str)
	if len(match) == 0 {
		return nil, errors.New("cannot match")
	}
	subMatchMap := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			subMatchMap[name] = match[i]
		}
	}

	return subMatchMap, nil
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
