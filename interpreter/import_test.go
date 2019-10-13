package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isImport(t *testing.T) {
	assert.True(t, isImport(`import "fmt"`))
	assert.False(t, isImport(`aaaa`))
}

func Test_extractImportData(t *testing.T) {
	assert.Equal(t, []ImportData{
		{
			Path:  `"fmt"`,
			Alias: "",
		},
	}, ExtractImportData(`import "fmt"`))
	assert.Equal(t, []ImportData{
		{
			Path:  `"fmt"`,
			Alias: "",
		},
		{
			Path:  `"os"`,
			Alias: "",
		},
	}, ExtractImportData("import (\n\"fmt\"\n\"os\"\n)"))
}
