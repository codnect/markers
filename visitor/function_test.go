package visitor

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type receiverInfo struct {
	name      string
	isPointer bool
	typeName  string
}

type functionInfo struct {
	markers    markers.MarkerValues
	isVariadic bool
	name       string
	fileName   string
	position   Position
	receiver   *receiverInfo
	params     []variableInfo
	results    []variableInfo
}

func (f functionInfo) String() string {
	var builder strings.Builder
	builder.WriteString("func ")

	if f.receiver != nil {
		builder.WriteString("(")
		builder.WriteString(f.receiver.name)
		builder.WriteString(" ")
		if f.receiver.isPointer {
			builder.WriteString("*")
		}
		builder.WriteString(f.receiver.typeName)
		builder.WriteString(") ")
	}

	builder.WriteString(f.name)
	builder.WriteString("(")

	if len(f.params) != 0 {
		for i := 0; i < len(f.params); i++ {
			param := f.params[i]
			if param.name != "" {
				builder.WriteString(param.name + " ")
			}

			builder.WriteString(param.typeName)

			if i != len(f.params)-1 {
				builder.WriteString(",")
			}
		}
	}

	builder.WriteString(") ")

	if len(f.results) > 1 {
		builder.WriteString("(")
	}

	if len(f.results) != 0 {
		for i := 0; i < len(f.results); i++ {
			result := f.results[i]
			if result.name != "" {
				builder.WriteString(result.name + " ")
			}

			builder.WriteString(result.typeName)

			if i != len(f.results)-1 {
				builder.WriteString(",")
			}
		}
	}

	if len(f.results) > 1 {
		builder.WriteString(")")
	}

	return builder.String()
}

// functions
var (
	breadFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Bread",
				},
			},
		},
		name:     "Bread",
		fileName: "dessert.go",
		position: Position{
			Line:   16,
			Column: 7,
		},
		isVariadic: false,
		params: []variableInfo{
			{
				name:     "i",
				typeName: "float64",
			},
			{
				name:     "k",
				typeName: "float64",
			},
		},
		results: []variableInfo{
			{
				name:     "",
				typeName: "struct{}",
			},
		},
	}

	macaronFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Macaron",
				},
			},
		},
		name:     "Macaron",
		fileName: "dessert.go",
		position: Position{
			Line:   133,
			Column: 9,
		},
		isVariadic: false,
		params: []variableInfo{
			{
				name:     "c",
				typeName: "complex128",
			},
		},
		results: []variableInfo{
			{
				name:     "",
				typeName: "chan string",
			},
			{
				name:     "",
				typeName: "fmt.Stringer",
			},
		},
	}

	makeACakeFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:function-level": {
				FunctionLevel{
					Name: "MakeACake",
				},
			},
		},
		name:     "MakeACake",
		fileName: "dessert.go",
		position: Position{
			Line:   113,
			Column: 1,
		},
		isVariadic: false,
		params: []variableInfo{
			{
				name:     "s",
				typeName: "interface{}",
			},
		},
		results: []variableInfo{
			{
				name:     "",
				typeName: "error",
			},
		},
	}

	biscuitCakeFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:function-level": {
				FunctionLevel{
					Name: "BiscuitCake",
				},
			},
		},
		name:     "BiscuitCake",
		fileName: "dessert.go",
		position: Position{
			Line:   119,
			Column: 1,
		},
		isVariadic: true,
		params: []variableInfo{
			{
				name:     "s",
				typeName: "string",
			},
			{
				name:     "arr",
				typeName: "[]int",
			},
			{
				name:     "v",
				typeName: "int16",
			},
		},
		results: []variableInfo{
			{
				name:     "i",
				typeName: "int",
			},
			{
				name:     "b",
				typeName: "bool",
			},
		},
	}

	funfettiFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Funfetti",
				},
			},
		},
		name:     "Funfetti",
		fileName: "dessert.go",
		position: Position{
			Line:   51,
			Column: 10,
		},
		isVariadic: false,
		params: []variableInfo{
			{
				name:     "v",
				typeName: "rune",
			},
		},
		results: []variableInfo{
			{
				name:     "",
				typeName: "byte",
			},
		},
	}

	iceCreamFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "IceCream",
				},
			},
		},
		name:     "IceCream",
		fileName: "dessert.go",
		position: Position{
			Line:   84,
			Column: 10,
		},
		isVariadic: true,
		params: []variableInfo{
			{
				name:     "s",
				typeName: "string",
			},
			{
				name:     "v",
				typeName: "bool",
			},
		},
		results: []variableInfo{
			{
				name:     "r",
				typeName: "string",
			},
		},
	}

	cupCakeFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "CupCake",
				},
			},
		},
		name:     "CupCake",
		fileName: "dessert.go",
		position: Position{
			Line:   88,
			Column: 9,
		},
		isVariadic: false,
		params: []variableInfo{
			{
				name:     "a",
				typeName: "[]int",
			},
			{
				name:     "b",
				typeName: "bool",
			},
		},
		results: []variableInfo{
			{
				name:     "",
				typeName: "float32",
			},
		},
	}

	tartFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Tart",
				},
			},
		},
		name:     "Tart",
		fileName: "dessert.go",
		position: Position{
			Line:   92,
			Column: 6,
		},
		isVariadic: false,
		params: []variableInfo{
			{
				name:     "s",
				typeName: "interface{}",
			},
		},
		results: []variableInfo{},
	}

	donutFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Donut",
				},
			},
		},
		name:     "Donut",
		fileName: "dessert.go",
		position: Position{
			Line:   96,
			Column: 7,
		},
		isVariadic: false,
		params:     []variableInfo{},
		results: []variableInfo{
			{
				name:     "",
				typeName: "error",
			},
		},
	}

	puddingFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Pudding",
				},
			},
		},
		name:     "Pudding",
		fileName: "dessert.go",
		position: Position{
			Line:   100,
			Column: 9,
		},
		isVariadic: false,
		params:     []variableInfo{},
		results: []variableInfo{
			{
				name:     "",
				typeName: "[5]string",
			},
		},
	}

	pieFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Pie",
				},
			},
		},
		name:     "Pie",
		fileName: "dessert.go",
		position: Position{
			Line:   104,
			Column: 5,
		},
		isVariadic: false,
		params:     []variableInfo{},
		results: []variableInfo{
			{
				name:     "",
				typeName: "interface{}",
			},
		},
	}

	muffinFunction = functionInfo{
		markers: markers.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "muffin",
				},
			},
		},
		name:     "muffin",
		fileName: "dessert.go",
		position: Position{
			Line:   108,
			Column: 8,
		},
		isVariadic: false,
		params:     []variableInfo{},
		results: []variableInfo{
			{
				name:     "",
				typeName: "string",
			},
			{
				name:     "",
				typeName: "error",
			},
		},
	}

	eatMethod = functionInfo{
		markers: markers.MarkerValues{
			"marker:struct-method-level": {
				StructMethodLevel{
					Name: "Eat",
				},
			},
		},
		name:     "Eat",
		fileName: "dessert.go",
		position: Position{
			Line:   24,
			Column: 1,
		},
		receiver: &receiverInfo{
			name:      "c",
			isPointer: true,
			typeName:  "FriedCookie",
		},
		isVariadic: false,
		params:     []variableInfo{},
		results: []variableInfo{
			{
				name:     "",
				typeName: "bool",
			},
		},
	}

	buyMethod = functionInfo{
		markers: markers.MarkerValues{
			"marker:struct-method-level": {
				StructMethodLevel{
					Name: "Buy",
				},
			},
		},
		name:     "Buy",
		fileName: "dessert.go",
		position: Position{
			Line:   42,
			Column: 1,
		},
		receiver: &receiverInfo{
			name:      "c",
			isPointer: true,
			typeName:  "FriedCookie",
		},
		isVariadic: false,
		params: []variableInfo{
			{
				name:     "i",
				typeName: "int",
			},
		},
		results: []variableInfo{},
	}

	fortuneCookieMethod = functionInfo{
		markers: markers.MarkerValues{
			"marker:struct-method-level": {
				StructMethodLevel{
					Name: "FortuneCookie",
				},
			},
		},
		name:     "FortuneCookie",
		fileName: "dessert.go",
		position: Position{
			Line:   67,
			Column: 1,
		},
		receiver: &receiverInfo{
			name:      "c",
			isPointer: true,
			typeName:  "cookie",
		},
		isVariadic: false,
		params: []variableInfo{
			{
				name:     "v",
				typeName: "interface{}",
			},
		},
		results: []variableInfo{
			{
				name:     "",
				typeName: "[]string",
			},
		},
	}

	oreoMethod = functionInfo{
		markers: markers.MarkerValues{
			"marker:struct-method-level": {
				StructMethodLevel{
					Name: "Oreo",
				},
			},
		},
		name:     "Oreo",
		fileName: "dessert.go",
		position: Position{
			Line:   73,
			Column: 1,
		},
		receiver: &receiverInfo{
			name:      "c",
			isPointer: true,
			typeName:  "cookie",
		},
		isVariadic: true,
		params: []variableInfo{
			{
				name:     "a",
				typeName: "[]interface{}",
			},
			{
				name:     "v",
				typeName: "string",
			},
		},
		results: []variableInfo{
			{
				name:     "",
				typeName: "error",
			},
		},
	}

	genericFunction = functionInfo{
		markers:  markers.MarkerValues{},
		name:     "GenericFunction",
		fileName: "generics.go",
		position: Position{
			Line:   3,
			Column: 1,
		},
		isVariadic: false,
		params:     []variableInfo{},
		results:    []variableInfo{},
	}
)

func assertFunctions(t *testing.T, descriptor string, actualMethods *Functions, expectedMethods map[string]functionInfo) bool {

	if actualMethods.Len() != len(expectedMethods) {
		t.Errorf("the number of the methods of %s should be %d, but got %d", descriptor, len(expectedMethods), actualMethods.Len())
		return false
	}

	for expectedMethodName, expectedMethod := range expectedMethods {
		actualMethod, ok := actualMethods.FindByName(expectedMethodName)

		if !ok {
			t.Errorf("method with name %s is not found for %s", expectedMethodName, descriptor)
			continue
		}

		if expectedMethodName != actualMethod.Name() {
			t.Errorf("the name of the function should be %s, but got %s", expectedMethodName, actualMethod.Name())
		}

		if expectedMethod.fileName != actualMethod.File().Name() {
			t.Errorf("the file name for function %s should be %s, but got %s", expectedMethodName, expectedMethod.fileName, actualMethod.File().Name())
		}

		if expectedMethod.String() != actualMethod.String() {
			t.Errorf("the signature of the function %s should be %s, but got %s", actualMethod.name, expectedMethod.String(), actualMethod.String())
		}

		if expectedMethod.isVariadic && !actualMethod.IsVariadic() {
			t.Errorf("the function %s should be a variadic function for %s", expectedMethodName, descriptor)
		} else if !expectedMethod.isVariadic && actualMethod.IsVariadic() {
			t.Errorf("the function %s should not be a variadic function for %s", expectedMethodName, descriptor)
		}

		typeParam := actualMethod.TypeParams()
		if typeParam != nil {
			typeParam.Len()
		}

		assert.Equal(t, actualMethod, actualMethod.Underlying())

		assert.Equal(t, expectedMethod.position, actualMethod.Position(), "the position of the function %s for %s should be %w, but got %w",
			expectedMethodName, descriptor, expectedMethod.position, actualMethod.Position())

		assertFunctionParameters(t, expectedMethod.params, actualMethod.Params(), fmt.Sprintf("function %s (%s)", expectedMethodName, descriptor))

		assertFunctionResult(t, expectedMethod.results, actualMethod.Results(), fmt.Sprintf("function %s (%s)", expectedMethodName, descriptor))

		assertMarkers(t, expectedMethod.markers, actualMethod.markers, fmt.Sprintf("function %s (%s)", expectedMethodName, descriptor))
	}

	return true
}

func assertFunctionParameters(t *testing.T, expectedParams []variableInfo, actualParams Variables, msg string) {
	if actualParams.Len() != len(expectedParams) {
		t.Errorf("the number of the %s parameters should be %d, but got %d", msg, len(expectedParams), actualParams.Len())
		return
	}

	for index := 0; index < actualParams.Len(); index++ {
		actualFunctionParam := actualParams.At(index)
		expectedFunctionParam := expectedParams[index]

		if expectedFunctionParam.name != actualFunctionParam.Name() {
			t.Errorf("at index %d, the parameter name of the %s should be %s, but got %s", index, msg, expectedFunctionParam.name, actualFunctionParam.name)
		}

		if expectedFunctionParam.typeName != actualFunctionParam.Type().Name() {
			t.Errorf("at index %d, the parameter type name of the %s should be %s, but got %s", index, msg, expectedFunctionParam.typeName, actualFunctionParam.Type().Name())
		}
	}
}

func assertFunctionResult(t *testing.T, expectedResults []variableInfo, actualResults Variables, msg string) {
	if actualResults.Len() != len(expectedResults) {
		t.Errorf("the number of the %s results should be %d, but got %d", msg, len(expectedResults), actualResults.Len())
		return
	}

	for index := 0; index < actualResults.Len(); index++ {
		actualFunctionParam := actualResults.At(index)
		expectedFunctionParam := expectedResults[index]

		if expectedFunctionParam.name != actualFunctionParam.Name() {
			t.Errorf("at index %d, the parameter result of the %s should be %s, but got %s", index, msg, expectedFunctionParam.name, actualFunctionParam.name)
		}
	}

}
