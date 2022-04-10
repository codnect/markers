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
			markers            marker.MarkerValues
			numExplicitMethods int
			numMethods         int
			embeddedTypes      []string
		}
		structs map[string]struct {
			markers marker.MarkerValues
		}
	}{
		"dessert.go": {
			interfaces: map[string]struct {
				markers            marker.MarkerValues
				numExplicitMethods int
				numMethods         int
				embeddedTypes      []string
			}{
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

		if !assertInterfaces(t, file, testCase.interfaces) {
			return nil
		}

		if !assertStructs(t, file, testCase.structs) {
			return nil
		}

		return nil
	})

	if err != nil {
		t.Errorf("traverval is not completed successfully")
	}
}

func assertInterfaces(t *testing.T, file *File, interfaces map[string]struct {
	markers            marker.MarkerValues
	numExplicitMethods int
	numMethods         int
	embeddedTypes      []string
}) bool {

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

		if actualInterface.Markers().Count() != expectedInterface.markers.Count() {
			t.Errorf("the number of the interface %s markers should be %d, but got %d", expectedInterfaceName, expectedInterface.markers.Count(), actualInterface.Markers().Count())
			continue
		}

		for markerName, markerValues := range expectedInterface.markers {
			if actualInterface.Markers().CountByName(markerName) != len(markerValues) {
				t.Errorf("the number of the marker %s should be %d, but got %d", markerName, len(markerValues), actualInterface.Markers().CountByName(markerName))
				continue
			}

			actualMarkerValues := actualInterface.Markers().AllMarkers(markerName)

			for index, expectedMarkerValue := range markerValues {
				actualMarker := actualMarkerValues[index]
				assert.Equal(t, expectedMarkerValue, actualMarker)
			}
		}
	}

	return true
}

func assertStructs(t *testing.T, file *File, structs map[string]struct {
	markers marker.MarkerValues
}) bool {

	if len(structs) != file.Structs().Len() {
		t.Errorf("structs count is wrong!")
		return false
	}

	for expectedStructName, expectedStruct := range structs {
		actualStruct, ok := file.Structs().FindByName(expectedStructName)
		if !ok {
			t.Errorf("struct with name %s is not found", expectedStruct)
			continue
		}

		if actualStruct.Markers().Count() != expectedStruct.markers.Count() {
			t.Errorf("the number of the struct %s markers should be %d, but got %d", expectedStructName, expectedStruct.markers.Count(), actualStruct.Markers().Count())
			continue
		}

		for markerName, markerValues := range expectedStruct.markers {
			if actualStruct.Markers().CountByName(markerName) != len(markerValues) {
				t.Errorf("the number of the marker %s should be %d, but got %d", markerName, len(markerValues), actualStruct.Markers().CountByName(markerName))
				continue
			}

			actualMarkerValues := actualStruct.Markers().AllMarkers(markerName)

			for index, expectedMarkerValue := range markerValues {
				actualMarker := actualMarkerValues[index]
				assert.Equal(t, expectedMarkerValue, actualMarker)
			}
		}
	}

	return true
}
