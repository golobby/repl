package session

import (
	"errors"
	"regexp"
)

type Type uint8

const (
	Import = Type(iota)
	TypeDecl
	FuncDecl
	VarDecl
	Shell
	Print
	Unknown
	Empty
)

func Parse(code string) (Type, error) {
	if isEmpty(code) {
		return Empty, nil
	} else if isShellCommand(code) {
		return Shell, nil
	} else if isImport(code) {
		return Import, nil
	} else if IsFuncDecl(code) {
		return FuncDecl, nil
	} else if isTypeDecl(code) {
		return TypeDecl, nil
	} else if isPrint(code) {
		return Print, nil
	} else if isVarDecl(code) {
		return VarDecl, nil
	} else {
		return Unknown, nil
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
