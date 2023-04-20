package visitor

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/markers"
	"github.com/procyon-projects/markers/packages"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"runtime"
	"strings"
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

func TestEachFile_ShouldReturnErrorIfCollectorIsNil(t *testing.T) {
	err := EachFile(nil, nil, nil)
	assert.Equal(t, "collector cannot be nil", err.Error())
}

func TestEachFile_ShouldReturnErrorIfPkgsIsNil(t *testing.T) {
	err := EachFile(&markers.Collector{}, nil, nil)
	assert.Equal(t, "packages cannot be nil", err.Error())
}

func TestEachFile_ShouldReturnErrorIfTraversedPkgIsNil(t *testing.T) {
	err := EachFile(&markers.Collector{}, []*packages.Package{nil}, nil)
	assert.Equal(t, markers.ErrorList{errors.New("pkg(package) cannot be nil")}, err.(markers.ErrorList))
}

func TestVisitor_VisitPackage(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	path := filepath.Dir(file)
	lastSlashIndex := strings.LastIndex(path, "/")
	path = path[:lastSlashIndex]

	markerList := []struct {
		Name   string
		Level  markers.TargetLevel
		Output interface{}
	}{
		{Name: "test-marker:package-level", Level: markers.PackageLevel, Output: &PackageLevel{}},
		{Name: "test-marker:interface-type-level", Level: markers.InterfaceTypeLevel, Output: &InterfaceTypeLevel{}},
		{Name: "test-marker:interface-method-level", Level: markers.InterfaceMethodLevel, Output: &InterfaceMethodLevel{}},
		{Name: "test-marker:function-level", Level: markers.FunctionLevel, Output: &FunctionLevel{}},
		{Name: "test-marker:struct-type-level", Level: markers.StructTypeLevel, Output: &StructTypeLevel{}},
		{Name: "test-marker:struct-method-level", Level: markers.StructMethodLevel, Output: &StructMethodLevel{}},
		{Name: "test-marker:struct-field-level", Level: markers.FieldLevel, Output: &StructFieldLevel{}},
	}

	testCasePkgs := map[string]map[string]testFile{
		"github.com/procyon-projects/markers/test/menu": {
			"coffee.go": {
				path:        fmt.Sprintf("%s/test/menu/coffee.go", path),
				constants:   coffeeConstants,
				customTypes: coffeeCustomTypes,
				functions: map[string]functionInfo{
					"PrintCookie": printCookieMethod,
				},
				importMarkers: []importMarkerInfo{
					{
						pkg:   "github.com/procyon-projects/markers",
						value: "marker",
					},
					{
						pkg:   "github.com/procyon-projects/test-markers",
						value: "test-marker",
					},
				},
				fileMarkers: []fileMarkerInfo{
					PackageLevel{
						Name: "coffee.go",
					},
				},
			},
			"fresh.go": {
				path:        fmt.Sprintf("%s/test/menu/fresh.go", path),
				constants:   freshConstants,
				customTypes: freshCustomTypes,
				importMarkers: []importMarkerInfo{
					{
						pkg:   "github.com/procyon-projects/markers",
						value: "marker",
					},
					{
						pkg:   "github.com/procyon-projects/test-markers",
						value: "test-marker",
					},
				},
				fileMarkers: []fileMarkerInfo{
					PackageLevel{
						Name: "fresh.go",
					},
				},
			},
			"dessert.go": {
				path: fmt.Sprintf("%s/test/menu/dessert.go", path),
				imports: []importInfo{
					{
						name:       "",
						path:       "fmt",
						sideEffect: false,
						file:       "dessert.go",
						position:   Position{Line: 8, Column: 2},
					},
					{
						name:       "_",
						path:       "strings",
						sideEffect: true,
						file:       "dessert.go",
						position:   Position{Line: 9, Column: 2},
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
				importMarkers: []importMarkerInfo{
					{
						pkg:   "github.com/procyon-projects/markers",
						value: "marker",
					},
					{
						pkg:   "github.com/procyon-projects/test-markers",
						value: "test-marker",
					},
				},
				fileMarkers: []fileMarkerInfo{
					PackageLevel{
						Name: "dessert.go",
					},
				},
			},
		},
		"github.com/procyon-projects/markers/test/any": {
			"error.go": {
				path:      fmt.Sprintf("%s/test/any/error.go", path),
				constants: []constantInfo{},
				functions: map[string]functionInfo{
					"Print":    printErrorMethod,
					"ToErrors": toErrorsMethod,
				},
			},
			"other.go": {
				path:        fmt.Sprintf("%s/test/any/other.go", path),
				constants:   []constantInfo{},
				customTypes: errorCustomTypes,
			},
			"permission.go": {
				path:        fmt.Sprintf("%s/test/any/permission.go", path),
				constants:   permissionConstants,
				customTypes: permissionCustomTypes,
				importMarkers: []importMarkerInfo{
					{
						pkg:   "github.com/procyon-projects/markers",
						value: "marker",
						alias: "test",
					},
					{
						pkg:   "github.com/procyon-projects/test-markers",
						value: "test-marker",
					},
				},
				fileMarkers: []fileMarkerInfo{
					PackageLevel{
						Name: "permission.go",
					},
				},
			},
			"math.go": {
				path:      fmt.Sprintf("%s/test/any/math.go", path),
				constants: mathConstants,
			},
			"generics.go": {
				path:      fmt.Sprintf("%s/test/any/generics.go", path),
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
				path: fmt.Sprintf("%s/test/any/method.go", path),
				functions: map[string]functionInfo{
					"Print": printHttpHandlerMethod,
				},
			},
			"string.go": {
				path: fmt.Sprintf("%s/test/any/string.go", path),
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
		err := registry.Register(m.Name, "github.com/procyon-projects/test-markers", m.Level, m.Output)
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

		if testCase.path != file.Path() {
			t.Errorf("file path %s shoud be %s, but got %s", file.name, testCase.path, file.Path())
		}

		if !assertImports(t, file, testCase.imports, testCase.importMarkers, testCase.fileMarkers) {
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
