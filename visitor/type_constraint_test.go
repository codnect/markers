package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeConstraint_Name(t *testing.T) {
	c := &TypeConstraint{
		typ:           basicTypesMap["bool"],
		tildeOperator: true,
	}

	assert.Equal(t, "bool", c.Name())
}

func TestTypeConstraint_String(t *testing.T) {
	c := &TypeConstraint{
		typ:           basicTypesMap["bool"],
		tildeOperator: true,
	}

	assert.Equal(t, "~bool", c.String())
}

func TestTypeConstraint_Type(t *testing.T) {
	typ := basicTypesMap["byte"]
	c := &TypeConstraint{
		typ:           typ,
		tildeOperator: true,
	}

	assert.Equal(t, typ, c.Type())
}

func TestTypeConstraint_Underlying(t *testing.T) {
	c := &TypeConstraint{
		typ:           basicTypesMap["bool"],
		tildeOperator: true,
	}

	assert.Equal(t, c, c.Underlying())
}

func TestTypeConstraints_At(t *testing.T) {
	constraint := &TypeConstraint{
		typ:           basicTypesMap["bool"],
		tildeOperator: true,
	}

	typeConstraints := &TypeConstraints{
		elements: []*TypeConstraint{constraint},
	}

	assert.NotNil(t, typeConstraints.At(0))
	assert.Nil(t, typeConstraints.At(1))
}

func TestTypeConstraints_Len(t *testing.T) {
	constraint := &TypeConstraint{
		typ:           basicTypesMap["bool"],
		tildeOperator: true,
	}

	typeConstraints := &TypeConstraints{
		elements: []*TypeConstraint{constraint},
	}

	assert.Equal(t, len(typeConstraints.elements), typeConstraints.Len())
}
