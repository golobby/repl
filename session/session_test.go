package session

import (
	"io/ioutil"
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func Test_newSession(t *testing.T) {
	monkey.Patch(createTmpDir, func(wd string) (string, error) {
		return "somedir/Session", nil
	})
	monkey.Patch(os.Chdir, func(string) error {
		return nil
	})
	monkey.Patch(getModuleNameOfCurrentProject, func(string) string {
		return "tmpmodule"
	})
	monkey.Unpatch(createTmpDir)
	monkey.Unpatch(os.Chdir)
	monkey.Unpatch(getModuleNameOfCurrentProject)
}

func Test_addCode(t *testing.T) {
	s := &Session{}
	//err := s.Add("fmt.Println(12)")
	//assert.NoError(t, err)
	//assert.Equal(t, []string{"fmt.Println(12)"}, s.code)
	//err = s.Add("fmt.Println(13)")
	//assert.NoError(t, err)
	//assert.Equal(t, []string{"fmt.Println(13)"}, s.code)
	//err = s.Add("type user struct{")
	//assert.NoError(t, err)
	//err = s.Add("Name string")
	//assert.NoError(t, err)
	//err = s.Add("}")
	//assert.NoError(t, err)
	//assert.Equal(t, []string{"type user struct{\nName string\n}"}, s.typesAndMethods)
	//err = s.Add("")
	//assert.NoError(t, err)
	//assert.Equal(t, []string{}, s.code)
	s.code = []string{}
	err := s.Add("fmt.Println(")
	assert.NoError(t, err)
	err = s.Add(`"Salam"`)
	assert.NoError(t, err)
	err = s.Add(`)`)
	assert.NoError(t, err)
	assert.Equal(t, []string{"fmt.Println(\n\"Salam\"\n)"}, s.code)
}
func Test_addImport(t *testing.T) {
	s := &Session{}
	s.addImport(`import "fmt"`)
	s.addImport(`import "os/exec"`)
	assert.Equal(t, []string{`import "fmt"`, `import "os/exec"`}, s.imports)
}
func Test_removeLastCode(t *testing.T) {
	s := &Session{}
	s.code = append(s.code, "some ok code", "some code caused error")
	s.removeLastCode()
	assert.Equal(t, []string{"some ok code"}, s.code)
}
func Test_removeTmpCodes(t *testing.T) {
	s := &Session{}
	s.code = append(s.code, `a := 1+2`)
	s.code = append(s.code, `fmt.Println("aaa")`)
	s.tmpCodes = append(s.tmpCodes, 1)
	s.removeTmpCodes()
	assert.Equal(t, []string{"a := 1+2"}, s.code)
}
func Test_validGoFileFromSession(t *testing.T) {
	s := &Session{}
	s.addImport(`import "fmt"`)
	s.Add(`fmt.Println("hey")`)
	s.Add(`var a int`)
	assert.Equal(t, "package main\nimport \"fmt\"\n\n\nvar (\n)\nfunc main() {\nvar a int\n}", s.String())
}

func Test_add_print(t *testing.T) {
	s := &Session{}
	s.Add(`fmt.Println("Salam")`)
	assert.Equal(t, []string{`fmt.Println("Salam")`}, s.code)
	assert.Equal(t, []int{0}, s.tmpCodes)
}

func Test_add_isImport(t *testing.T) {
	s := &Session{}
	s.Add(`import "github.com"`)
	assert.Equal(t, s.imports, []string{`import "github.com"`})
}
func Test_add_function_call(t *testing.T) {
	s := &Session{}
	s.Add(`someFunc("salam man be to yare ghadimi")`)
	assert.Equal(t, s.code, []string{`someFunc("salam man be to yare ghadimi")`})
}
func Test_add_continue_mode(t *testing.T) {
	s := &Session{}
	s.Add("fmt.Println(")
	s.Add("2")
	s.Add(")")
	assert.Equal(t, []string{"fmt.Println(\n2\n)"}, s.code)
}

func Test_checkIfErrIsNotDecl(t *testing.T) {
	assert.True(t, checkIfErrIsNotDecl(`"fmt" imported and not used`))
	assert.True(t, checkIfErrIsNotDecl(`a declared and not used`))
	assert.False(t, checkIfErrIsNotDecl("not able to compile"))
}
func Test_createReplaceRequireClause_with_moduleName(t *testing.T) {
	moduleName := "shell"
	localPath := "inja"
	assert.Equal(t, "replace shell => inja", createReplaceRequireClause(moduleName, localPath))
}

func Test_createReplaceRequireClause_without_moduleName(t *testing.T) {
	moduleName := ""
	localPath := "inja"
	assert.Equal(t, "", createReplaceRequireClause(moduleName, localPath))
}
func Test_wrapInPrint(t *testing.T) {
	assert.Equal(t, `fmt.Printf("<%T> %+v\n", 1+2, 1+2)`, wrapInPrint("1+2"))
	assert.Equal(t, `fmt.Printf("<%T> %+v\n", "Hello", "Hello")`, wrapInPrint(`"Hello"`))

}
func Test_multiplyString(t *testing.T) {
	assert.Equal(t, "", multiplyString("...", 0))
	assert.Equal(t, "...", multiplyString("...", 1))
	assert.Equal(t, "......", multiplyString("...", 2))
	assert.Equal(t, ".........", multiplyString("...", 3))

}

func Test_getModuleNameOfCurrentProject_in_go_project(t *testing.T) {
	monkey.Patch(ioutil.ReadFile, func(string) ([]byte, error) {
		return []byte(`module somemodule
go 1.13`), nil
	})
	moduleName := getModuleNameOfCurrentProject("somedir")
	assert.Equal(t, moduleName, "somemodule")
	monkey.Unpatch(ioutil.ReadFile)
}

func Test_getModuleNameOfCurrentProject_not_in_go_project(t *testing.T) {
	monkey.Patch(ioutil.ReadFile, func(string) ([]byte, error) {
		return nil, os.ErrNotExist
	})
	moduleName := getModuleNameOfCurrentProject("somedir")
	assert.Equal(t, moduleName, "")
	monkey.Unpatch(ioutil.ReadFile)
}
