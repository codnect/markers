package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPointer_Name(t *testing.T) {
	p := &Pointer{
		base: basicTypesMap["bool"],
	}

	assert.Equal(t, "*bool", p.Name())
}

func TestPointer_String(t *testing.T) {
	p := &Pointer{
		base: basicTypesMap["bool"],
	}

	assert.Equal(t, "*bool", p.String())
}

func TestPointer_Elem(t *testing.T) {
	elem := basicTypesMap["byte"]
	p := &Pointer{
		base: elem,
	}

	assert.Equal(t, elem, p.Elem())
	assert.Equal(t, "byte", p.Elem().Name())
	assert.Equal(t, "byte", p.Elem().String())
}

func TestPointer_Underlying(t *testing.T) {
	p := &Pointer{
		base: basicTypesMap["bool"],
	}

	assert.Equal(t, p, p.Underlying())
}
