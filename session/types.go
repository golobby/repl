package session

import (
	"fmt"
	"strings"
)

func (s *Session) addType(name string, code string) {
	s.types[name] = code
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
