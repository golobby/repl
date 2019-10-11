package session

import (
	"fmt"
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

func (s *Session) typesAsString() string {
	var types []string
	for k, v := range s.types {
		types = append(types, fmt.Sprintf("%s => %s", k, v))
	}
	return strings.Join(types, "\n")
}

func (s *Session) typesForSource() string {
	var ts []string
	for _, v := range s.types {
		ts = append(ts, v)
	}
	return strings.Join(ts, "\n")
}

func ExtractTypeName(code string) string {
	regex := regexp.MustCompile(`type\s+(?P<name>[a-zA-Z0-9_]+)(.|\s)+`)
	matched, err := reSubMatchMap(regex, code)
	if err != nil {
		return ""
	}
	typeName, _ := matched["name"]
	return typeName
}
