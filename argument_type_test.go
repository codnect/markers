package marker

import (
	"reflect"
	"testing"
)

func TestGetArgumentTypeInfo(t *testing.T) {
	var rawValue []byte
	var interfaceValue interface{}
	testCases := []struct {
		Type              reflect.Type
		ShouldReturnError bool
		ExpectedType      ArgumentType
		ExpectedItemType  *ArgumentTypeInfo
	}{

		{
			Type:              reflect.TypeOf(make([]interface{}, 0)),
			ShouldReturnError: false,
			ExpectedType:      SliceType,
			ExpectedItemType:  &ArgumentTypeInfo{ActualType: AnyType},
		},
		{
			Type:              reflect.TypeOf(&rawValue),
			ShouldReturnError: false,
			ExpectedType:      RawType,
		},
		{
			Type:              reflect.TypeOf(&interfaceValue),
			ShouldReturnError: false,
			ExpectedType:      AnyType,
		},
		{
			Type:              reflect.TypeOf(true),
			ShouldReturnError: false,
			ExpectedType:      BoolType,
		},
		{
			Type:              reflect.TypeOf(uint8(0)),
			ShouldReturnError: false,
			ExpectedType:      UnsignedIntegerType,
		},
		{
			Type:              reflect.TypeOf(uint16(0)),
			ShouldReturnError: false,
			ExpectedType:      UnsignedIntegerType,
		},
		{
			Type:              reflect.TypeOf(uint(0)),
			ShouldReturnError: false,
			ExpectedType:      UnsignedIntegerType,
		},
		{
			Type:              reflect.TypeOf(uint32(0)),
			ShouldReturnError: false,
			ExpectedType:      UnsignedIntegerType,
		},
		{
			Type:              reflect.TypeOf(uint64(0)),
			ShouldReturnError: false,
			ExpectedType:      UnsignedIntegerType,
		},
		{
			Type:              reflect.TypeOf(int8(0)),
			ShouldReturnError: false,
			ExpectedType:      SignedIntegerType,
		},
		{
			Type:              reflect.TypeOf(int16(0)),
			ShouldReturnError: false,
			ExpectedType:      SignedIntegerType,
		},
		{
			Type:              reflect.TypeOf(0),
			ShouldReturnError: false,
			ExpectedType:      SignedIntegerType,
		},
		{
			Type:              reflect.TypeOf(int32(0)),
			ShouldReturnError: false,
			ExpectedType:      SignedIntegerType,
		},
		{
			Type:              reflect.TypeOf(int64(0)),
			ShouldReturnError: false,
			ExpectedType:      SignedIntegerType,
		},
		{
			Type:              reflect.TypeOf("test"),
			ShouldReturnError: false,
			ExpectedType:      StringType,
		},
		{
			Type:              reflect.TypeOf(make([]bool, 0)),
			ShouldReturnError: false,
			ExpectedType:      SliceType,
			ExpectedItemType:  &ArgumentTypeInfo{ActualType: BoolType},
		},
		{
			Type:              reflect.TypeOf(make([]interface{}, 0)),
			ShouldReturnError: false,
			ExpectedType:      SliceType,
			ExpectedItemType:  &ArgumentTypeInfo{ActualType: AnyType},
		},
		{
			Type:              reflect.TypeOf(make(map[string]int)),
			ShouldReturnError: false,
			ExpectedType:      MapType,
			ExpectedItemType:  &ArgumentTypeInfo{ActualType: SignedIntegerType},
		},
		{
			Type:              reflect.TypeOf(make(map[string]interface{})),
			ShouldReturnError: false,
			ExpectedType:      MapType,
			ExpectedItemType:  &ArgumentTypeInfo{ActualType: AnyType},
		},
		{
			Type: reflect.TypeOf(&struct {
			}{}),
			ShouldReturnError: true,
			ExpectedType:      InvalidType,
		},
	}

	for index, testCase := range testCases {
		typeInfo, err := GetArgumentTypeInfo(testCase.Type)
		if testCase.ShouldReturnError && err == nil {
			t.Errorf("%d. test case must have an error.", index)
		}

		if typeInfo.ActualType != testCase.ExpectedType {
			t.Errorf("actual type is not equal to expected, got %q; want %q", typeInfo.ActualType, testCase.ExpectedType)
		}

		if testCase.ExpectedItemType != nil && typeInfo.ItemType.ActualType != testCase.ExpectedItemType.ActualType {
			t.Errorf("item type is not equal to expected, got %q; want %q", typeInfo.ItemType.ActualType, testCase.ExpectedItemType.ActualType)
		}
	}
}
