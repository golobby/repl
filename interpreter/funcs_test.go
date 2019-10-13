package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExtractFuncName(t *testing.T) {
	assert.Equal(t, "name", ExtractFuncName(`func name(a int) string {
		return ""
}`))
}
