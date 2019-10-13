package session

import (
	"go/token"
	"regexp"
	"strings"
)

func (s *Session) addType(name string, code string) {
	s.types[name] = code
}

func isTypeDecl(code string) bool {
	matched, err := regexp.Match("type .+", []byte(code))
	if err != nil {
		return false
	}
	return matched
}

func (s *Session) typesForSource() string {
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
