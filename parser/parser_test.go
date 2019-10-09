package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	typ, err := Parse(``)
	assert.NoError(t, err)
	assert.Equal(t, StmtEmpty, typ)
	typ, err = Parse(`// salam`)
	assert.NoError(t, err)
	assert.Equal(t, StmtTypeComment, typ)
	typ, err = Parse(`import "fmt"`)
	assert.NoError(t, err)
	assert.Equal(t, StmtTypeImport, typ)
	typ, err = Parse(`func name() string { return "" }`)
	assert.NoError(t, err)
	assert.Equal(t, StmtTypeFuncDecl, typ)
	typ, err = Parse(`type user struct{ Name string }`)
	assert.NoError(t, err)
	assert.Equal(t, StmtTypeTypeDecl, typ)
	typ, err = Parse(`fmt.Println("aleyk")`)
	assert.NoError(t, err)
	assert.Equal(t, StmtTypePrint, typ)
	typ, err = Parse(`fmt.Println`)
	assert.NoError(t, err)
	assert.Equal(t, StmtTypeExpr, typ)
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
