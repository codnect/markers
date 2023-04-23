package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestImports_AtShouldReturnNilIfIndexIsOutOfRange(t *testing.T) {
	imports := &Imports{}
	assert.Nil(t, imports.At(0))
}

func TestImports_AtShouldReturnNilIfPathIsNotFound(t *testing.T) {
	imports := &Imports{}
	imp, ok := imports.FindByPath("anyPath")

	assert.Nil(t, imp)
	assert.False(t, ok)
}
