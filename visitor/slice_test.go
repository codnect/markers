package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSlice_Name(t *testing.T) {
	s := &Slice{
		elem: basicTypesMap["rune"],
	}

	assert.Equal(t, "[]rune", s.Name())
}

func TestSlice_String(t *testing.T) {
	s := &Slice{
		elem: basicTypesMap["byte"],
	}

	assert.Equal(t, "[]byte", s.String())
}

func TestSlice_Elem(t *testing.T) {
	elem := basicTypesMap["byte"]
	s := &Slice{
		elem: elem,
	}

	assert.Equal(t, elem, s.Elem())
	assert.Equal(t, "byte", s.Elem().Name())
	assert.Equal(t, "byte", s.Elem().String())
}

func TestSlice_Underlying(t *testing.T) {
	s := &Slice{
		elem: basicTypesMap["bool"],
	}

	assert.Equal(t, s, s.Underlying())
}
