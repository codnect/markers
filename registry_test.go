package marker

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
		err := registry.Register(testCase.MarkerName, testCase.TargetLevel, testCase.Output)
		assert.Nil(t, err)

		definition, ok := registry.definitionMap[testCase.MarkerName]
		if !ok {
			t.Error("marker has not been registered successfully")
		}

		if definition.Name != testCase.MarkerName {
			t.Errorf("marker name is not equal to expected, got %q; want %q", definition.Name, testCase.MarkerName)
		}

		if definition.Level != testCase.TargetLevel {
			t.Errorf("target level is not equal to expected, got %q; want %q", definition.Level, testCase.TargetLevel)
		}
	}

	assert.Len(t, registry.definitionMap, len(testCases))
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
		newDefinition, _ := MakeDefinition(testCase.MarkerName, testCase.TargetLevel, testCase.Output)
		err := registry.RegisterWithDefinition(newDefinition)
		assert.Nil(t, err)

		definition, ok := registry.definitionMap[testCase.MarkerName]
		if !ok {
			t.Error("marker has not been registered successfully")
		}

		if definition.Name != testCase.MarkerName {
			t.Errorf("marker name is not equal to expected, got %q; want %q", definition.Name, testCase.MarkerName)
		}

		if definition.Level != testCase.TargetLevel {
			t.Errorf("target level is not equal to expected, got %q; want %q", definition.Level, testCase.TargetLevel)
		}
	}

	assert.Len(t, registry.definitionMap, len(testCases))
}

func TestRegistry_RegisterMarkerAlreadyRegistered(t *testing.T) {
	registry := NewRegistry()
	registry.Register("marker:test", TypeLevel, &testTypeLevelMarker{})
	err := registry.Register("marker:test", MethodLevel, &testTypeLevelMarker{})
	assert.Len(t, registry.definitionMap, 1)
	assert.NotNil(t, err)
	assert.Equal(t, "there is already registered definition : marker:test", err.Error())
}

func TestRegistry_RegisterMarkerWithEmptyName(t *testing.T) {
	registry := NewRegistry()
	err := registry.Register("", MethodLevel, &testMarker{})
	assert.Len(t, registry.definitionMap, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "marker name cannot be empty", err.Error())
}

func TestRegistry_RegisterMarkerWithoutLevel(t *testing.T) {
	registry := NewRegistry()
	err := registry.Register("marker:test", 0, &testMarker{})
	assert.Len(t, registry.definitionMap, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "specify target levels for the definition : marker:test", err.Error())
}
