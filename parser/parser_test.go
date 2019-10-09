package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

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
