package interpreter

import (
	"fmt"
	"go/token"
	"strings"
)

type ImportData struct {
	Path string
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
	return fmt.Sprintf("%s", i.Path)
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
	var imports []ImportData
	importLines := strings.Split(im, "\n")
	for _, i := range importLines {
		if strings.Contains(i, "import") {
			path := strings.TrimSpace(i[strings.Index(im, "import")+len("import"):])
			imports = append(imports, ImportData{Path: path})
		}
	}
	return imports
}
