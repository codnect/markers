package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testFile struct {
	constants  map[string]struct{}
	interfaces map[string]interfaceInfo
	structs    map[string]structInfo
	functions  map[string]functionInfo
	imports    []importInfo
}

type importInfo struct {
	name       string
	path       string
	sideEffect bool
	position   Position
}

func sideEffects(imports []importInfo) []importInfo {
	result := make([]importInfo, 0)
	for _, importItem := range imports {
		if importItem.sideEffect {
			result = append(result, importItem)
		}
	}

	return result
}

func assertImports(t *testing.T, file *File, expectedImports []importInfo) bool {
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

		if actualImport.SideEffect() && !expectedImport.sideEffect {
			t.Errorf("import with path %s in file %s is not an import side effect, but should be an import side effect", expectedImport.path, file.name)
		} else if !actualImport.SideEffect() && expectedImport.sideEffect {
			t.Errorf("import with path %s in file %s is an import side effect, but should not be an import side effect", expectedImport.path, file.name)
		}

		assert.Equal(t, expectedImport.position, actualImport.Position(), "position for import with path %s in file %s should be %w, but got %w", expectedImport.position, fileImport.Position())
	}

	return true
}
