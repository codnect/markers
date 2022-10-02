package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasic_Kind(t *testing.T) {
	b := basicTypesMap["bool"]
	assert.Equal(t, Bool, b.Kind())
}

func TestBasic_Name(t *testing.T) {
	b := basicTypesMap["bool"]
	assert.Equal(t, "bool", b.Name())
}

func TestBasic_String(t *testing.T) {
	b := basicTypesMap["bool"]
	assert.Equal(t, "bool", b.String())
}

func TestBasic_Underlying(t *testing.T) {
	b := basicTypesMap["bool"]
	assert.Equal(t, b, b.Underlying())
}
