package visitor

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"strings"
	"testing"
)

type receiverInfo struct {
	name      string
	isPointer bool
	typeName  string
}

type functionInfo struct {
	markers    marker.MarkerValues
	isVariadic bool
	name       string
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
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Bread",
				},
			},
		},
		name:       "Bread",
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
		name:       "Macaron",
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
		markers: marker.MarkerValues{
			"marker:function-level": {
				FunctionLevel{
					Name: "MakeACake",
				},
			},
		},
		name:       "MakeACake",
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
		name:       "BiscuitCake",
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
		name:       "Funfetti",
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
		name:       "IceCream",
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
		name:       "CupCake",
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
		name:       "Tart",
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
		name:       "Donut",
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
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Pudding",
				},
			},
		},
		name:       "Pudding",
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
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Pie",
				},
			},
		},
		name:       "Pie",
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
		name:       "muffin",
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
		name: "Eat",
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
		markers: marker.MarkerValues{
			"marker:struct-method-level": {
				StructMethodLevel{
					Name: "Buy",
				},
			},
		},
		name: "Buy",
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
		markers: marker.MarkerValues{
			"marker:struct-method-level": {
				StructMethodLevel{
					Name: "FortuneCookie",
				},
			},
		},
		name: "FortuneCookie",
		receiver: &receiverInfo{
			name:      "c",
			isPointer: true,
			typeName:  "Cookie",
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
		name: "Oreo",
		receiver: &receiverInfo{
			name:      "c",
			isPointer: true,
			typeName:  "Cookie",
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

func assertFunctions(t *testing.T, descriptor string, actualMethods *Functions, expectedMethods map[string]functionInfo) bool {

	if actualMethods.Len() != len(expectedMethods) {
		t.Errorf("the number of the methods should be %d, but got %d", len(expectedMethods), actualMethods.Len())
		return false
	}

	for expectedMethodName, expectedMethod := range expectedMethods {
		actualMethod, ok := actualMethods.FindByName(expectedMethodName)

		if !ok {
			t.Errorf("method with name %s is not found for %s", expectedMethodName, descriptor)
			continue
		}

		if expectedMethod.String() != actualMethod.String() {
			t.Errorf("the signature of the function %s should be %s, but got %s", actualMethod.name, expectedMethod.String(), actualMethod.String())
		}

		if expectedMethod.isVariadic && !actualMethod.IsVariadic() {
			t.Errorf("the function %s should be a variadic function for %s", expectedMethodName, descriptor)
		} else if !expectedMethod.isVariadic && actualMethod.IsVariadic() {
			t.Errorf("the function %s should not be a variadic function for %s", expectedMethodName, descriptor)
		}

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
