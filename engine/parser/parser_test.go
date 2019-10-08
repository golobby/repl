package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {

}
func Test_shouldContinue(t *testing.T) {

	code1 := "fmt.Println(\n"
	assert.True(t, ShouldContinue(code1))
	code2 := "fmt.Println(fmt.Sprint(2"
	assert.True(t, ShouldContinue(code2))
	code3 := "{fmt.Print("
	assert.True(t, ShouldContinue(code3))
	code4 := "fmt.Println(22)"
	assert.False(t, ShouldContinue(code4))
}
