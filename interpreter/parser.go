package interpreter

import (
	"go/scanner"
	"go/token"
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

func createScannerFor(code string) scanner.Scanner {
	var s scanner.Scanner
	fs := token.NewFileSet()
	s.Init(fs.AddFile("", fs.Base(), len(code)), []byte(code), nil, scanner.ScanComments)
	return s
}
func tokenizerAndLiterizer(code string) ([]token.Token, []string) {
	s := createScannerFor(code)
	tokens := []token.Token{}
	lits := []string{}
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		tokens = append(tokens, tok)
		lits = append(lits, lit)
	}
	return tokens, lits
}

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

func isPrint(code string) bool {
	tokens, lits := tokenizerAndLiterizer(code)
	for i, t := range tokens {
		if t == token.IDENT &&
			(lits[i] == "Println" || lits[i] == "Printf" || lits[i] == "Print" || lits[i] == "println") || lits[i] == "printf" || lits[i] == "print" {
			return true
		}
	}
	return false
}
