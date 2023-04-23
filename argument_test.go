package markers

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestExtractArgument(t *testing.T) {
	reflMarker := reflect.TypeOf(Marker{})

	arg, err := extractArgument(reflMarker.Field(0))
	assert.Nil(t, err)
	assert.Equal(t, "Value", arg.Name)
	assert.Empty(t, arg.Default)
	assert.True(t, arg.Required)
	assert.False(t, arg.Deprecated)

	assert.Equal(t, StringType, arg.TypeInfo.ActualType)
	assert.False(t, arg.TypeInfo.IsPointer)
	assert.Nil(t, arg.TypeInfo.ItemType)
	assert.Empty(t, arg.TypeInfo.Enum)

	arg, err = extractArgument(reflMarker.Field(2))
	assert.Nil(t, err)
	assert.Equal(t, "Repeatable", arg.Name)
	assert.Empty(t, arg.Default)
	assert.False(t, arg.Required)
	assert.False(t, arg.Deprecated)

	assert.Equal(t, BoolType, arg.TypeInfo.ActualType)
	assert.False(t, arg.TypeInfo.IsPointer)
	assert.Nil(t, arg.TypeInfo.ItemType)
	assert.Empty(t, arg.TypeInfo.Enum)

	arg, err = extractArgument(reflMarker.Field(4))
	assert.Nil(t, err)
	assert.Equal(t, "Targets", arg.Name)
	assert.Empty(t, arg.Default)
	assert.True(t, arg.Required)
	assert.False(t, arg.Deprecated)

	assert.Equal(t, SliceType, arg.TypeInfo.ActualType)
	assert.Equal(t, StringType, arg.TypeInfo.ItemType.ActualType)
	assert.False(t, arg.TypeInfo.IsPointer)
	assert.Equal(t, map[string]interface{}{
		"FIELD_LEVEL":            "FIELD_LEVEL",
		"FUNCTION_LEVEL":         "FUNCTION_LEVEL",
		"INTERFACE_METHOD_LEVEL": "INTERFACE_METHOD_LEVEL",
		"INTERFACE_TYPE_LEVEL":   "INTERFACE_TYPE_LEVEL",
		"PACKAGE_LEVEL":          "PACKAGE_LEVEL",
		"STRUCT_METHOD_LEVEL":    "STRUCT_METHOD_LEVEL",
		"STRUCT_TYPE_LEVEL":      "STRUCT_TYPE_LEVEL",
	}, arg.TypeInfo.Enum)
}
