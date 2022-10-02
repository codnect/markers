package visitor

import "testing"

type testFile struct {
	interfaces map[string]interfaceInfo
	structs    map[string]structInfo
	functions  map[string]functionInfo
	imports    []importInfo
}

type importInfo struct {
	name       string
	path       string
	sideEffect bool
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

	if len(file.Imports().SideEffects()) != len(sideEffects(expectedImports)) {
		t.Errorf("the number of the side effect imports in file %s should be %d, but got %d", file.name, len(sideEffects(expectedImports)), len(file.Imports().SideEffects()))
	}

	return true
}
