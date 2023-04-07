package visitor

import (
	"fmt"
	"github.com/procyon-projects/markers"
	"github.com/procyon-projects/markers/packages"
	"github.com/stretchr/testify/assert"
	"testing"
)

type PackageLevel struct {
	Name string `marker:"Name"`
}

type StructTypeLevel struct {
	Name string `marker:"Name"`
	Any  any    `marker:"Any"`
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
	name      string
	typeName  string
	isPointer bool
}

func (v variableInfo) String() string {
	if v.isPointer {
		return fmt.Sprintf("*%s", v.typeName)
	}

	return v.typeName
}

func TestVisitor_VisitPackage(t *testing.T) {
	markerList := []struct {
		Name   string
		Level  markers.TargetLevel
		Output interface{}
	}{
		{Name: "marker:package-level", Level: markers.PackageLevel, Output: &PackageLevel{}},
		{Name: "marker:interface-type-level", Level: markers.InterfaceTypeLevel, Output: &InterfaceTypeLevel{}},
		{Name: "marker:interface-method-level", Level: markers.InterfaceMethodLevel, Output: &InterfaceMethodLevel{}},
		{Name: "marker:function-level", Level: markers.FunctionLevel, Output: &FunctionLevel{}},
		{Name: "marker:struct-type-level", Level: markers.StructTypeLevel, Output: &StructTypeLevel{}},
		{Name: "marker:struct-method-level", Level: markers.StructMethodLevel, Output: &StructMethodLevel{}},
		{Name: "marker:struct-field-level", Level: markers.FieldLevel, Output: &StructFieldLevel{}},
	}

	testCasePkgs := map[string]map[string]testFile{
		"github.com/procyon-projects/markers/test/menu": {
			"coffee.go": {
				constants:   coffeeConstants,
				customTypes: coffeeCustomTypes,
				functions: map[string]functionInfo{
					"PrintCookie": printCookieMethod,
				},
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
						file:       "dessert.go",
						position:   Position{Line: 7, Column: 2},
					},
					{
						name:       "_",
						path:       "strings",
						sideEffect: true,
						file:       "dessert.go",
						position:   Position{Line: 8, Column: 2},
					},
				},
				functions: map[string]functionInfo{
					"MakeACake":     makeACakeFunction,
					"BiscuitCake":   biscuitCakeFunction,
					"Eat":           eatMethod,
					"Buy":           buyMethod,
					"FortuneCookie": fortuneCookieMethod,
					"Oreo":          oreoMethod,
				},
				interfaces: map[string]interfaceInfo{
					"BakeryShop":        bakeryShopInterface,
					"Dessert":           dessertInterface,
					"newYearsEveCookie": newYearsEveCookieInterface,
					"SweetShop":         sweetShopInterface,
				},
				structs: map[string]structInfo{
					"FriedCookie": friedCookieStruct,
					"cookie":      cookieStruct,
				},
			},
		},
		"github.com/procyon-projects/markers/test/any": {
			"error.go": {
				constants:   []constantInfo{},
				customTypes: errorCustomTypes,
				functions: map[string]functionInfo{
					"Print":    printErrorMethod,
					"ToErrors": toErrorsMethod,
				},
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
					"Index":           indexMethod,
				},
				interfaces: map[string]interfaceInfo{
					"Repository":     repositoryInterface,
					"Number":         numberInterface,
					"EventPublisher": eventPublisherInterface,
				},
				structs: map[string]structInfo{
					"Controller":     controllerStruct,
					"TestController": testControllerStruct,
				},
				imports: []importInfo{
					{
						name:       "",
						path:       "context",
						file:       "generics.go",
						sideEffect: false,
						position:   Position{Line: 4, Column: 2},
					},
					{
						name:       "",
						path:       "golang.org/x/exp/constraints",
						file:       "generics.go",
						sideEffect: false,
						position:   Position{Line: 5, Column: 2},
					},
				},
				customTypes: genericsCustomTypes,
			},
			"method.go": {
				functions: map[string]functionInfo{
					"Print": printHttpHandlerMethod,
				},
			},
			"string.go": {
				imports: []importInfo{
					{
						name:       "",
						path:       "net/http",
						sideEffect: false,
						file:       "string.go",
						position:   Position{Line: 3, Column: 8},
					},
				},
				constants: stringConstants,
			},
		},
	}

	result, _ := packages.LoadPackages("../test/...")
	registry := markers.NewRegistry()

	for _, m := range markerList {
		err := registry.Register(m.Name, "github.com/procyon-projects/markers", m.Level, m.Output)
		if err != nil {
			t.Errorf("marker %s could not be registered", m.Name)
			return
		}
	}

	collector := markers.NewCollector(registry)

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

func assertMarkers(t *testing.T, expectedMarkers markers.Values, actualMarkers markers.Values, msg string) {
	if actualMarkers.Count() != expectedMarkers.Count() {
		t.Errorf("the number of the markers of the %s should be %d, but got %d", msg, expectedMarkers.Count(), actualMarkers.Count())
		return
	}

	for markerName, markerValues := range expectedMarkers {
		if actualMarkers.CountByName(markerName) != len(markerValues) {
			t.Errorf("%s: the number of the marker %s should be %d, but got %d", msg, markerName, len(markerValues), actualMarkers.CountByName(markerName))
			continue
		}

		actualMarkerValues, _ := actualMarkers.FindByName(markerName)

		for index, expectedMarkerValue := range markerValues {
			actualMarker := actualMarkerValues[index]
			assert.Equal(t, expectedMarkerValue, actualMarker, "%s", msg)
		}
	}
}
