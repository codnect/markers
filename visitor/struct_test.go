package visitor

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fieldInfo struct {
	name            string
	typeName        string
	isExported      bool
	isEmbeddedField bool
}

type structInfo struct {
	fileName          string
	isExported        bool
	position          Position
	markers           marker.MarkerValues
	methods           map[string]functionInfo
	allMethods        map[string]functionInfo
	fields            map[string]fieldInfo
	embeddedFields    map[string]fieldInfo
	numFields         int
	totalFields       int
	numEmbeddedFields int
	implements        map[string]struct{}
}

// structs
var (
	friedCookieStruct = structInfo{
		markers: marker.MarkerValues{
			"marker:struct-type-level": {
				StructTypeLevel{
					Name: "FriedCookie",
				},
			},
		},
		fileName:   "dessert.go",
		isExported: true,
		position: Position{
			Line:   30,
			Column: 6,
		},
		methods: map[string]functionInfo{
			"Eat": eatMethod,
			"Buy": buyMethod,
		},
		allMethods: map[string]functionInfo{
			"Eat":           eatMethod,
			"Buy":           buyMethod,
			"Oreo":          oreoMethod,
			"FortuneCookie": fortuneCookieMethod,
		},
		fields: map[string]fieldInfo{
			"cookie": {
				isExported:      false,
				isEmbeddedField: true,
				typeName:        "cookie",
			},
			"cookieDough": {
				isExported:      false,
				isEmbeddedField: false,
				typeName:        "any",
			},
		},
		embeddedFields: map[string]fieldInfo{
			"cookie": {
				isExported:      false,
				isEmbeddedField: true,
				typeName:        "cookie",
			},
		},
		numFields:         2,
		totalFields:       3,
		numEmbeddedFields: 1,
	}

	cookieStruct = structInfo{
		markers: marker.MarkerValues{
			"marker:struct-type-level": {
				StructTypeLevel{
					Name: "cookie",
					Any: map[string]interface{}{
						"key": "value",
					},
				},
			},
		},
		fileName:   "dessert.go",
		isExported: false,
		position: Position{
			Line:   56,
			Column: 6,
		},
		methods: map[string]functionInfo{
			"FortuneCookie": fortuneCookieMethod,
			"Oreo":          oreoMethod,
		},
		allMethods: map[string]functionInfo{
			"FortuneCookie": fortuneCookieMethod,
			"Oreo":          oreoMethod,
		},
		fields: map[string]fieldInfo{
			"ChocolateChip": {
				isExported:      true,
				isEmbeddedField: false,
				typeName:        "string",
			},
			"tripleChocolateCookie": {
				isExported:      false,
				isEmbeddedField: false,
				typeName:        "map[string]error",
			},
		},
		embeddedFields:    map[string]fieldInfo{},
		numFields:         2,
		totalFields:       2,
		numEmbeddedFields: 0,
	}
)

func assertStructs(t *testing.T, file *File, structs map[string]structInfo) bool {

	if len(structs) != file.Structs().Len() {
		t.Errorf("the number of the functions should be %d, but got %d", len(structs), file.Structs().Len())
		return false
	}

	index := 0
	for expectedStructName, expectedStruct := range structs {
		actualStruct, ok := file.Structs().FindByName(expectedStructName)
		if !ok {
			t.Errorf("struct with name %s is not found", expectedStructName)
			continue
		}

		if actualStruct.NamedType() == nil {
			t.Errorf("NamedType() for struct %s should not return nil", actualStruct.Name())
		}

		if expectedStruct.fileName != actualStruct.File().Name() {
			t.Errorf("the file name for struct %s should be %s, but got %s", expectedStructName, expectedStruct.fileName, actualStruct.File().Name())
		}

		if file.Structs().elements[index] != file.Structs().At(index) {
			t.Errorf("struct with name %s does not match with struct at index %d", actualStruct.Name(), index)
			continue
		}

		if actualStruct.IsExported() && !expectedStruct.isExported {
			t.Errorf("struct with name %s is exported, but should be unexported", actualStruct.Name())
		} else if !actualStruct.IsExported() && expectedStruct.isExported {
			t.Errorf("struct with name %s is not exported, but should be exported", actualStruct.Name())
		}

		if actualStruct.NumMethods() == 0 && !actualStruct.IsEmpty() {
			t.Errorf("the struct %s should be empty", actualStruct.Name())
		} else if actualStruct.NumMethods() != 0 && actualStruct.IsEmpty() {
			t.Errorf("the struct %s should not be empty", actualStruct.Name())
		}

		if actualStruct.NumMethodsInHierarchy() != len(expectedStruct.allMethods) {
			t.Errorf("the number of the methods of the struct %s should be %d, but got %d", expectedStructName, len(expectedStruct.allMethods), actualStruct.NumMethodsInHierarchy())
		}

		if actualStruct.NumMethods() != len(expectedStruct.methods) {
			t.Errorf("the number of the methods of the struct %s should be %d, but got %d", expectedStructName, len(expectedStruct.methods), actualStruct.NumMethods())
		}

		if actualStruct.NumFields() != expectedStruct.numFields {
			t.Errorf("the number of the fields of the struct %s should be %d, but got %d", expectedStructName, expectedStruct.totalFields-expectedStruct.numEmbeddedFields, actualStruct.NumFields())
		}

		if actualStruct.NumFieldsInHierarchy() != expectedStruct.totalFields {
			t.Errorf("the number of the all fields of the struct %s should be %d, but got %d", expectedStructName, expectedStruct.totalFields, actualStruct.NumFields())
		}

		if actualStruct.NumEmbeddedFields() != expectedStruct.numEmbeddedFields {
			t.Errorf("the number of the embededed fields of the struct %s should be %d, but got %d", expectedStructName, expectedStruct.numEmbeddedFields, actualStruct.NumFields())
		}

		assert.Equal(t, actualStruct, actualStruct.Underlying())

		assert.Equal(t, expectedStruct.position, actualStruct.Position(), "the position of the struct %s should be %w, but got %w",
			expectedStructName, expectedStruct.position, actualStruct.Position())

		assertFunctions(t, fmt.Sprintf("struct %s", actualStruct.Name()), actualStruct.Methods(), expectedStruct.methods)
		assertFunctions(t, fmt.Sprintf("struct %s", actualStruct.Name()), actualStruct.MethodsInHierarchy(), expectedStruct.allMethods)
		assertStructFields(t, actualStruct.Name(), actualStruct.EmbeddedFields(), expectedStruct.embeddedFields)
		assertStructFields(t, actualStruct.Name(), actualStruct.Fields(), expectedStruct.fields)
		assertMarkers(t, expectedStruct.markers, actualStruct.Markers(), fmt.Sprintf("struct %s", expectedStructName))

		index++
	}

	return true
}

func assertStructFields(t *testing.T, structName string, actualFields *Fields, expectedFields map[string]fieldInfo) bool {

	for expectedFieldName, expectedField := range expectedFields {
		actualField, ok := actualFields.FindByName(expectedFieldName)

		if !ok {
			t.Errorf("field with name %s for struct %s is not found", expectedFieldName, structName)
			continue
		}

		if actualField.Name() != expectedFieldName {
			t.Errorf("field name for struct %s shoud be %s, but got %s", structName, expectedFieldName, actualField.Name())
		}

		if actualField.Type().Name() != expectedField.typeName {
			t.Errorf("type of field with name %s for struct %s shoud be %s, but got %s", actualField.Name(), structName, expectedField.typeName, actualField.Type().Name())
		}

		if actualField.IsExported() && !expectedField.isExported {
			t.Errorf("field with name %s for struct %s is exported, but should be unexported field", expectedFieldName, structName)
		} else if !actualField.IsExported() && expectedField.isExported {
			t.Errorf("field with name %s for struct %s is not exported, but should be exported field", expectedFieldName, structName)
		}

		if actualField.IsEmbedded() && !expectedField.isEmbeddedField {
			t.Errorf("field with name %s for struct %s is embedded, but should be not embedded field", expectedFieldName, structName)
		} else if !actualField.IsEmbedded() && expectedField.isEmbeddedField {
			t.Errorf("field with name %s for struct %s is not embedded, but should be embedded field", expectedFieldName, structName)
		}
	}

	return true
}
