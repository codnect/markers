package visitor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type customTypeInfo struct {
	name               string
	underlyingTypeName string
	isExported         bool
	methods            map[string]functionInfo
	stringValue        string
}

var (
	errorCustomTypes = map[string]customTypeInfo{
		"errorList": {
			name:               "errorList",
			underlyingTypeName: "[]error",
			isExported:         false,
			methods: map[string]functionInfo{
				"Print":    printErrorMethod,
				"ToErrors": toErrorsMethod,
			},
			stringValue: "any.errorList",
		},
	}
	permissionCustomTypes = map[string]customTypeInfo{
		"Permission": {
			name:               "Permission",
			underlyingTypeName: "int",
			isExported:         true,
			stringValue:        "any.Permission",
		},
		"RequestMethod": {
			name:               "RequestMethod",
			underlyingTypeName: "string",
			isExported:         true,
			stringValue:        "any.RequestMethod",
		},
		"Chan": {
			name:               "Chan",
			underlyingTypeName: "int",
			isExported:         true,
			stringValue:        "any.Chan",
		},
	}
	coffeeCustomTypes = map[string]customTypeInfo{
		"Coffee": {
			name:               "Coffee",
			underlyingTypeName: "int",
			isExported:         true,
			stringValue:        "menu.Coffee",
		},
	}
	freshCustomTypes = map[string]customTypeInfo{
		"Lemonade": {
			name:               "Lemonade",
			underlyingTypeName: "uint",
			isExported:         true,
			stringValue:        "menu.Lemonade",
		},
	}
	genericsCustomTypes = map[string]customTypeInfo{
		"HttpHandler": {
			name:               "HttpHandler",
			underlyingTypeName: "func (ctx C) K",
			isExported:         true,
			methods: map[string]functionInfo{
				"Print": printHttpHandlerMethod,
			},
			stringValue: "any.HttpHandler[C context.Context,K string|int]",
		},
	}
)

func assertCustomTypes(t *testing.T, file *File, customTypes map[string]customTypeInfo) bool {
	if file.CustomTypes().Len() != len(customTypes) {
		t.Errorf("the number of the custom types in file %s should be %d, but got %d", file.Name(), len(customTypes), file.CustomTypes().Len())
	}

	assert.Equal(t, file.CustomTypes().elements, file.CustomTypes().ToSlice(), "ToSlice should return %w, but got %w", file.Constants().elements, file.Constants().ToSlice())

	index := 0
	for expectedCustomTypeName, expectedCustomType := range customTypes {
		fileCustomType, ok := file.CustomTypes().FindByName(expectedCustomTypeName)
		if !ok {
			t.Errorf("custom type with name %s is not found", expectedCustomTypeName)
			continue
		}

		if file.CustomTypes().elements[index] != file.CustomTypes().At(index) {
			t.Errorf("custom type with name %s does not match with custom type at index %d", fileCustomType.Name(), index)
			continue
		}

		actualCustomType, exists := file.CustomTypes().FindByName(expectedCustomTypeName)
		if !exists || actualCustomType == nil {
			t.Errorf("custom type with name %s in file %s is not found", expectedCustomTypeName, file.name)
			continue
		}

		assert.Equal(t, fileCustomType, actualCustomType, "CustomTypes.At should return %w, but got %w", fileCustomType, actualCustomType)

		if expectedCustomType.name != actualCustomType.Name() {
			t.Errorf("custom type name in file %s shoud be %s, but got %s", file.name, expectedCustomTypeName, actualCustomType.Name())
		}

		if expectedCustomType.underlyingTypeName != actualCustomType.Underlying().String() {
			t.Errorf("underlying type of custom type %s in file %s shoud be %s, but got %s", file.name, expectedCustomType.name, expectedCustomType.underlyingTypeName, actualCustomType.Underlying().String())
		}

		if actualCustomType.IsExported() && !expectedCustomType.isExported {
			t.Errorf("custom type with name %s is exported, but should be unexported field", expectedCustomTypeName)
		} else if !actualCustomType.IsExported() && expectedCustomType.isExported {
			t.Errorf("custom type with name %s is not exported, but should be exported field", expectedCustomTypeName)
		}

		if expectedCustomType.stringValue != actualCustomType.String() {
			t.Errorf("Output returning from String() method for custom type with name %s does not equal to %s, but got %s", expectedCustomTypeName, expectedCustomType.stringValue, actualCustomType.String())
		}

		if actualCustomType.NumMethods() != len(expectedCustomType.methods) {
			t.Errorf("the number of the methods of the custom type %s should be %d, but got %d", expectedCustomTypeName, len(expectedCustomType.methods), actualCustomType.NumMethods())
		}

		assertFunctions(t, fmt.Sprintf("custom type %s", actualCustomType.Name()), actualCustomType.Methods(), expectedCustomType.methods)
		index++
	}

	return true
}

func TestCustomTypes_AtShouldReturnNilIfIndexIsOutOfRange(t *testing.T) {
	customTypes := &CustomTypes{}
	assert.Nil(t, customTypes.At(0))
}
