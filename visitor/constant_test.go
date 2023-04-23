package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type constantInfo struct {
	name       string
	position   Position
	value      any
	typeName   string
	isExported bool
}

var (
	coffeeConstants = []constantInfo{
		{
			name: "Cappuccino",
			position: Position{
				Line:   10,
				Column: 2,
			},
			value:      -1,
			typeName:   "Coffee",
			isExported: true,
		},
		{
			name: "Americano",
			position: Position{
				Line:   11,
				Column: 2,
			},
			value:      -2,
			typeName:   "Coffee",
			isExported: true,
		},
		{
			name: "Latte",
			position: Position{
				Line:   12,
				Column: 2,
			},
			value:      -3,
			typeName:   "Coffee",
			isExported: true,
		},
		{
			name: "TurkishCoffee",
			position: Position{
				Line:   13,
				Column: 2,
			},
			value:      -4,
			typeName:   "Coffee",
			isExported: true,
		},
	}
	freshConstants = []constantInfo{
		{
			name: "ClassicLemonade",
			position: Position{
				Line:   10,
				Column: 2,
			},
			value:      0,
			typeName:   "Lemonade",
			isExported: true,
		},
		{
			name: "BlueberryLemonade",
			position: Position{
				Line:   11,
				Column: 2,
			},
			value:      1,
			typeName:   "Lemonade",
			isExported: true,
		},
		{
			name: "WatermelonLemonade",
			position: Position{
				Line:   12,
				Column: 2,
			},
			value:      2,
			typeName:   "Lemonade",
			isExported: true,
		},
		{
			name: "MangoLemonade",
			position: Position{
				Line:   13,
				Column: 2,
			},
			value:      3,
			typeName:   "Lemonade",
			isExported: true,
		},
		{
			name: "StrawberryLemonade",
			position: Position{
				Line:   14,
				Column: 2,
			},
			value:      4,
			typeName:   "Lemonade",
			isExported: true,
		},
	}
	stringConstants = []constantInfo{
		{
			name: "StringOperation",
			position: Position{
				Line:   5,
				Column: 7,
			},
			value:      "AnyString",
			typeName:   "string",
			isExported: true,
		},
		{
			name: "methods",
			position: Position{
				Line:   6,
				Column: 7,
			},
			value:      "GETPUT",
			typeName:   "string",
			isExported: false,
		},
	}
	permissionConstants = []constantInfo{
		{
			name: "Read",
			position: Position{
				Line:   10,
				Column: 2,
			},
			value:      1,
			typeName:   "Permission",
			isExported: true,
		},
		{
			name: "Write",
			position: Position{
				Line:   11,
				Column: 2,
			},
			value:      2,
			typeName:   "Permission",
			isExported: true,
		},
		{
			name: "ReadWrite",
			position: Position{
				Line:   12,
				Column: 2,
			},
			value:      3,
			typeName:   "Permission",
			isExported: true,
		},
		{
			name: "RequestGet",
			position: Position{
				Line:   18,
				Column: 2,
			},
			value:      "GET",
			typeName:   "RequestMethod",
			isExported: true,
		},
		{
			name: "RequestPost",
			position: Position{
				Line:   19,
				Column: 2,
			},
			value:      "POST",
			typeName:   "RequestMethod",
			isExported: true,
		},
		{
			name: "RequestPatch",
			position: Position{
				Line:   20,
				Column: 2,
			},
			value:      "PATCH",
			typeName:   "RequestMethod",
			isExported: true,
		},
		{
			name: "RequestDelete",
			position: Position{
				Line:   21,
				Column: 2,
			},
			value:      "DELETE",
			typeName:   "RequestMethod",
			isExported: true,
		},
		{
			name: "SendDir",
			position: Position{
				Line:   27,
				Column: 2,
			},
			value:      2,
			typeName:   "Chan",
			isExported: true,
		},
		{
			name: "ReceiveDir",
			position: Position{
				Line:   28,
				Column: 2,
			},
			value:      1,
			typeName:   "Chan",
			isExported: true,
		},
		{
			name: "BothDir",
			position: Position{
				Line:   29,
				Column: 2,
			},
			value:      3,
			typeName:   "Chan",
			isExported: true,
		},
	}

	mathConstants = []constantInfo{
		{
			name: "IntegerMathOperation",
			position: Position{
				Line:   3,
				Column: 7,
			},
			value:      -4,
			typeName:   "untyped int",
			isExported: true,
		},
		{
			name: "floatMathOperation",
			position: Position{
				Line:   4,
				Column: 7,
			},
			value:      5.4,
			typeName:   "untyped int",
			isExported: false,
		},
		{
			name: "ModOperation",
			position: Position{
				Line:   5,
				Column: 7,
			},
			value:      1,
			typeName:   "untyped int",
			isExported: true,
		},
		{
			name: "equalOperation",
			position: Position{
				Line:   6,
				Column: 7,
			},
			value:      true,
			typeName:   "bool",
			isExported: false,
		},
		{
			name: "NotEqualOperation",
			position: Position{
				Line:   7,
				Column: 7,
			},
			value:      false,
			typeName:   "bool",
			isExported: true,
		},
		{
			name: "GreaterThan",
			position: Position{
				Line:   8,
				Column: 7,
			},
			value:      true,
			typeName:   "bool",
			isExported: true,
		},
		{
			name: "GreaterThanOrEqual",
			position: Position{
				Line:   9,
				Column: 7,
			},
			value:      true,
			typeName:   "bool",
			isExported: true,
		},
		{
			name: "LessThan",
			position: Position{
				Line:   10,
				Column: 7,
			},
			value:      true,
			typeName:   "bool",
			isExported: true,
		},
		{
			name: "LessThanOrEqual",
			position: Position{
				Line:   11,
				Column: 7,
			},
			value:      true,
			typeName:   "bool",
			isExported: true,
		},
		{
			name: "XorOperation",
			position: Position{
				Line:   12,
				Column: 7,
			},
			value:      6,
			typeName:   "untyped int",
			isExported: true,
		},
		{
			name: "AndNotOperation",
			position: Position{
				Line:   13,
				Column: 7,
			},
			value:      4,
			typeName:   "untyped int",
			isExported: true,
		},
		{
			name: "AndOperation",
			position: Position{
				Line:   14,
				Column: 7,
			},
			value:      0,
			typeName:   "untyped int",
			isExported: true,
		},
		{
			name: "orOperation",
			position: Position{
				Line:   15,
				Column: 7,
			},
			value:      6,
			typeName:   "untyped int",
			isExported: false,
		},
	}
)

func assertConstants(t *testing.T, file *File, constants []constantInfo) bool {
	if file.Constants().Len() != len(constants) {
		t.Errorf("the number of the constants in file %s should be %d, but got %d", file.Name(), len(constants), file.Constants().Len())
	}

	assert.Equal(t, file.Constants().elements, file.Constants().ToSlice(), "ToSlice should return %w, but got %w", file.Constants().elements, file.Constants().ToSlice())

	for index, expectedConstant := range constants {
		fileConstant := file.Constants().At(index)

		actualConstant, exists := file.Constants().FindByName(expectedConstant.name)
		if !exists || actualConstant == nil {
			t.Errorf("constant with name %s in file %s is not found", expectedConstant.name, file.name)
			continue
		}

		assert.Equal(t, fileConstant, actualConstant, "Constants.At should return %w, but got %w", fileConstant, actualConstant)

		if expectedConstant.name != actualConstant.Name() {
			t.Errorf("constant name in file %s shoud be %s, but got %s", file.name, expectedConstant.name, actualConstant.Name())
		}

		if expectedConstant.value != actualConstant.Value() {
			t.Errorf("value of constant %s in file %s shoud be %s, but got %s", actualConstant.Name(), file.name, expectedConstant.value, actualConstant.Value())
		}

		if expectedConstant.typeName != actualConstant.Type().Name() {
			t.Errorf("type name of constant %s in file %s shoud be %s, but got %s", actualConstant.Name(), file.name, expectedConstant.typeName, actualConstant.Type().Name())
		}

		if actualConstant.IsExported() && !expectedConstant.isExported {
			t.Errorf("constant with name %s in file %s is exported, but should be unexported field", expectedConstant.name, file.name)
		} else if !actualConstant.IsExported() && expectedConstant.isExported {
			t.Errorf("constant with name %s in file %s is not exported, but should be exported field", expectedConstant.name, file.name)
		}

		assert.Equal(t, expectedConstant.position, actualConstant.Position(), "the position of constant %s in file %s should be %w, but got %w", expectedConstant.name, actualConstant.File().Name(), expectedConstant.position, actualConstant.Position())
	}

	return true
}

func TestConstants_AtShouldReturnNilIfIndexIsOutOfRange(t *testing.T) {
	constants := &Constants{}
	assert.Nil(t, constants.At(0))
}
