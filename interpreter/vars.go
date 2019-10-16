package interpreter

import (
	"fmt"
	"go/token"
	"strings"
)

type Var struct {
	Name  string
	Type  string
	Value string
}
type Vars map[string]Var

func (v Var) String() string {
	if v.Value != "" {
		return fmt.Sprintf("%s %s = %s", v.Name, v.Type, v.Value)
	}
	return fmt.Sprintf("%s %s", v.Name, v.Type)
}

func (s *Interpreter) addVar(v Var) {
	if strings.Contains(v.Name, ",") {
		// multiple variable definition
		splitted := strings.Split(v.Name, ",")
		for _, sv := range splitted {
			_, exists := s.vars[sv]
			if exists {
				delete(s.vars, sv)
			}
		}
	} else {
		for varName, value := range s.vars {
			if strings.Contains(varName, ",") {
				spl := strings.Split(varName, ",")
				for idx := range spl {
					if spl[idx] == v.Name {
						spl[idx] = "_"
					}
				}
				delete(s.vars, varName)
				s.vars[strings.Join(spl, ",")] = value
			}
		}
	}

	s.vars[strings.TrimSpace(v.Name)] = v
}

func (vs Vars) String() string {
	var sets []string
	for _, v := range vs {
		sets = append(sets, v.String())
	}
	return strings.Join(sets, "\n\t")
}

func isVarDecl(code string) bool {
	tokens, _ := tokenizerAndLiterizer(code)
	for _, t := range tokens {
		if t == token.DEFINE || t == token.VAR {
			return true
		}
	}
	return false
}

func NewVar(code string) Var {
	tokens, lits := tokenizerAndLiterizer(code)
	for _, t := range tokens {
		if t == token.DEFINE {
			return extractDataFromVarWithDefine(tokens, lits)
		} else if t == token.VAR {
			return extractDataFromVarWithVar(tokens, lits)
		}
	}
	return Var{}
}

func extractDataFromVarWithVar(tokens []token.Token, lits []string) Var {
	var names []string
	var values []string
	for idx, tok := range tokens {
		if tok == token.ASSIGN {
			for i := 0; i < idx; i++ {
				if tokens[i] == token.VAR {
					continue
				}
				names = append(names, lits[i])
			}
			for i := idx + 1; i < len(tokens)-1; i++ {
				if lits[i] == "" {
					values = append(values, tokens[i].String())
				} else {
					values = append(values, lits[i])
				}
			}
		}
	}

	return Var{Name: strings.Join(names, ""), Value: strings.Join(values, "")}
}
func extractDataFromVarWithDefine(tokens []token.Token, lits []string) Var {
	var names []string
	var values []string
	for idx, tok := range tokens {
		if tok == token.DEFINE {
			for i := 0; i < idx; i++ {
				if tokens[i] == token.IDENT {
					names = append(names, lits[i])
				}
			}
			for i := idx + 1; i < len(tokens)-1; i++ {
				if tokens[i] == token.IDENT || tokens[i] == token.STRING {
					values = append(values, lits[i])
					continue
				}
				values = append(values, tokens[i].String())
			}
			break
		}
	}
	return Var{Name: strings.Join(names, ", "), Value: strings.Join(values, "")}
}
