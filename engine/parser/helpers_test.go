package parser

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	assert.True(t, isFunc(code))
}

func Test_isFunc_false(t *testing.T) {
	code := `unc (User) Name string{
		return u.name
}`
	assert.False(t, isFunc(code))
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
	assert.True(t, true, isPrint("fmt.Println()"))
	assert.True(t, true, isPrint("fmt.Printf()"))
	assert.True(t, true, isPrint("println()"))
	assert.True(t, true, isPrint("printf()"))

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
	matched := reSubMatchMap(regx, "func thisFunc(somearg string) (string, error)")
	assert.Equal(t, matched["args"], "somearg string")
	assert.Equal(t, matched["returns"], "string, error")
}
func Test_isFuncCall(t *testing.T) {
	assert.True(t, isFunctionCall("json.Marshal()"))
	assert.True(t, isFunctionCall("println()"))

	assert.False(t, isFunctionCall("2*3(1+2)"))
}

func Test_isExpr(t *testing.T) {
	//assert.True(t, isExpr("1+2"))
	//assert.True(t, isExpr(`"Hello World"`))
	//assert.False(t, isExpr("var x int"))
	assert.False(t, isExpr("x:=2"))
}
func Test_isComment(t *testing.T) {
	assert.True(t, isComment("// salam"))
	assert.True(t, isComment("/* salam */"))
	assert.True(t, isComment(`//fmt.Println("Hello")`))
	assert.False(t, isComment(`fmt.Println("Hello")`))
}
