package interpreter

import (
	"go/token"
	"strings"
)

func (s *Interpreter) addType(name string, code string) {
	s.types[name] = code
}

func isTypeDecl(code string) bool {
	tokens, _ := tokenizerAndLiterizer(code)
	for _, t := range tokens {
		if t == token.TYPE {
			return true
		}
	}
	return false
}

func (s *Interpreter) typesForSource() string {
	var ts []string
	for _, v := range s.types {
		ts = append(ts, v)
	}
	return strings.Join(ts, "\n\t")
}

func ExtractTypeName(code string) string {
	tokens, lits := tokenizerAndLiterizer(code)
	for i, t := range tokens {
		if t == token.IDENT {
			return lits[i]
		}
	}
	return ""
}
