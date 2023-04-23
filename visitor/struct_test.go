package visitor

import (
	"fmt"
	"github.com/procyon-projects/markers"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fieldInfo struct {
	name            string
	typeName        string
	stringValue     string
	isExported      bool
	isEmbeddedField bool
	markers         markers.Values
}

type structInfo struct {
	fileName               string
	isExported             bool
	position               Position
	markers                markers.Values
	methods                map[string]functionInfo
	allMethods             map[string]functionInfo
	fields                 map[string]fieldInfo
	embeddedFields         map[string]fieldInfo
	numFields              int
	totalFields            int
	numEmbeddedFields      int
	stringValue            string
	isAnonymous            bool
	implementsInterfaces   []string
	noImplementsInterfaces []string
}

// structs
var (
	controllerStruct = structInfo{
		markers:     markers.Values{},
		stringValue: "any.Controller[C context.Context,T any|int,Y ~int]",
		fileName:    "generics.go",
		isExported:  true,
		position: Position{
			Line:   17,
			Column: 6,
		},
		methods: map[string]functionInfo{
			"Index": indexMethod,
		},
		allMethods: map[string]functionInfo{
			"Index": indexMethod,
		},
		fields: map[string]fieldInfo{
			"AnyField1": {
				isExported:      true,
				isEmbeddedField: false,
				typeName:        "string",
				stringValue:     "string",
			},
			"AnyField2": {
				isExported:      true,
				isEmbeddedField: false,
				typeName:        "int",
				stringValue:     "int",
			},
		},
		embeddedFields:    map[string]fieldInfo{},
		numFields:         2,
		totalFields:       2,
		numEmbeddedFields: 0,
	}
	testControllerStruct = structInfo{
		markers:     markers.Values{},
		stringValue: "any.TestController",
		fileName:    "generics.go",
		isExported:  true,
		position: Position{
			Line:   26,
			Column: 6,
		},
		methods: map[string]functionInfo{},
		allMethods: map[string]functionInfo{
			"Index": indexMethod,
		},
		fields: map[string]fieldInfo{
			"Controller": {
				isExported:      true,
				isEmbeddedField: true,
				typeName:        "Controller",
				stringValue:     "Controller[context.Context,int16,int]",
			},
			"BaseController": {
				isExported:      true,
				isEmbeddedField: true,
				typeName:        "BaseController",
				stringValue:     "BaseController[int]",
			},
		},
		embeddedFields: map[string]fieldInfo{
			"Controller": {
				isExported:      true,
				isEmbeddedField: true,
				typeName:        "Controller",
				stringValue:     "Controller[context.Context,int16,int]",
			},
			"BaseController": {
				isExported:      true,
				isEmbeddedField: true,
				typeName:        "BaseController",
				stringValue:     "BaseController[int]",
			},
		},
		numFields:         2,
		totalFields:       2,
		numEmbeddedFields: 2,
	}

	friedCookieStruct = structInfo{
		markers: markers.Values{
			"test-marker:struct-type-level": {
				StructTypeLevel{
					Name: "FriedCookie",
				},
			},
		},
		stringValue: "menu.FriedCookie",
		fileName:    "dessert.go",
		isExported:  true,
		position: Position{
			Line:   31,
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
			"PrintCookie":   printCookieMethod,
		},
		fields: map[string]fieldInfo{
			"cookie": {
				isExported:      false,
				isEmbeddedField: true,
				typeName:        "cookie",
				stringValue:     "menu.cookie",
			},
			"cookieDough": {
				isExported:      false,
				isEmbeddedField: false,
				typeName:        "any",
				stringValue:     "any",
				markers: markers.Values{
					"test-marker:struct-field-level": {
						StructFieldLevel{
							Name: "CookieDough",
						},
					},
				},
			},
			"anonymousStruct": {
				isExported:      false,
				isEmbeddedField: false,
				typeName:        "struct{}",
				stringValue:     "struct{}",
			},
			"emptyInterface": {
				isExported:      false,
				isEmbeddedField: false,
				typeName:        "interface{}",
				stringValue:     "interface{}",
			},
		},
		embeddedFields: map[string]fieldInfo{
			"cookie": {
				isExported:      false,
				isEmbeddedField: true,
				typeName:        "cookie",
				stringValue:     "menu.cookie",
			},
		},
		implementsInterfaces:   []string{"Meal"},
		noImplementsInterfaces: []string{"Dessert"},
		numFields:              4,
		totalFields:            5,
		numEmbeddedFields:      1,
	}

	cookieStruct = structInfo{
		markers: markers.Values{
			"test-marker:struct-type-level": {
				StructTypeLevel{
					Name: "cookie",
					Any: map[string]interface{}{
						"key": "value",
					},
				},
			},
		},
		stringValue: "menu.cookie",
		fileName:    "dessert.go",
		isExported:  false,
		position: Position{
			Line:   61,
			Column: 6,
		},
		methods: map[string]functionInfo{
			"FortuneCookie": fortuneCookieMethod,
			"Oreo":          oreoMethod,
			"PrintCookie":   printCookieMethod,
		},
		allMethods: map[string]functionInfo{
			"FortuneCookie": fortuneCookieMethod,
			"Oreo":          oreoMethod,
			"PrintCookie":   printCookieMethod,
		},
		fields: map[string]fieldInfo{
			"ChocolateChip": {
				isExported:      true,
				isEmbeddedField: false,
				typeName:        "string",
				stringValue:     "string",
				markers: markers.Values{
					"test-marker:struct-field-level": {
						StructFieldLevel{
							Name: "ChocolateChip",
						},
					},
				},
			},
			"tripleChocolateCookie": {
				isExported:      false,
				isEmbeddedField: false,
				typeName:        "map[string]error",
				stringValue:     "map[string]error",
				markers: markers.Values{
					"test-marker:struct-field-level": {
						StructFieldLevel{
							Name: "tripleChocolateCookie",
						},
					},
				},
			},
		},
		embeddedFields:    map[string]fieldInfo{},
		numFields:         2,
		totalFields:       2,
		numEmbeddedFields: 0,
	}
	baseControllerStruct = structInfo{
		markers:     markers.Values{},
		stringValue: "any.BaseController[M any]",
		fileName:    "generics.go",
		isExported:  true,
		position: Position{
			Line:   42,
			Column: 6,
		},
		methods:           map[string]functionInfo{},
		allMethods:        map[string]functionInfo{},
		fields:            map[string]fieldInfo{},
		embeddedFields:    map[string]fieldInfo{},
		numFields:         0,
		totalFields:       0,
		numEmbeddedFields: 0,
	}
)

func assertStructs(t *testing.T, file *File, structs map[string]structInfo) bool {

	if len(structs) != file.Structs().Len() {
		t.Errorf("the number of the structs in file %s should be %d, but got %d", file.Name(), len(structs), file.Structs().Len())
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

		if actualStruct.IsAnonymous() && !expectedStruct.isAnonymous {
			t.Errorf("struct with name %s is anonymous, but should be anonymous", actualStruct.Name())
		} else if !actualStruct.IsAnonymous() && expectedStruct.isAnonymous {
			t.Errorf("struct with name %s is not anonymous, but should be anonymous", actualStruct.Name())
		}

		if actualStruct.NumFields() == 0 && !actualStruct.IsEmpty() {
			t.Errorf("the struct %s should be empty", actualStruct.Name())
		} else if actualStruct.NumFields() != 0 && actualStruct.IsEmpty() {
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
			t.Errorf("the number of the all fields of the struct %s should be %d, but got %d", expectedStructName, expectedStruct.totalFields, actualStruct.NumFieldsInHierarchy())
		}

		if actualStruct.NumEmbeddedFields() != expectedStruct.numEmbeddedFields {
			t.Errorf("the number of the embededed fields of the struct %s should be %d, but got %d", expectedStructName, expectedStruct.numEmbeddedFields, actualStruct.NumFields())
		}

		if expectedStruct.stringValue != actualStruct.String() {
			t.Errorf("Output returning from String() method for struct type with name %s does not equal to %s, but got %s", expectedStructName, expectedStruct.stringValue, actualStruct.String())
		}

		assert.Equal(t, actualStruct, actualStruct.Underlying())

		assert.Equal(t, expectedStruct.position, actualStruct.Position(), "the position of the struct %s should be %w, but got %w",
			expectedStructName, expectedStruct.position, actualStruct.Position())

		assertFunctions(t, fmt.Sprintf("struct %s", actualStruct.Name()), actualStruct.Methods(), expectedStruct.methods)
		assertFunctions(t, fmt.Sprintf("struct %s", actualStruct.Name()), actualStruct.MethodsInHierarchy(), expectedStruct.allMethods)
		assertStructFields(t, actualStruct.Name(), actualStruct.EmbeddedFields(), expectedStruct.embeddedFields)
		assertStructFields(t, actualStruct.Name(), actualStruct.Fields(), expectedStruct.fields)
		assertMarkers(t, expectedStruct.markers, actualStruct.Markers(), fmt.Sprintf("struct %s", expectedStructName))

		for _, interfaceName := range expectedStruct.implementsInterfaces {
			iface, exists := file.Interfaces().FindByName(interfaceName)

			if !exists {
				t.Errorf("the interface %s should exists in file %s and the struct %s should implement it", interfaceName, file.Name(), actualStruct.Name())
				continue
			}

			if !actualStruct.Implements(iface) {
				t.Errorf(" the struct %s should implement the interface %s", actualStruct.Name(), interfaceName)
				continue
			}
		}

		for _, interfaceName := range expectedStruct.noImplementsInterfaces {
			iface, exists := file.Interfaces().FindByName(interfaceName)

			if !exists {
				t.Errorf("the interface %s should exists in file %s and the struct %s should implement it", interfaceName, file.Name(), actualStruct.Name())
				continue
			}

			if actualStruct.Implements(iface) {
				t.Errorf(" the struct %s should not implement the interface %s", actualStruct.Name(), interfaceName)
				continue
			}
		}
		index++
	}

	return true
}

func assertStructFields(t *testing.T, structName string, actualFields *Fields, expectedFields map[string]fieldInfo) bool {
	if actualFields.Len() != len(expectedFields) {
		t.Errorf("the number of the fields of struct %s should be %d, but got %d", structName, len(expectedFields), actualFields.Len())
		return false
	}

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
			t.Errorf("type of field with name %s for struct %s shoud be '%s', but got %s", actualField.Name(), structName, expectedField.typeName, actualField.Type().Name())
		}

		if actualField.Type().String() != expectedField.stringValue {
			t.Errorf("String() result shoud be '%s' for the field with name %s in struct %s, but got %s", expectedField.stringValue, actualField.Name(), structName, actualField.Type().String())
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

		assertMarkers(t, expectedField.markers, actualField.Markers(), fmt.Sprintf("field %s in struct %s", expectedFieldName, structName))
	}

	return true
}

func TestStructs_AtShouldReturnNilIfIndexIsOutOfRange(t *testing.T) {
	structs := &Structs{}
	assert.Nil(t, structs.At(0))
}

func TestFields_AtShouldReturnNilIfIndexIsOutOfRange(t *testing.T) {
	fields := &Fields{}
	assert.Nil(t, fields.At(0))
}
func TestFields_AtShouldReturnFieldIfIndexIsBetweenRange(t *testing.T) {
	fields := &Fields{
		elements: []*Field{
			{
				name: "anyField",
			},
		},
	}
	field := fields.At(0)
	assert.NotNil(t, field)
	assert.Equal(t, "anyField", field.Name())
}
