package interpreter

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	typ, err := Parse(``)
	assert.NoError(t, err)
	assert.Equal(t, Empty, typ)
	typ, err = Parse(`import "fmt"`)
	assert.NoError(t, err)
	assert.Equal(t, Import, typ)
	typ, err = Parse(`func name() string { return "" }`)
	assert.NoError(t, err)
	assert.Equal(t, FuncDecl, typ)
	typ, err = Parse(`type user struct{ Name string }`)
	assert.NoError(t, err)
	assert.Equal(t, TypeDecl, typ)
	typ, err = Parse(`fmt.Println("aleyk")`)
	assert.NoError(t, err)
	assert.Equal(t, Print, typ)
}
func Test_shouldContinue(t *testing.T) {

	code1 := "fmt.Println(\n"
	i, sc := ShouldContinue(code1)
	assert.True(t, sc)
	assert.Equal(t, 1, i)
	code2 := "fmt.Println(fmt.Sprint(2"
	i, sc = ShouldContinue(code2)
	assert.True(t, sc)
	assert.Equal(t, 2, i)
	code3 := "{fmt.Print("
	i, sc = ShouldContinue(code3)
	assert.True(t, sc)
	assert.Equal(t, 2, i)
	code4 := "fmt.Println(22)"
	i, sc = ShouldContinue(code4)
	assert.False(t, sc)
}
func Test_isImport_true(t *testing.T) {
	code := `import "github.com/golobby"`
	assert.True(t, isImport(code))
}
func Test_isImport_false(t *testing.T) {
	code := `impor "g"`
	assert.False(t, isImport(code))
}

func Test_isFunc_true(t *testing.T) {
	code := `func (User) Name() string{
		return u.name
}`
	assert.True(t, IsFuncDecl(code))
}

func Test_isFunc_false(t *testing.T) {
	code := `unc (User) Name string{
		return u.name
}`
	assert.False(t, IsFuncDecl(code))
}

func Test_isTypeDecl_true(t *testing.T) {
	code := `type user struct{}`
	assert.True(t, isTypeDecl(code))
}
func Test_isTypeDecl_false(t *testing.T) {
	code := `tpe user str`
	assert.False(t, isTypeDecl(code))
}

func Test_isPrint_true(t *testing.T) {
	assert.True(t, isPrint(`fmt.Println("Salam")`))
	assert.True(t, isPrint("fmt.Printf()"))
	assert.True(t, isPrint("println()"))
	assert.True(t, isPrint("printf()"))
	assert.True(t, isPrint(`fmt.Println(
"Salam"
)`))

}

func Test_isPrint_false(t *testing.T) {
	assert.True(t, true, isPrint("fmt.Fprintln()"))
	assert.True(t, true, isPrint("fmt.Sprintf()"))

}
func Test_isShellCommand(t *testing.T) {
	assert.False(t, isShellCommand(""))
	assert.True(t, isShellCommand(":help"))
	assert.True(t, isShellCommand(":doc"))
	assert.True(t, isShellCommand(":pp"))
}
func Test_reSubMatch(t *testing.T) {
	regx := regexp.MustCompile("func\\s*.+\\s*\\((?P<args>.*)\\) \\((?P<returns>.*)\\)")
	matched, err := reSubMatchMap(regx, "func thisFunc(somearg string) (string, error)")
	assert.NoError(t, err)
	assert.Equal(t, matched["args"], "somearg string")
	assert.Equal(t, matched["returns"], "string, error")
}
