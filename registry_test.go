package markers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testMarker struct {
}

type testTypeLevelMarker struct {
}

type testFunctionLevelMarker struct {
}

type testMethodFunctionLevelMarker struct {
}

type testStructInterfaceTypeLevelMarker struct {
}

func TestRegistry_Register(t *testing.T) {
	testCases := []struct {
		MarkerName  string
		TargetLevel TargetLevel
		Output      interface{}
	}{
		{"marker:type-level", TypeLevel, &testTypeLevelMarker{}},
		{"marker:function-level", FunctionLevel, &testFunctionLevelMarker{}},
		{"marker:method-function-level", MethodLevel | FunctionLevel, &testMethodFunctionLevelMarker{}},
		{"marker:struct-interface-level", StructTypeLevel | InterfaceTypeLevel, &testStructInterfaceTypeLevelMarker{}},
	}

	registry := NewRegistry()

	for _, testCase := range testCases {
		err := registry.Register(testCase.MarkerName, "anyPkg", testCase.TargetLevel, testCase.Output)
		assert.Nil(t, err)

		definition, ok := registry.packageMap["anyPkg"][testCase.MarkerName]
		if !ok {
			t.Error("marker has not been registered successfully")
		}

		if definition.Name != testCase.MarkerName {
			t.Errorf("marker name is not equal to expected, got %q; want %q", definition.Name, testCase.MarkerName)
		}

		if definition.TargetLevel != testCase.TargetLevel {
			t.Errorf("target level is not equal to expected, got %q; want %q", definition.TargetLevel, testCase.TargetLevel)
		}
	}

	assert.Len(t, registry.packageMap["anyPkg"], len(testCases))
}

func TestRegistry_RegisterWithDefinition(t *testing.T) {
	testCases := []struct {
		MarkerName  string
		TargetLevel TargetLevel
		Output      interface{}
	}{
		{"marker:type-level", TypeLevel, &testTypeLevelMarker{}},
		{"marker:function-level", FunctionLevel, &testFunctionLevelMarker{}},
		{"marker:method-function-level", MethodLevel | FunctionLevel, &testMethodFunctionLevelMarker{}},
		{"marker:struct-interface-level", StructTypeLevel | InterfaceTypeLevel, &testStructInterfaceTypeLevelMarker{}},
	}

	registry := NewRegistry()

	for _, testCase := range testCases {
		newDefinition, _ := MakeDefinition(testCase.MarkerName, "anyPkg", testCase.TargetLevel, testCase.Output)
		err := registry.RegisterWithDefinition(newDefinition)
		assert.Nil(t, err)

		definition, ok := registry.packageMap["anyPkg"][testCase.MarkerName]
		if !ok {
			t.Error("marker has not been registered successfully")
		}

		if definition.Name != testCase.MarkerName {
			t.Errorf("marker name is not equal to expected, got %q; want %q", definition.Name, testCase.MarkerName)
		}

		if definition.TargetLevel != testCase.TargetLevel {
			t.Errorf("target level is not equal to expected, got %q; want %q", definition.TargetLevel, testCase.TargetLevel)
		}
	}

	assert.Len(t, registry.packageMap["anyPkg"], len(testCases))
}

func TestRegistry_RegisterMarkerAlreadyRegistered(t *testing.T) {
	registry := NewRegistry()
	registry.Register("marker:test", "anyPkg", TypeLevel, &testTypeLevelMarker{})
	err := registry.Register("marker:test", "anyPkg", MethodLevel, &testTypeLevelMarker{})
	assert.Len(t, registry.packageMap["anyPkg"], 1)
	assert.NotNil(t, err)
	assert.Equal(t, "there is already registered definition : marker:test", err.Error())
}

func TestRegistry_RegisterMarkerWithEmptyName(t *testing.T) {
	registry := NewRegistry()
	err := registry.Register("", "anyPkg", MethodLevel, &testMarker{})
	assert.Len(t, registry.packageMap["anyPkg"], 0)
	assert.NotNil(t, err)
	assert.Equal(t, "marker name cannot be empty", err.Error())
}

func TestRegistry_RegisterMarkerWithoutLevel(t *testing.T) {
	registry := NewRegistry()
	err := registry.Register("marker:test", "anyPkg", 0, &testMarker{})
	assert.Len(t, registry.packageMap["anyPkg"], 0)
	assert.NotNil(t, err)
	assert.Equal(t, "specify target levels for the definition: marker:test", err.Error())
}
