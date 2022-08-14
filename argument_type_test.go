package marker

import (
	"reflect"
	"testing"
)

func TestGetArgumentTypeInfo(t *testing.T) {
	var rawValue []byte
	var interfaceValue interface{}
	testCases := []struct {
		Type             reflect.Type
		MustHaveError    bool
		ExpectedType     ArgumentType
		ExpectedItemType *ArgumentTypeInfo
	}{

		{
			Type:             reflect.TypeOf(make([]interface{}, 0)),
			MustHaveError:    false,
			ExpectedType:     SliceType,
			ExpectedItemType: &ArgumentTypeInfo{ActualType: AnyType},
		},
		{
			Type:          reflect.TypeOf(&rawValue),
			MustHaveError: false,
			ExpectedType:  RawType,
		},
		{
			Type:          reflect.TypeOf(&interfaceValue),
			MustHaveError: false,
			ExpectedType:  AnyType,
		},
		{
			Type:          reflect.TypeOf(true),
			MustHaveError: false,
			ExpectedType:  BoolType,
		},
		{
			Type:          reflect.TypeOf(uint8(0)),
			MustHaveError: false,
			ExpectedType:  IntegerType,
		},
		{
			Type:          reflect.TypeOf(uint16(0)),
			MustHaveError: false,
			ExpectedType:  IntegerType,
		},
		{
			Type:          reflect.TypeOf(uint(0)),
			MustHaveError: false,
			ExpectedType:  IntegerType,
		},
		{
			Type:          reflect.TypeOf(uint32(0)),
			MustHaveError: false,
			ExpectedType:  IntegerType,
		},
		{
			Type:          reflect.TypeOf(uint64(0)),
			MustHaveError: false,
			ExpectedType:  IntegerType,
		},
		{
			Type:          reflect.TypeOf(int8(0)),
			MustHaveError: false,
			ExpectedType:  IntegerType,
		},
		{
			Type:          reflect.TypeOf(int16(0)),
			MustHaveError: false,
			ExpectedType:  IntegerType,
		},
		{
			Type:          reflect.TypeOf(0),
			MustHaveError: false,
			ExpectedType:  IntegerType,
		},
		{
			Type:          reflect.TypeOf(int32(0)),
			MustHaveError: false,
			ExpectedType:  IntegerType,
		},
		{
			Type:          reflect.TypeOf(int64(0)),
			MustHaveError: false,
			ExpectedType:  IntegerType,
		},
		{
			Type:          reflect.TypeOf("test"),
			MustHaveError: false,
			ExpectedType:  StringType,
		},
		{
			Type:             reflect.TypeOf(make([]bool, 0)),
			MustHaveError:    false,
			ExpectedType:     SliceType,
			ExpectedItemType: &ArgumentTypeInfo{ActualType: BoolType},
		},
		{
			Type:             reflect.TypeOf(make([]interface{}, 0)),
			MustHaveError:    false,
			ExpectedType:     SliceType,
			ExpectedItemType: &ArgumentTypeInfo{ActualType: AnyType},
		},
		{
			Type:             reflect.TypeOf(make(map[string]int)),
			MustHaveError:    false,
			ExpectedType:     MapType,
			ExpectedItemType: &ArgumentTypeInfo{ActualType: IntegerType},
		},
		{
			Type:             reflect.TypeOf(make(map[string]interface{})),
			MustHaveError:    false,
			ExpectedType:     MapType,
			ExpectedItemType: &ArgumentTypeInfo{ActualType: AnyType},
		},
		{
			Type: reflect.TypeOf(&struct {
			}{}),
			MustHaveError: true,
			ExpectedType:  InvalidType,
		},
	}

	for index, testCase := range testCases {
		typeInfo, err := GetArgumentTypeInfo(testCase.Type)
		if testCase.MustHaveError && err == nil {
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
