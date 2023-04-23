package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypes_FindByNameShouldReturnTypeIfItExists(t *testing.T) {
	types := &Types{
		elements: []Type{
			&Struct{
				name: "test",
			},
		},
	}

	typ, ok := types.FindByName("test")
	assert.True(t, ok)
	assert.NotNil(t, typ)
	assert.Equal(t, "test", typ.Name())
}

func TestTypes_FindByNameShouldReturnNilIfItDoesExist(t *testing.T) {
	types := &Types{
		elements: []Type{},
	}

	typ, ok := types.FindByName("test")
	assert.False(t, ok)
	assert.Nil(t, typ)
}

func TestTypes_AtShouldReturnNilIfGivenIndexIsOutOfRange(t *testing.T) {
	types := &Types{
		elements: []Type{},
	}

	typ := types.At(-1)
	assert.Nil(t, typ)
}

func TestTypeSets_AtShouldReturnNilIfGivenIndexIsOutOfRange(t *testing.T) {
	typeSets := &TypeSets{}

	typ := typeSets.At(-1)
	assert.Nil(t, typ)
}
