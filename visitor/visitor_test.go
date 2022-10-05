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

type variableInfo struct {
	name     string
	typeName string
}

func TestVisitor_VisitPackage(t *testing.T) {
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

	testCasePkgs := map[string]map[string]testFile{
		"github.com/procyon-projects/marker/test/menu": {
			"coffee.go": {
				constants:   coffeeConstants,
				customTypes: coffeeCustomTypes,
			},
			"fresh.go": {
				constants:   freshConstants,
				customTypes: freshCustomTypes,
			},
			"dessert.go": {
				imports: []importInfo{
					{
						name:       "",
						path:       "fmt",
						sideEffect: false,
						position:   Position{Line: 7, Column: 2},
					},
					{
						name:       "_",
						path:       "strings",
						sideEffect: true,
						position:   Position{Line: 8, Column: 2},
					},
				},
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
		},
		"github.com/procyon-projects/marker/test/any": {
			"error.go": {
				constants:   []constantInfo{},
				customTypes: errorCustomTypes,
			},
			"permission.go": {
				constants:   permissionConstants,
				customTypes: permissionCustomTypes,
			},
			"math.go": {
				constants: mathConstants,
			},
			"generics.go": {
				constants: []constantInfo{},
				functions: map[string]functionInfo{
					"GenericFunction": genericFunction,
				},
			},
			"string.go": {
				imports: []importInfo{
					{
						name:       "",
						path:       "net/http",
						sideEffect: false,
						position:   Position{Line: 3, Column: 8},
					},
				},
				constants: stringConstants,
			},
		},
	}

	result, _ := packages.LoadPackages("../test/...")
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
		if _, isTestCasePkg := testCasePkgs[file.pkg.ID]; !isTestCasePkg {
			return nil
		}

		testCase, exists := testCasePkgs[file.pkg.ID][file.Name()]

		if !exists {
			t.Errorf("file %s not found in test cases", file.Name())
			return nil
		}

		if !assertImports(t, file, testCase.imports) {
			return nil
		}

		if !assertConstants(t, file, testCase.constants) {
			return nil
		}

		if !assertCustomTypes(t, file, testCase.customTypes) {
			return nil
		}

		if !assertInterfaces(t, file, testCase.interfaces) {
			return nil
		}

		if !assertStructs(t, file, testCase.structs) {
			return nil
		}

		if !assertFunctions(t, fmt.Sprintf("file %s", file.Name()), file.Functions(), testCase.functions) {
			return nil
		}

		return nil
	})

	if err != nil {
		t.Errorf("traverval is not completed successfully")
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
