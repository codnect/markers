package visitor

import (
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

	testCases := map[string]struct {
		interfaces map[string]struct {
			markers marker.MarkerValues
		}
		structs map[string]struct {
			markers marker.MarkerValues
		}
	}{
		"dessert.go": {
			interfaces: map[string]struct {
				markers marker.MarkerValues
			}{
				"Dessert": {
					markers: marker.MarkerValues{
						"marker:interface-type-level": {
							InterfaceTypeLevel{
								Name: "Dessert",
							},
						},
					},
				},
			},
			structs: map[string]struct {
				markers marker.MarkerValues
			}{
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

		if len(testCase.interfaces) != file.Interfaces().Len() {
			t.Errorf("interface count is wrong!")
			return nil
		}

		for expectedInterfaceName, expectedInterface := range testCase.interfaces {
			actualInterface, ok := file.Interfaces().FindByName(expectedInterfaceName)
			if !ok {
				t.Errorf("interface with name %s is not found", expectedInterfaceName)
				continue
			}

			if actualInterface.Markers().Count() != expectedInterface.markers.Count() {
				t.Errorf("marker count is wrong!")
				continue
			}

			for markerName, markerValues := range expectedInterface.markers {
				if actualInterface.Markers().CountByName(markerName) != len(markerValues) {
					t.Errorf("marker count is wrong!")
					continue
				}

				actualMarkerValues := actualInterface.Markers().AllMarkers(markerName)

				for index, expectedMarkerValue := range markerValues {
					actualMarker := actualMarkerValues[index]
					assert.Equal(t, expectedMarkerValue, actualMarker)
				}
			}
		}

		if len(testCase.structs) != file.Structs().Len() {
			t.Errorf("structs count is wrong!")
			return nil
		}

		for expectedStructName, expectedStruct := range testCase.structs {
			actualStruct, ok := file.Structs().FindByName(expectedStructName)
			if !ok {
				t.Errorf("struct with name %s is not found", expectedStruct)
				continue
			}

			if actualStruct.Markers().Count() != expectedStruct.markers.Count() {
				t.Errorf("marker count is wrong!")
				continue
			}

			for markerName, markerValues := range expectedStruct.markers {
				if actualStruct.Markers().CountByName(markerName) != len(markerValues) {
					t.Errorf("marker count is wrong!")
					continue
				}

				actualMarkerValues := actualStruct.Markers().AllMarkers(markerName)

				for index, expectedMarkerValue := range markerValues {
					actualMarker := actualMarkerValues[index]
					assert.Equal(t, expectedMarkerValue, actualMarker)
				}
			}
		}

		return nil
	})

	if err != nil {
		t.Errorf("traverval is not completed successfully")
	}
}
