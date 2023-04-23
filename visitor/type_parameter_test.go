package visitor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeParameter_Name(t *testing.T) {
	constraint := &TypeConstraint{
		typ:           basicTypesMap["bool"],
		tildeOperator: true,
	}
	typeConstraints := &TypeConstraints{
		[]*TypeConstraint{constraint},
	}
	typeParameter := &TypeParameter{
		name:        "anyTypeParameter",
		constraints: typeConstraints,
	}

	assert.Equal(t, typeParameter.name, typeParameter.Name())
}

func TestTypeParameter_String(t *testing.T) {
	constraint := &TypeConstraint{
		typ:           basicTypesMap["bool"],
		tildeOperator: true,
	}
	typeConstraints := &TypeConstraints{
		[]*TypeConstraint{constraint},
	}
	typeParameter := &TypeParameter{
		name:        "anyTypeParameter",
		constraints: typeConstraints,
	}

	assert.Equal(t, fmt.Sprintf("%s ~%s", typeParameter.name, "bool"), typeParameter.String())
}

func TestTypeParameter_TypeConstraints(t *testing.T) {
	constraint := &TypeConstraint{
		typ:           basicTypesMap["bool"],
		tildeOperator: true,
	}
	typeConstraints := &TypeConstraints{
		[]*TypeConstraint{constraint},
	}
	typeParameter := &TypeParameter{
		name:        "anyTypeParameter",
		constraints: typeConstraints,
	}
	assert.Equal(t, typeConstraints, typeParameter.TypeConstraints())
}

func TestTypeParameter_Underlying(t *testing.T) {
	constraint := &TypeConstraint{
		typ:           basicTypesMap["bool"],
		tildeOperator: true,
	}
	typeConstraints := &TypeConstraints{
		[]*TypeConstraint{constraint},
	}
	typeParameter := &TypeParameter{
		name:        "anyTypeParameter",
		constraints: typeConstraints,
	}
	assert.Equal(t, typeParameter, typeParameter.Underlying())
}

func TestTypeParameters_At(t *testing.T) {
	constraint := &TypeConstraint{
		typ:           basicTypesMap["bool"],
		tildeOperator: true,
	}
	typeConstraints := &TypeConstraints{
		[]*TypeConstraint{constraint},
	}
	typeParameter := &TypeParameter{
		name:        "anyTypeParameter",
		constraints: typeConstraints,
	}
	typeParameters := &TypeParameters{
		elements: []*TypeParameter{typeParameter},
	}

	assert.NotNil(t, typeParameters.At(0))
	assert.Nil(t, typeParameters.At(1))
}

func TestTypeParameters_Len(t *testing.T) {
	constraint := &TypeConstraint{
		typ:           basicTypesMap["bool"],
		tildeOperator: true,
	}
	typeConstraints := &TypeConstraints{
		[]*TypeConstraint{constraint},
	}
	typeParameter := &TypeParameter{
		name:        "anyTypeParameter",
		constraints: typeConstraints,
	}
	typeParameters := &TypeParameters{
		elements: []*TypeParameter{typeParameter},
	}

	assert.Equal(t, len(typeParameters.elements), typeParameters.Len())
}
