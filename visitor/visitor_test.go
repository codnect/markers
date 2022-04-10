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
	markers            marker.MarkerValues
	numExplicitMethods int
	numMethods         int
	embeddedTypes      []string
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
				"MakeACake": {
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
				},
				"BiscuitCake": {
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
				},
			},
			interfaces: map[string]interfaceInfo{
				"BakeryShop": {
					markers: marker.MarkerValues{
						"marker:interface-type-level": {
							InterfaceTypeLevel{
								Name: "BakeryShop",
							},
						},
					},
					numExplicitMethods: 1,
					numMethods:         8,
					embeddedTypes:      []string{"Dessert"},
				},
				"Dessert": {
					markers: marker.MarkerValues{
						"marker:interface-type-level": {
							InterfaceTypeLevel{
								Name: "Dessert",
							},
						},
					},
					numExplicitMethods: 7,
					numMethods:         7,
				},
				"NewYearsEveCookie": {
					markers: marker.MarkerValues{
						"marker:interface-type-level": {
							InterfaceTypeLevel{
								Name: "NewYearsEveCookie",
							},
						},
					},
					numExplicitMethods: 1,
					numMethods:         1,
				},
				"SweetShop": {
					markers: marker.MarkerValues{
						"marker:interface-type-level": {
							InterfaceTypeLevel{
								Name: "SweetShop",
							},
						},
					},
					numExplicitMethods: 1,
					numMethods:         9,
					embeddedTypes:      []string{"NewYearsEveCookie", "Dessert"},
				},
			},
			structs: map[string]structInfo{
				"FriedCookie": {
					markers: marker.MarkerValues{
						"marker:struct-type-level": {
							StructTypeLevel{
								Name: "FriedCookie",
							},
						},
					},
				},
				"Cookie": {
					markers: marker.MarkerValues{
						"marker:struct-type-level": {
							StructTypeLevel{
								Name: "Cookie",
							},
						},
					},
				},
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

	err := EachFile(collector, result.GetPackages(), func(file *File, err error) error {
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

		if actualInterface.NumMethods() != expectedInterface.numMethods {
			t.Errorf("the number of the methods of the interface %s should be %d, but got %d", expectedInterfaceName, expectedInterface.numMethods, actualInterface.NumMethods())
			continue
		}

		if actualInterface.NumExplicitMethods() != expectedInterface.numExplicitMethods {
			t.Errorf("the number of the explicit methods of the interface %s should be %d, but got %d", expectedInterfaceName, expectedInterface.numExplicitMethods, actualInterface.NumExplicitMethods())
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

func assertFunctionParameters(t *testing.T, expectedParams []variableInfo, actualParams *Tuple, msg string) {
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
	}
}

func assertFunctionResult(t *testing.T, expectedResults []variableInfo, actualResults *Tuple, msg string) {
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
