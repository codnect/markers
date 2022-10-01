package visitor

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"testing"
)

type functionInfo struct {
	markers    marker.MarkerValues
	isVariadic bool
	params     []variableInfo
	results    []variableInfo
}

// functions
var (
	breadFunction = functionInfo{
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Bread",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Macaron",
				},
			},
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
				typeName: "book",
			},
		},
	}

	makeACakeFunction = functionInfo{
		markers: marker.MarkerValues{
			"marker:function-level": {
				FunctionLevel{
					Name: "MakeACake",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:function-level": {
				FunctionLevel{
					Name: "BiscuitCake",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Funfetti",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "IceCream",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "CupCake",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Tart",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Donut",
				},
			},
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

	puddingFunction = functionInfo{
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Pudding",
				},
			},
		},
		isVariadic: false,
		params:     []variableInfo{},
		results: []variableInfo{
			{
				name:     "",
				typeName: "[]string",
			},
		},
	}

	pieFunction = functionInfo{
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Pie",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "muffin",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:struct-method-level": {
				StructMethodLevel{
					Name: "Eat",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:struct-method-level": {
				StructMethodLevel{
					Name: "Buy",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:struct-method-level": {
				StructMethodLevel{
					Name: "FortuneCookie",
				},
			},
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
		markers: marker.MarkerValues{
			"marker:struct-method-level": {
				StructMethodLevel{
					Name: "Oreo",
				},
			},
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
)

func assertFunctions(t *testing.T, descriptior string, actualMethods *Functions, expectedMethods map[string]functionInfo) bool {

	if actualMethods.Len() != len(expectedMethods) {
		t.Errorf("the number of the methods should be %d, but got %d", len(expectedMethods), actualMethods.Len())
		return false
	}

	for expectedMethodName, expectedMethod := range expectedMethods {
		actualMethod, ok := actualMethods.FindByName(expectedMethodName)

		if !ok {
			t.Errorf("method with name %s is not found for %s", expectedMethodName, descriptior)
			continue
		}

		if expectedMethod.isVariadic && !actualMethod.IsVariadic() {
			t.Errorf("the function %s should be a variadic function for %s", expectedMethodName, descriptior)
		} else if !expectedMethod.isVariadic && actualMethod.IsVariadic() {
			t.Errorf("the function %s should not be a variadic function for %s", expectedMethodName, descriptior)
		}

		assertFunctionParameters(t, expectedMethod.params, actualMethod.Params(), fmt.Sprintf("function %s (%s)", expectedMethodName, descriptior))

		assertFunctionResult(t, expectedMethod.results, actualMethod.Results(), fmt.Sprintf("function %s (%s)", expectedMethodName, descriptior))

		assertMarkers(t, expectedMethod.markers, actualMethod.markers, fmt.Sprintf("function %s (%s)", expectedMethodName, descriptior))
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
