package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isTypeDecl(t *testing.T) {
	assert.True(t, isTypeDecl("type user struct{\n Name string \n}"))
	assert.True(t, isTypeDecl("type user interface{\n Name() string \n}"))
	assert.False(t, isTypeDecl("aaaaa"))
}

func Test_ExtractTypeName(t *testing.T) {
	assert.Equal(t, "user", ExtractTypeName("type user struct{\n Name string \n}"))
	assert.Equal(t, "user", ExtractTypeName("type user interface{\n Name() string \n}"))
}
