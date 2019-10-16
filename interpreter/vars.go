package interpreter

import (
	"fmt"
	"go/token"
	"strings"
)

type Var struct {
	Name  string
	Value string
}
type Vars map[string]Var

func (v Var) String() string {
	if v.Value != "" {
		return fmt.Sprintf("%s = %s", v.Name, v.Value)
	}
	return fmt.Sprintf("%s", v.Name)
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
				varNameSplitted := strings.Split(varName, ",")
				for idx := range varNameSplitted {
					if varNameSplitted[idx] == v.Name {
						varNameSplitted[idx] = "_"
					}
				}
				delete(s.vars, varName)
				s.vars[strings.Join(varNameSplitted, ",")] = value
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
			return extractDataFromVarWithVar(code)
		}
	}
	return Var{}
}

func extractDataFromVarWithVar(code string) Var {
	var names string
	var values string
	if strings.Contains(code, "=") {
		names = strings.Split(code, "=")[0]
		names = names[3:]
		names = strings.TrimSpace(names)
		values = strings.Split(code, "=")[1]
		values = strings.TrimSpace(values)
		return Var{Name: names, Value: values}
	} else {
		return Var{Name: code[strings.Index(code, "var")+len("var"):]}
	}

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
				if tokens[i] == token.IDENT || tokens[i] == token.STRING || tokens[i] == token.INT || tokens[i] == token.FLOAT {
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
