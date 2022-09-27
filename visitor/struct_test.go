package visitor

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"testing"
)

type fieldInfo struct {
	name            string
	typeName        string
	isExported      bool
	isEmbeddedField bool
}

type structInfo struct {
	markers           marker.MarkerValues
	methods           map[string]functionInfo
	fields            map[string]fieldInfo
	numFields         int
	numAllFields      int
	numEmbeddedFields int
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
		methods: map[string]functionInfo{
			"Eat": eatMethod,
			"Buy": buyMethod,
		},
		fields: map[string]fieldInfo{
			"Cookie": {
				isExported:      true,
				isEmbeddedField: true,
			},
			"cookieDough": {
				isExported:      false,
				isEmbeddedField: false,
			},
		},
		numFields:         2,
		numAllFields:      3,
		numEmbeddedFields: 1,
	}

	cookieStruct = structInfo{
		markers: marker.MarkerValues{
			"marker:struct-type-level": {
				StructTypeLevel{
					Name: "Cookie",
				},
			},
		},
		methods: map[string]functionInfo{
			"FortuneCookie": fortuneCookieMethod,
			"Oreo":          oreoMethod,
		},
		fields: map[string]fieldInfo{
			"ChocolateChip": {
				isExported:      true,
				isEmbeddedField: false,
			},
		},
		numFields:         2,
		numAllFields:      2,
		numEmbeddedFields: 0,
	}
)

func assertStructs(t *testing.T, file *File, structs map[string]structInfo) bool {

	if len(structs) != file.Structs().Len() {
		t.Errorf("the number of the functions should be %d, but got %d", len(structs), file.Structs().Len())
		return false
	}

	for expectedStructName, expectedStruct := range structs {
		actualStruct, ok := file.Structs().FindByName(expectedStructName)
		if !ok {
			t.Errorf("struct with name %s is not found", expectedStructName)
			continue
		}

		if actualStruct.NumAllMethods() != len(expectedStruct.methods) {
			t.Errorf("the number of the methods of the struct %s should be %d, but got %d", expectedStructName, len(expectedStruct.methods), actualStruct.NumMethods())
		}

		if actualStruct.NumMethods() != len(expectedStruct.methods) {
			t.Errorf("the number of the methods of the struct %s should be %d, but got %d", expectedStructName, len(expectedStruct.methods), actualStruct.NumMethods())
		}

		if actualStruct.NumFields() != expectedStruct.numFields {
			t.Errorf("the number of the fields of the struct %s should be %d, but got %d", expectedStructName, expectedStruct.numAllFields-expectedStruct.numEmbeddedFields, actualStruct.NumFields())
		}

		if actualStruct.NumAllFields() != expectedStruct.numAllFields {
			t.Errorf("the number of the all fields of the struct %s should be %d, but got %d", expectedStructName, expectedStruct.numAllFields, actualStruct.NumFields())
		}

		if actualStruct.NumEmbeddedFields() != expectedStruct.numEmbeddedFields {
			t.Errorf("the number of the embededed fields of the struct %s should be %d, but got %d", expectedStructName, expectedStruct.numEmbeddedFields, actualStruct.NumFields())
		}

		assertFunctions(t, actualStruct.Methods(), expectedStruct.methods)
		assertStructFields(t, actualStruct.Name(), actualStruct.Fields(), expectedStruct.fields)
		assertMarkers(t, expectedStruct.markers, actualStruct.Markers(), fmt.Sprintf("struct %s", expectedStructName))
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

		if actualField.IsExported() && !expectedField.isExported {
			t.Errorf("field with name %s for struct %s is exported, but should be unexported field", expectedFieldName, structName)
		} else if !actualField.IsExported() && expectedField.isExported {
			t.Errorf("field with name %s for struct %s is not exported, but should be exported field", expectedFieldName, structName)
		}

		if actualField.IsEmbedded() && !expectedField.isExported {
			t.Errorf("field with name %s for struct %s is embedded, but should be not embedded field", expectedFieldName, structName)
		} else if !actualField.IsEmbedded() && expectedField.isEmbeddedField {
			t.Errorf("field with name %s for struct %s is not embedded, but should be embedded field", expectedFieldName, structName)
		}
	}

	return true
}
