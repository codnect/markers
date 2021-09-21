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
	registry := NewRegistry()
	registry.Register("marker:type-level", TypeLevel, &testTypeLevelMarker{})
	registry.Register("marker:function-level", FunctionLevel, &testFunctionLevelMarker{})
	registry.Register("marker:method-function-level", MethodLevel|FunctionLevel, &testMethodFunctionLevelMarker{})
	registry.Register("marker:struct-interface-level", StructTypeLevel|InterfaceTypeLevel, &testStructInterfaceTypeLevelMarker{})

	assert.Len(t, registry.definitionMap, 4)

	definition := registry.definitionMap["marker:type-level"]
	assert.NotNil(t, definition)
	assert.Equal(t, definition.Name, "marker:type-level")
	assert.Equal(t, definition.Level, TypeLevel)

	definition = registry.definitionMap["marker:function-level"]
	assert.NotNil(t, definition)
	assert.Equal(t, definition.Name, "marker:function-level")
	assert.Equal(t, definition.Level, FunctionLevel)

	definition = registry.definitionMap["marker:method-function-level"]
	assert.NotNil(t, definition)
	assert.Equal(t, definition.Name, "marker:method-function-level")
	assert.Equal(t, definition.Level, MethodLevel|FunctionLevel)

	definition = registry.definitionMap["marker:struct-interface-level"]
	assert.NotNil(t, definition)
	assert.Equal(t, definition.Name, "marker:struct-interface-level")
	assert.Equal(t, definition.Level, StructTypeLevel|InterfaceTypeLevel)
}

func TestRegistry_RegisterWithDefinition(t *testing.T) {
	registry := NewRegistry()
	definition, _ := MakeDefinition("marker:type-level", TypeLevel, &testTypeLevelMarker{})
	registry.RegisterWithDefinition(definition)
	definition, _ = MakeDefinition("marker:function-level", FunctionLevel, &testFunctionLevelMarker{})
	registry.RegisterWithDefinition(definition)
	definition, _ = MakeDefinition("marker:method-function-level", MethodLevel|FunctionLevel, &testMethodFunctionLevelMarker{})
	registry.RegisterWithDefinition(definition)
	definition, _ = MakeDefinition("marker:struct-interface-level", StructTypeLevel|InterfaceTypeLevel, &testStructInterfaceTypeLevelMarker{})
	registry.RegisterWithDefinition(definition)

	assert.Len(t, registry.definitionMap, 4)

	definition = registry.definitionMap["marker:type-level"]
	assert.NotNil(t, definition)
	assert.Equal(t, definition.Name, "marker:type-level")
	assert.Equal(t, definition.Level, TypeLevel)

	definition = registry.definitionMap["marker:function-level"]
	assert.NotNil(t, definition)
	assert.Equal(t, definition.Name, "marker:function-level")
	assert.Equal(t, definition.Level, FunctionLevel)

	definition = registry.definitionMap["marker:method-function-level"]
	assert.NotNil(t, definition)
	assert.Equal(t, definition.Name, "marker:method-function-level")
	assert.Equal(t, definition.Level, MethodLevel|FunctionLevel)

	definition = registry.definitionMap["marker:struct-interface-level"]
	assert.NotNil(t, definition)
	assert.Equal(t, definition.Name, "marker:struct-interface-level")
	assert.Equal(t, definition.Level, StructTypeLevel|InterfaceTypeLevel)
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
