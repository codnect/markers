package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenericType_Name(t *testing.T) {
	g := &GenericType{
		rawType: &Struct{
			name: "Test",
		},
		arguments: TypeSets{basicTypesMap["int"], basicTypesMap["string"]},
	}

	assert.Equal(t, "Test", g.Name())
}

func TestGenericType_String(t *testing.T) {
	g := &GenericType{
		rawType: &Struct{
			name: "Test",
		},
		arguments: TypeSets{basicTypesMap["int"], basicTypesMap["string"]},
	}

	assert.Equal(t, "Test[int,string]", g.String())
}

func TestGenericType_RawType(t *testing.T) {
	rawType := &Struct{
		name: "",
		typeParams: &TypeParameters{
			[]*TypeParameter{},
		},
	}
	g := &GenericType{
		rawType: rawType,
	}

	assert.Equal(t, rawType, g.RawType())
	assert.Equal(t, "struct{}", g.RawType().Name())
	assert.Equal(t, "struct{}", g.RawType().String())
}

func TestGenericType_Underlying(t *testing.T) {
	g := &GenericType{
		rawType: &Struct{
			name: "Test",
		},
		arguments: TypeSets{basicTypesMap["int"], basicTypesMap["string"]},
	}

	assert.Equal(t, g, g.Underlying())
}
