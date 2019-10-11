package session

import (
	"fmt"
	"strings"
)

func (s *Session) addVar(name string, value string) {
	s.vars[name] = value
}

func (s *Session) varsForSource() string {
	var sets []string
	for k, v := range s.vars {
		sets = append(sets, fmt.Sprintf("%s = %s", k, v))
	}
	return "var (\n" + strings.Join(sets, "\n") + ")"
}

func (s *Session) varsString() string {
	var sets []string
	for k, v := range s.vars {
		sets = append(sets, fmt.Sprintf("%s => %s", k, v))
	}
	return strings.Join(sets, "\n")
}
