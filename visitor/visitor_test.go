package visitor

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"github.com/stretchr/testify/assert"
	"testing"
)

type PackageLevel struct {
	Name string `marker:"Name"`
}

type StructTypeLevel struct {
	Name string `marker:"Name"`
}

type StructMethodLevel struct {
	Name string `marker:"Name"`
}

type StructFieldLevel struct {
	Name string `marker:"Name"`
}

type InterfaceTypeLevel struct {
	Name string `marker:"Name"`
}

type InterfaceMethodLevel struct {
	Name string `marker:"Name"`
}

type FunctionLevel struct {
	Name string `marker:"Name"`
}

type testFile struct {
	interfaces map[string]interfaceInfo
	structs    map[string]structInfo
	functions  map[string]functionInfo
}

type interfaceInfo struct {
	markers         marker.MarkerValues
	explicitMethods map[string]functionInfo
	methods         map[string]functionInfo
	embeddedTypes   []string
}

type structInfo struct {
	markers marker.MarkerValues
}

type functionInfo struct {
	markers    marker.MarkerValues
	isVariadic bool
	params     []variableInfo
	results    []variableInfo
}

type variableInfo struct {
	name     string
	typeName string
}

// functions
var (
	breadFunction = functionInfo{
		markers: marker.MarkerValues{
			"marker:function-level": {
				FunctionLevel{
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
			"marker:function-level": {
				FunctionLevel{
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
		isVariadic: true,
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
			{
				name:     "b",
				typeName: "bool",
			},
		},
		results: []variableInfo{},
	}

	donutFunction = functionInfo{
		markers: marker.MarkerValues{
			"marker:interface-method-level": {
				InterfaceMethodLevel{
					Name: "Tart",
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
				typeName: "interface{}",
			},
		},
	}
)

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
	}

	cookieStruct = structInfo{
		markers: marker.MarkerValues{
			"marker:struct-type-level": {
				StructTypeLevel{
					Name: "Cookie",
				},
			},
		},
	}
)

// interfaces
var (
	bakeryShopInterface = interfaceInfo{
		markers: marker.MarkerValues{
			"marker:interface-type-level": {
				InterfaceTypeLevel{
					Name: "BakeryShop",
				},
			},
		},
		explicitMethods: map[string]functionInfo{
			"Bread": breadFunction,
		},
		methods: map[string]functionInfo{
			"IceCream": iceCreamFunction,
			"CupCake":  cupCakeFunction,
			"Tart":     tartFunction,
			"Donut":    donutFunction,
			"Pudding":  puddingFunction,
			"Pie":      pieFunction,
			"muffin":   muffinFunction,
			"Bread":    breadFunction,
		},
		embeddedTypes: []string{"Dessert"},
	}

	dessertInterface = interfaceInfo{
		markers: marker.MarkerValues{
			"marker:interface-type-level": {
				InterfaceTypeLevel{
					Name: "Dessert",
				},
			},
		},
		explicitMethods: map[string]functionInfo{
			"IceCream": iceCreamFunction,
			"CupCake":  cupCakeFunction,
			"Tart":     tartFunction,
			"Donut":    donutFunction,
			"Pudding":  puddingFunction,
			"Pie":      pieFunction,
			"muffin":   muffinFunction,
		},
		methods: map[string]functionInfo{
			"IceCream": iceCreamFunction,
			"CupCake":  cupCakeFunction,
			"Tart":     tartFunction,
			"Donut":    donutFunction,
			"Pudding":  puddingFunction,
			"Pie":      pieFunction,
			"muffin":   muffinFunction,
		},
	}

	newYearsEveCookieInterface = interfaceInfo{
		markers: marker.MarkerValues{
			"marker:interface-type-level": {
				InterfaceTypeLevel{
					Name: "NewYearsEveCookie",
				},
			},
		},
		methods: map[string]functionInfo{
			"Funfetti": funfettiFunction,
		},
		explicitMethods: map[string]functionInfo{
			"Funfetti": funfettiFunction,
		},
	}

	sweetShopInterface = interfaceInfo{
		markers: marker.MarkerValues{
			"marker:interface-type-level": {
				InterfaceTypeLevel{
					Name: "SweetShop",
				},
			},
		},
		explicitMethods: map[string]functionInfo{
			"Macaron": macaronFunction,
		},
		methods: map[string]functionInfo{
			"Funfetti": funfettiFunction,
			"Macaron":  macaronFunction,
			"IceCream": iceCreamFunction,
			"CupCake":  cupCakeFunction,
			"Tart":     tartFunction,
			"Donut":    donutFunction,
			"Pudding":  puddingFunction,
			"Pie":      pieFunction,
			"muffin":   muffinFunction,
		},
		embeddedTypes: []string{"NewYearsEveCookie", "Dessert"},
	}
)

func TestVisitor_VisitPackage1(t *testing.T) {
	markers := []struct {
		Name   string
		Level  marker.TargetLevel
		Output interface{}
	}{
		{Name: "marker:package-level", Level: marker.PackageLevel, Output: &PackageLevel{}},
		{Name: "marker:interface-type-level", Level: marker.InterfaceTypeLevel, Output: &InterfaceTypeLevel{}},
		{Name: "marker:interface-method-level", Level: marker.InterfaceMethodLevel, Output: &InterfaceMethodLevel{}},
		{Name: "marker:function-level", Level: marker.FunctionLevel, Output: &FunctionLevel{}},
		{Name: "marker:struct-type-level", Level: marker.StructTypeLevel, Output: &StructTypeLevel{}},
		{Name: "marker:struct-method-level", Level: marker.StructMethodLevel, Output: &StructMethodLevel{}},
		{Name: "marker:struct-field-level", Level: marker.FieldLevel, Output: &StructFieldLevel{}},
	}

	testCases := map[string]testFile{
		"dessert.go": {
			functions: map[string]functionInfo{
				"MakeACake":   makeACakeFunction,
				"BiscuitCake": biscuitCakeFunction,
			},
			interfaces: map[string]interfaceInfo{
				"BakeryShop":        bakeryShopInterface,
				"Dessert":           dessertInterface,
				"NewYearsEveCookie": newYearsEveCookieInterface,
				"SweetShop":         sweetShopInterface,
			},
			structs: map[string]structInfo{
				"FriedCookie": friedCookieStruct,
				"Cookie":      cookieStruct,
			},
		},
	}

	result, _ := packages.LoadPackages("../test/package1")
	registry := marker.NewRegistry()

	for _, m := range markers {
		err := registry.Register(m.Name, "github.com/procyon-projects/marker", m.Level, m.Output)
		if err != nil {
			t.Errorf("marker %s could not be registered", m.Name)
			return
		}
	}

	collector := marker.NewCollector(registry)

	err := EachFile(collector, result.Packages(), func(file *File, err error) error {
		if file.pkg.ID == "builtin" {
			return nil
		}

		testCase := testCases[file.Name()]

		if !assertInterfaces(t, file, testCase.interfaces) {
			return nil
		}

		if !assertStructs(t, file, testCase.structs) {
			return nil
		}

		if !assertFunctions(t, file, testCase.functions) {
			return nil
		}

		return nil
	})

	if err != nil {
		t.Errorf("traverval is not completed successfully")
	}
}

func assertInterfaces(t *testing.T, file *File, interfaces map[string]interfaceInfo) bool {

	if len(interfaces) != file.Interfaces().Len() {
		t.Errorf("the number of the interface should be %d, but got %d", len(interfaces), file.Interfaces().Len())
		return false
	}

	for expectedInterfaceName, expectedInterface := range interfaces {
		actualInterface, ok := file.Interfaces().FindByName(expectedInterfaceName)

		if !ok {
			t.Errorf("interface with name %s is not found", expectedInterfaceName)
			continue
		}

		if actualInterface.NumMethods() != len(expectedInterface.methods) {
			t.Errorf("the number of the methods of the interface %s should be %d, but got %d", expectedInterfaceName, len(expectedInterface.methods), actualInterface.NumMethods())
			continue
		}

		if actualInterface.NumExplicitMethods() != len(expectedInterface.explicitMethods) {
			t.Errorf("the number of the explicit methods of the interface %s should be %d, but got %d", expectedInterfaceName, len(expectedInterface.explicitMethods), actualInterface.NumExplicitMethods())
			continue
		}

		if actualInterface.NumEmbeddedTypes() != len(expectedInterface.embeddedTypes) {
			t.Errorf("the number of the embedded types of the interface %s should be %d, but got %d", expectedInterfaceName, len(expectedInterface.embeddedTypes), actualInterface.NumEmbeddedTypes())
			continue
		}

		for index, expectedEmbeddedType := range expectedInterface.embeddedTypes {
			actualEmbeddedType := actualInterface.EmbeddedTypes()[index]
			if expectedEmbeddedType != actualEmbeddedType.Name() {
				t.Errorf("the interface %s should have the embedded type %s at index %d, but got %s", expectedInterfaceName, expectedEmbeddedType, index, actualEmbeddedType.Name())
				continue
			}
		}

		assertMarkers(t, expectedInterface.markers, actualInterface.Markers(), fmt.Sprintf("interface %s", expectedInterfaceName))

	}

	return true
}

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

		assertMarkers(t, expectedStruct.markers, actualStruct.Markers(), fmt.Sprintf("struct %s", expectedStructName))
	}

	return true
}

func assertFunctions(t *testing.T, file *File, functions map[string]functionInfo) bool {

	if len(functions) != file.Functions().Len() {
		t.Errorf("the number of the functions should be %d, but got %d", len(functions), file.Functions().Len())
		return false
	}

	for expectedFunctionName, expectedFunction := range functions {
		actualFunction, ok := file.Functions().FindByName(expectedFunctionName)

		if !ok {
			t.Errorf("function with name %s is not found", expectedFunctionName)
			continue
		}

		if actualFunction.Receiver() != nil {
			t.Errorf("the receiver of the function %s should be nil", expectedFunctionName)
			continue
		}

		if expectedFunction.isVariadic && !actualFunction.IsVariadic() {
			t.Errorf("the function %s should be a variadic function", expectedFunctionName)
			continue
		} else if !expectedFunction.isVariadic && actualFunction.IsVariadic() {
			t.Errorf("the function %s should not be a variadic function", expectedFunctionName)
			continue
		}

		assertFunctionParameters(t, expectedFunction.params, actualFunction.Params(), fmt.Sprintf("function %s", expectedFunctionName))

		assertFunctionResult(t, expectedFunction.results, actualFunction.Results(), fmt.Sprintf("function %s", expectedFunctionName))

		assertMarkers(t, expectedFunction.markers, actualFunction.markers, fmt.Sprintf("function %s", expectedFunctionName))
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

func assertMarkers(t *testing.T, expectedMarkers marker.MarkerValues, actualMarkers marker.MarkerValues, msg string) {
	if actualMarkers.Count() != expectedMarkers.Count() {
		t.Errorf("the number of the markers of the %s should be %d, but got %d", msg, expectedMarkers.Count(), actualMarkers.Count())
		return
	}

	for markerName, markerValues := range expectedMarkers {
		if actualMarkers.CountByName(markerName) != len(markerValues) {
			t.Errorf("%s: the number of the marker %s should be %d, but got %d", msg, markerName, len(markerValues), actualMarkers.CountByName(markerName))
			continue
		}

		actualMarkerValues := actualMarkers.AllMarkers(markerName)

		for index, expectedMarkerValue := range markerValues {
			actualMarker := actualMarkerValues[index]
			assert.Equal(t, expectedMarkerValue, actualMarker, "%s", msg)
		}
	}
}
