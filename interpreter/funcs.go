package interpreter

import (
	"fmt"
	"go/token"
	"strings"
)

func (s *Interpreter) addFunc(name string, value string) {
	s.funcs[name] = value
}

func IsFuncDecl(code string) bool {
	tokens, _ := tokenizerAndLiterizer(code)

	for idx, tok := range tokens {
		if idx == 0 && tok == token.FUNC {
			return true
		}
	}
	return false
}

func ExtractFuncName(code string) string {
	tokens, lits := tokenizerAndLiterizer(code)
	var inParen int
	for idx, tok := range tokens {
		if tok == token.LPAREN {
			inParen++
			continue
		}
		if tok == token.RPAREN {
			inParen--
			continue
		}
		if inParen == 0 && tok == token.IDENT {
			return lits[idx]
		}
	}
	return ""
}

func (s *Interpreter) funcsForSource() string {
	var fs []string
	for _, v := range s.funcs {
		fs = append(fs, v)
	}
	return strings.Join(fs, "\n")
}

func (s *Interpreter) funcsAsString() string {
	var funcs []string
	for k, v := range s.funcs {
		funcs = append(funcs, fmt.Sprintf("%s => %s", k, v))
	}
	return strings.Join(funcs, "\n")
}
