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
	sess := &session{}
	sess.addCode("fmt.Println(12)")
	sess.addCode("fmt.Println(13)")
	assert.Equal(t, []string{"fmt.Println(12)", "fmt.Println(13)"}, sess.code)
}
func Test_addImport(t *testing.T) {
	sess := &session{}
	sess.addImport(`import "fmt"`)
	sess.addImport(`import "os/exec"`)
	assert.Equal(t, []string{`import "fmt"`, `import "os/exec"`}, sess.imports)
}
func Test_removeLastCode(t *testing.T) {
	sess := &session{}
	sess.code = append(sess.code, "some ok code", "some code caused error")
	sess.removeLastCode()
	assert.Equal(t, []string{"some ok code"}, sess.code)
}
func Test_removeTmpCodes(t *testing.T) {
	sess := &session{}
	sess.code = append(sess.code, `a := 1+2`)
	sess.code = append(sess.code, `fmt.Println("aaa")`)
	sess.tmpCodes = append(sess.tmpCodes, 1)
	sess.removeTmpCodes()
	assert.Equal(t, []string{"a := 1+2", ""}, sess.code)
}
func Test_validGoFileFromSession(t *testing.T) {
	sess := &session{}
	sess.addImport(`import "fmt"`)
	sess.addCode(`fmt.Println("hey")`)
	sess.addCode(`var a int`)
	assert.Equal(t, `package main
import "fmt"

func main() {
fmt.Println("hey")
var a int
}`, sess.validGoFileFromSession())
}
