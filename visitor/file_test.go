package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testFile struct {
	path          string
	constants     []constantInfo
	interfaces    map[string]interfaceInfo
	structs       map[string]structInfo
	functions     map[string]functionInfo
	imports       []importInfo
	customTypes   map[string]customTypeInfo
	importMarkers []importMarkerInfo
	fileMarkers   []fileMarkerInfo
}

type importInfo struct {
	name       string
	path       string
	sideEffect bool
	file       string
	position   Position
}

type importMarkerInfo struct {
	value string
	pkg   string
	alias string
}

type fileMarkerInfo any

func sideEffects(imports []importInfo) []importInfo {
	result := make([]importInfo, 0)
	for _, importItem := range imports {
		if importItem.sideEffect {
			result = append(result, importItem)
		}
	}

	return result
}

func assertImports(t *testing.T, file *File, expectedImports []importInfo, expectedImportMarkers []importMarkerInfo, fileMarkers []fileMarkerInfo) bool {
	if file.Imports().Len() != len(expectedImports) {
		t.Errorf("the number of the imports in file %s should be %d, but got %d", file.name, len(expectedImports), file.Imports().Len())
	}

	assert.Equal(t, file.Imports().elements, file.Imports().ToSlice(), "ToSlice should return %w, but got %w", file.Imports().elements, file.Imports().ToSlice())

	if len(file.Imports().SideEffects()) != len(sideEffects(expectedImports)) {
		t.Errorf("the number of the side effect imports in file %s should be %d, but got %d", file.name, len(sideEffects(expectedImports)), len(file.Imports().SideEffects()))
	}

	for index, expectedImport := range expectedImports {
		fileImport := file.Imports().At(index)

		actualImport, exists := file.Imports().FindByName(expectedImport.name)
		if !expectedImport.sideEffect && (!exists || actualImport == nil) {
			t.Errorf("import with name %s in file %s  is not found", file.name, expectedImport.name)
			continue
		}

		actualImport, exists = file.Imports().FindByPath(expectedImport.path)
		if !exists || actualImport == nil {
			t.Errorf("import with name %s in file %s  is not found", file.name, expectedImport.name)
			continue
		}

		assert.Equal(t, fileImport, actualImport, "Imports.At should return %w, but got %w", fileImport, actualImport)

		if expectedImport.name != actualImport.Name() {
			t.Errorf("import name in file %s shoud be %s, but got %s", file.name, expectedImport.name, actualImport.Name())
		}

		if expectedImport.path != actualImport.Path() {
			t.Errorf("import path in file %s shoud be %s, but got %s", file.name, expectedImport.path, actualImport.Path())
		}

		if expectedImport.file != actualImport.File().Name() {
			t.Errorf("the file name for import '%s' should be %s, but got %s", expectedImport.path, expectedImport.file, actualImport.File().Name())
		}

		if actualImport.SideEffect() && !expectedImport.sideEffect {
			t.Errorf("import with path %s in file %s is not an import side effect, but should be an import side effect", expectedImport.path, file.name)
		} else if !actualImport.SideEffect() && expectedImport.sideEffect {
			t.Errorf("import with path %s in file %s is an import side effect, but should not be an import side effect", expectedImport.path, file.name)
		}

		assert.Equal(t, expectedImport.position, actualImport.Position(), "position for import with path %s in file %s should be %w, but got %w", expectedImport.name, "", expectedImport.position, fileImport.Position())
	}

	if file.Markers().Count() != len(fileMarkers) {
		t.Errorf("the number of the file markers in file %s should be %d, but got %d", file.name, len(fileMarkers), file.Markers().Count())
	}

	assertImportMarkers(t, file, expectedImportMarkers)

	return true
}

func assertImportMarkers(t *testing.T, file *File, expectedImportMarkers []importMarkerInfo) {

	if file.NumImportMarkers() != len(expectedImportMarkers) {
		t.Errorf("the number of the import markers in file %s should be %d, but got %d", file.name, len(expectedImportMarkers), len(file.ImportMarkers()))
	}

	for index, importMarker := range file.ImportMarkers() {
		expectedImportMarker := expectedImportMarkers[index]
		if importMarker.Pkg != expectedImportMarker.pkg {
			t.Errorf("the Pkg attribute of the import marker in file %s shoud be %s, but got %s", file.name, expectedImportMarker.pkg, importMarker.Pkg)
		}

		if importMarker.Value != expectedImportMarker.value {
			t.Errorf("the Value attribute of the import marker in file %s shoud be %s, but got %s", file.name, expectedImportMarker.value, importMarker.Value)
		}

		if importMarker.Alias != expectedImportMarker.alias {
			t.Errorf("the Alias attribute of the import marker in file %s shoud be %s, but got %s", file.name, expectedImportMarker.alias, importMarker.Alias)
		}
	}
}
