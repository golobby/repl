package session

import (
	"fmt"
	"go/token"
	"strings"
)

type ImportData struct {
	Path  string
	Alias string
}
type ImportDatas []ImportData

func (is ImportDatas) String() string {
	var imports []string
	for _, i := range is {
		imports = append(imports, i.String())
	}
	return fmt.Sprintf("import (\n%s\n)", strings.Join(imports, "\n"))
}
func (is ImportDatas) AsDump() string {
	var imports []string
	for _, i := range is {
		imports = append(imports, i.String())
	}
	return fmt.Sprintf("%s", strings.Join(imports, "\n"))
}

func (i ImportData) String() string {
	return fmt.Sprintf("%s %s", i.Alias, i.Path)
}

func isImport(im string) bool {
	tokens, _ := tokenizerAndLiterizer(im)
	for _, tok := range tokens {
		if tok == token.EOF {
			break
		}
		if tok == token.IMPORT {
			return true
		}
	}
	return false
}

func ExtractImportData(im string) []ImportData {
	tokens, lits := tokenizerAndLiterizer(im)
	var imports []ImportData
	var currentImport ImportData
	for idx, tok := range tokens {
		if tok == token.EOF {
			break
		}
		if tok == token.IMPORT || tok == token.LPAREN || tok == token.RPAREN {
			continue
		}
		if lits[idx] == "\n" {
			if currentImport.Path != "" {
				imports = append(imports, currentImport)
			}
			currentImport = ImportData{}
		}
		if tok == token.IDENT {
			currentImport.Alias = lits[idx]
		}
		if tok == token.STRING {
			currentImport.Path = lits[idx]
		}
	}
	return imports
}
