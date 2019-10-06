package main

import (
	"os"
	"testing"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func Test_newSession(t *testing.T) {
	monkey.Patch(createTmpDir, func(wd string) (string, error) {
		return "somedir/session", nil
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
	s := &session{}
	s.addCode("fmt.Println(12)")
	s.addCode("fmt.Println(13)")
	assert.Equal(t, []string{"fmt.Println(12)", "fmt.Println(13)"}, s.code)
}
func Test_addImport(t *testing.T) {
	s := &session{}
	s.addImport(`import "fmt"`)
	s.addImport(`import "os/exec"`)
	assert.Equal(t, []string{`import "fmt"`, `import "os/exec"`}, s.imports)
}
func Test_removeLastCode(t *testing.T) {
	s := &session{}
	s.code = append(s.code, "some ok code", "some code caused error")
	s.removeLastCode()
	assert.Equal(t, []string{"some ok code"}, s.code)
}
func Test_removeTmpCodes(t *testing.T) {
	s := &session{}
	s.code = append(s.code, `a := 1+2`)
	s.code = append(s.code, `fmt.Println("aaa")`)
	s.tmpCodes = append(s.tmpCodes, 1)
	s.removeTmpCodes()
	assert.Equal(t, []string{"a := 1+2", ""}, s.code)
}
func Test_validGoFileFromSession(t *testing.T) {
	s := &session{}
	s.addImport(`import "fmt"`)
	s.addCode(`fmt.Println("hey")`)
	s.addCode(`var a int`)
	assert.Equal(t, `package main
import "fmt"

func main() {
fmt.Println("hey")
var a int
}`, s.validGoFileFromSession())
}

func Test_shouldContinue(t *testing.T) {
	s := &session{}
	code1 := "fmt.Println(\n"
	assert.True(t, s.shouldContinue(code1))
	assert.Equal(t, 1, s.indents)
	s = &session{}
	code2 := "fmt.Println(fmt.Sprint(2"
	assert.True(t, s.shouldContinue(code2))
	assert.Equal(t, 2, s.indents)
	s = &session{}
	code3 := "{fmt.Print("
	assert.True(t, s.shouldContinue(code3))
	assert.Equal(t, 2, s.indents)
}
