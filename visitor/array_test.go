package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArray_Name(t *testing.T) {
	a := &Array{
		len:  5,
		elem: basicTypesMap["bool"],
	}

	assert.Equal(t, "[5]bool", a.Name())

	twoDimensionalArray := &Array{
		len: 5,
		elem: &Array{
			len:  6,
			elem: basicTypesMap["int32"],
		},
	}

	assert.Equal(t, "[5][6]int32", twoDimensionalArray.String())
}

func TestArray_String(t *testing.T) {
	a := &Array{
		len:  5,
		elem: basicTypesMap["bool"],
	}

	assert.Equal(t, "[5]bool", a.String())

	twoDimensionalArray := &Array{
		len: 5,
		elem: &Array{
			len:  6,
			elem: basicTypesMap["int32"],
		},
	}

	assert.Equal(t, "[5][6]int32", twoDimensionalArray.Name())
}

func TestArray_Elem(t *testing.T) {
	elem := basicTypesMap["byte"]
	a := &Array{
		len:  5,
		elem: elem,
	}

	assert.Equal(t, elem, a.Elem())
	assert.Equal(t, "byte", a.Elem().Name())
	assert.Equal(t, "byte", a.Elem().String())
}

func TestArray_Len(t *testing.T) {
	elem := basicTypesMap["byte"]
	a := &Array{
		len:  5,
		elem: elem,
	}

	assert.Equal(t, a.len, a.Len())
}

func TestArray_Underlying(t *testing.T) {
	elem := basicTypesMap["byte"]
	a := &Array{
		len:  5,
		elem: elem,
	}

	assert.Equal(t, a, a.Underlying())
}
