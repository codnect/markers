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
			imports: []importInfo{
				{
					name:       "_",
					path:       "fmt",
					sideEffect: true,
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

		if !assertImports(t, file, testCase.imports) {
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
