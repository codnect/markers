package visitor

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type customTypeInfo struct {
	name          string
	aliasTypeName string
	isExported    bool
}

var (
	errorCustomTypes = map[string]customTypeInfo{
		"errorList": {
			name:          "errorList",
			aliasTypeName: "[]error",
			isExported:    false,
		},
	}
	permissionCustomTypes = map[string]customTypeInfo{
		"Permission": {
			name:          "Permission",
			aliasTypeName: "int",
			isExported:    true,
		},
		"RequestMethod": {
			name:          "RequestMethod",
			aliasTypeName: "string",
			isExported:    true,
		},
		"Chan": {
			name:          "Chan",
			aliasTypeName: "int",
			isExported:    true,
		},
	}
	coffeeCustomTypes = map[string]customTypeInfo{
		"Coffee": {
			name:          "Coffee",
			aliasTypeName: "int",
			isExported:    true,
		},
	}
	freshCustomTypes = map[string]customTypeInfo{
		"Lemonade": {
			name:          "Lemonade",
			aliasTypeName: "uint",
			isExported:    true,
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

		if fileCustomType != file.CustomTypes().At(index) {
			t.Errorf("custom type with name %s does not match with custom type at index %d", fileCustomType.Name(), index)
			continue
		}

		actualCustomType, exists := file.CustomTypes().FindByName(expectedCustomTypeName)
		if !exists || actualCustomType == nil {
			t.Errorf("custom type with name %s in file %s is not found", expectedCustomTypeName, file.name)
			continue
		}

		assert.Equal(t, actualCustomType, actualCustomType.Underlying())
		assert.Equal(t, fileCustomType, actualCustomType, "CustomTypes.At should return %w, but got %w", fileCustomType, actualCustomType)

		if expectedCustomType.name != actualCustomType.Name() {
			t.Errorf("custom type name in file %s shoud be %s, but got %s", file.name, expectedCustomTypeName, actualCustomType.Name())
		}

		if expectedCustomType.aliasTypeName != actualCustomType.AliasType().Name() {
			t.Errorf("alias type of custom type %s in file %s shoud be %s, but got %s", file.name, expectedCustomType.name, expectedCustomType.aliasTypeName, actualCustomType.AliasType().Name())
		}

		customTypeStrValue := fmt.Sprintf("type %s %s", expectedCustomType.name, expectedCustomType.aliasTypeName)
		if customTypeStrValue != actualCustomType.String() {
			t.Errorf("String() method of custom type %s shoud return %s, but got %s", expectedCustomTypeName, customTypeStrValue, actualCustomType.String())
		}

		if actualCustomType.IsExported() && !expectedCustomType.isExported {
			t.Errorf("custom type with name %s is exported, but should be unexported field", expectedCustomTypeName)
		} else if !actualCustomType.IsExported() && expectedCustomType.isExported {
			t.Errorf("custom type with name %s is not exported, but should be exported field", expectedCustomTypeName)
		}

		/*if expectedConstant.value != actualConstant.Value() {
			t.Errorf("value of constant %s in file %s shoud be %s, but got %s", actualConstant.Name(), file.name, expectedConstant.value, actualConstant.Value())
		}

		if expectedConstant.typeName != actualConstant.Type().Name() {
			t.Errorf("type name of constant %s in file %s shoud be %s, but got %s", actualConstant.Name(), file.name, expectedConstant.typeName, actualConstant.Type().Name())
		}

		assert.Equal(t, expectedConstant.position, actualConstant.Position(), "the position of constant %s in file %s should be %w, but got %w", expectedConstant.name, actualConstant.File().Name(), expectedConstant.position, actualConstant.Position())
		*/
		index++
	}

	return true
}
