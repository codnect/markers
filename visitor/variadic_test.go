package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVariadic_Name(t *testing.T) {
	v := &Variadic{
		elem: basicTypesMap["bool"],
	}

	assert.Equal(t, "bool", v.Name())
}

func TestVariadic_String(t *testing.T) {
	v := &Variadic{
		elem: basicTypesMap["bool"],
	}

	assert.Equal(t, "...bool", v.String())
}

func TestVariadic_Elem(t *testing.T) {
	elem := basicTypesMap["byte"]
	a := &Variadic{
		elem: elem,
	}

	assert.Equal(t, elem, a.Elem())
	assert.Equal(t, "byte", a.Elem().Name())
	assert.Equal(t, "byte", a.Elem().String())
}

func TestVariadic_Underlying(t *testing.T) {
	elem := basicTypesMap["byte"]
	v := &Variadic{
		elem: elem,
	}

	assert.Equal(t, v, v.Underlying())
}
