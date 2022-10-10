package markers

import (
	"github.com/stretchr/testify/assert"
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
		typeInfo, err := ArgumentTypeInfoFromType(testCase.Type)
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

func TestArgumentTypeInfo_ParseString(t *testing.T) {
	typeInfo, err := ArgumentTypeInfoFromType(reflect.TypeOf("anyTest"))
	assert.Nil(t, err)
	assert.Equal(t, StringType, typeInfo.ActualType)

	strValue := ""

	scanner := NewScanner(" anyTestString ")
	scanner.Peek()

	err = typeInfo.parseString(scanner, reflect.ValueOf(&strValue))
	assert.Nil(t, err)
	assert.Equal(t, "anyTestString", strValue)

	scanner = NewScanner("\"anyTestString\"")
	scanner.Peek()

	err = typeInfo.parseString(scanner, reflect.ValueOf(&strValue))
	assert.Nil(t, err)
	assert.Equal(t, "anyTestString", strValue)

	scanner = NewScanner("`anyTestString`")
	scanner.Peek()

	err = typeInfo.parseString(scanner, reflect.ValueOf(&strValue))
	assert.Nil(t, err)
	assert.Equal(t, "anyTestString", strValue)

	scanner = NewScanner(" anyTestString ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&strValue))
	assert.Nil(t, err)
	assert.Equal(t, "anyTestString", strValue)

	scanner = NewScanner("\"anyTestString\"")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&strValue))
	assert.Nil(t, err)
	assert.Equal(t, "anyTestString", strValue)

	scanner = NewScanner("`anyTestString`")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&strValue))
	assert.Nil(t, err)
	assert.Equal(t, "anyTestString", strValue)
}

func TestArgumentTypeInfo_ParseBoolean(t *testing.T) {
	typeInfo, err := ArgumentTypeInfoFromType(reflect.TypeOf(true))
	assert.Nil(t, err)
	assert.Equal(t, BoolType, typeInfo.ActualType)

	boolValue := false

	scanner := NewScanner(" true ")
	scanner.Peek()

	err = typeInfo.parseBoolean(scanner, reflect.ValueOf(&boolValue))
	assert.Nil(t, err)
	assert.True(t, boolValue)

	scanner = NewScanner(" false ")
	scanner.Peek()

	err = typeInfo.parseBoolean(scanner, reflect.ValueOf(&boolValue))
	assert.Nil(t, err)
	assert.False(t, boolValue)

	scanner = NewScanner(" true ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&boolValue))
	assert.Nil(t, err)
	assert.True(t, boolValue)

	scanner = NewScanner(" false ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&boolValue))
	assert.Nil(t, err)
	assert.False(t, boolValue)
}

func TestArgumentTypeInfo_ParseInteger(t *testing.T) {
	typeInfo, err := ArgumentTypeInfoFromType(reflect.TypeOf(0))
	assert.Nil(t, err)
	assert.Equal(t, SignedIntegerType, typeInfo.ActualType)

	signedIntegerValue := 0

	scanner := NewScanner(" -091215 ")
	scanner.Peek()

	err = typeInfo.parseInteger(scanner, reflect.ValueOf(&signedIntegerValue))
	assert.Nil(t, err)
	assert.Equal(t, -91215, signedIntegerValue)

	scanner = NewScanner(" -070519 ")
	scanner.Peek()

	err = typeInfo.parseInteger(scanner, reflect.ValueOf(&signedIntegerValue))
	assert.Nil(t, err)
	assert.Equal(t, -70519, signedIntegerValue)

	typeInfo, err = ArgumentTypeInfoFromType(reflect.TypeOf(uint(0)))
	assert.Nil(t, err)
	assert.Equal(t, UnsignedIntegerType, typeInfo.ActualType)

	scanner = NewScanner(" -091215 ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&signedIntegerValue))
	assert.Nil(t, err)
	assert.Equal(t, -91215, signedIntegerValue)

	scanner = NewScanner(" -070519 ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&signedIntegerValue))
	assert.Nil(t, err)
	assert.Equal(t, -70519, signedIntegerValue)

	typeInfo, err = ArgumentTypeInfoFromType(reflect.TypeOf(uint(0)))
	assert.Nil(t, err)
	assert.Equal(t, UnsignedIntegerType, typeInfo.ActualType)

	unsignedIntegerValue := uint(0)

	scanner = NewScanner(" 091215 ")
	scanner.Peek()

	err = typeInfo.parseInteger(scanner, reflect.ValueOf(&unsignedIntegerValue))
	assert.Nil(t, err)
	assert.Equal(t, uint(91215), unsignedIntegerValue)

	scanner = NewScanner(" 070519 ")
	scanner.Peek()

	err = typeInfo.parseInteger(scanner, reflect.ValueOf(&unsignedIntegerValue))
	assert.Nil(t, err)
	assert.Equal(t, uint(70519), unsignedIntegerValue)

	scanner = NewScanner(" 091215 ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&unsignedIntegerValue))
	assert.Nil(t, err)
	assert.Equal(t, uint(91215), unsignedIntegerValue)

	scanner = NewScanner(" 070519 ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&unsignedIntegerValue))
	assert.Nil(t, err)
	assert.Equal(t, uint(70519), unsignedIntegerValue)
}

func TestArgumentTypeInfo_ParseMap(t *testing.T) {
	m := make(map[string]any)
	typeInfo, err := ArgumentTypeInfoFromType(reflect.TypeOf(&m))
	assert.Nil(t, err)
	assert.Equal(t, MapType, typeInfo.ActualType)
	assert.Equal(t, AnyType, typeInfo.ItemType.ActualType)

	scanner := NewScanner(" {anyKey1:123,anyKey2:true,anyKey3:\"anyValue1\",anyKey4:`anyValue2`} ")
	scanner.Peek()

	err = typeInfo.parseMap(scanner, reflect.ValueOf(&m))
	assert.Nil(t, err)
	assert.Contains(t, m, "anyKey1")
	assert.Equal(t, 123, m["anyKey1"])
	assert.Contains(t, m, "anyKey2")
	assert.Equal(t, true, m["anyKey2"])
	assert.Contains(t, m, "anyKey3")
	assert.Equal(t, "anyValue1", m["anyKey3"])
	assert.Contains(t, m, "anyKey4")
	assert.Equal(t, "anyValue2", m["anyKey4"])

	scanner = NewScanner(" {anyKey1:123,anyKey2:true,anyKey3:\"anyValue1\",anyKey4:`anyValue2`} ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&m))
	assert.Nil(t, err)
	assert.Contains(t, m, "anyKey1")
	assert.Equal(t, 123, m["anyKey1"])
	assert.Contains(t, m, "anyKey2")
	assert.Equal(t, true, m["anyKey2"])
	assert.Contains(t, m, "anyKey3")
	assert.Equal(t, "anyValue1", m["anyKey3"])
	assert.Contains(t, m, "anyKey4")
	assert.Equal(t, "anyValue2", m["anyKey4"])
}

func TestArgumentTypeInfo_ParseSlice(t *testing.T) {
	s := make([]int, 0)
	typeInfo, err := ArgumentTypeInfoFromType(reflect.TypeOf(&s))
	assert.Nil(t, err)
	assert.Equal(t, SliceType, typeInfo.ActualType)
	assert.Equal(t, SignedIntegerType, typeInfo.ItemType.ActualType)

	scanner := NewScanner(" 1;2;3;4;5 ")
	scanner.Peek()

	err = typeInfo.parseSlice(scanner, reflect.ValueOf(&s))
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s)

	scanner = NewScanner(" 1;2;3;4;5 ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&s))
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s)

	scanner = NewScanner(" {1,2,3,4,5} ")
	scanner.Peek()

	err = typeInfo.parseSlice(scanner, reflect.ValueOf(&s))
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s)

	scanner = NewScanner(" {1,2,3,4,5} ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&s))
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s)
}

func TestArgumentTypeInfo_TypeInference(t *testing.T) {
	var value any
	typeInfo, err := ArgumentTypeInfoFromType(reflect.TypeOf(&value))

	assert.Nil(t, err)
	assert.Equal(t, AnyType, typeInfo.ActualType)

	scanner := NewScanner(" {1,2,3,4,5} ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&value))
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{1, 2, 3, 4, 5}, value)

	scanner = NewScanner(" {anyValue1, anyValue2, anyValue3} ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&value))
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{"anyValue1", "anyValue2", "anyValue3"}, value)

	scanner = NewScanner(" {`anyValue1`, `anyValue2`, `anyValue3`} ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&value))
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{"anyValue1", "anyValue2", "anyValue3"}, value)

	scanner = NewScanner(" {\"anyValue1\", \"anyValue2\", \"anyValue3\"} ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&value))
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{"anyValue1", "anyValue2", "anyValue3"}, value)

	scanner = NewScanner(" {{1,2,3},{4,5,6}} ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&value))
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{[]interface{}{1, 2, 3}, []interface{}{4, 5, 6}}, value)

	scanner = NewScanner(" {{anyKey1:\"anyValue1\"},{anyKey2:32},{anyKey3:`anyValue3`}} ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&value))
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{map[string]interface{}{"anyKey1": "anyValue1"}, map[string]interface{}{"anyKey2": 32}, map[string]interface{}{"anyKey3": "anyValue3"}}, value)

	scanner = NewScanner(" {anyKey1:\"anyValue1\",anyKey2:true,anyKey3:23,anyKey4:{1,2,3},anyKey5:`anyValue5`} ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&value))
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{"anyKey1": "anyValue1", "anyKey2": true, "anyKey3": 23, "anyKey4": []interface{}{1, 2, 3}, "anyKey5": "anyValue5"}, value)

}

func TestDefinition_Parse(t *testing.T) {
	m := map[string]any{}
	typeInfo, err := ArgumentTypeInfoFromType(reflect.TypeOf(&m))
	assert.Nil(t, err)
	assert.Equal(t, MapType, typeInfo.ActualType)
	assert.Equal(t, AnyType, typeInfo.ItemType.ActualType)

	scanner := NewScanner(" {anyKey1:{anySubKey1:23,anySubKey2:\"anySubKeyValue2\"}} ")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&m))
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"anyKey1": map[string]interface{}{
			"anySubKey1": 23,
			"anySubKey2": "anySubKeyValue2",
		},
	}, m)

	s := make([]any, 0)
	typeInfo, err = ArgumentTypeInfoFromType(reflect.TypeOf(&s))
	assert.Nil(t, err)
	assert.Equal(t, SliceType, typeInfo.ActualType)
	assert.Equal(t, AnyType, typeInfo.ItemType.ActualType)

	scanner = NewScanner("{anyItem1, `anyItem2`, 2, true, -2}")
	scanner.Peek()

	err = typeInfo.Parse(scanner, reflect.ValueOf(&s))
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{"anyItem1", "anyItem2", 2, true, -2}, s)
}

func TestArgumentTypeInfo_MakeMapType(t *testing.T) {
	anyMap := map[string]any{}
	typeInfo, _ := ArgumentTypeInfoFromType(reflect.TypeOf(&anyMap))
	mapType, err := typeInfo.makeMapType()

	assert.Nil(t, err)
	result := reflect.MakeMap(mapType).Interface()

	assert.Equal(t, map[string]interface{}{}, result)

	intMap := map[string]int{}
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&intMap))
	mapType, err = typeInfo.makeMapType()

	assert.Nil(t, err)
	result = reflect.MakeMap(mapType).Interface()

	assert.Equal(t, map[string]int{}, result)

	boolMap := map[string]bool{}
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&boolMap))
	mapType, err = typeInfo.makeMapType()

	assert.Nil(t, err)
	result = reflect.MakeMap(mapType).Interface()

	assert.Equal(t, map[string]bool{}, result)

	uintMap := map[string]uint{}
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&uintMap))
	mapType, err = typeInfo.makeMapType()

	assert.Nil(t, err)
	result = reflect.MakeMap(mapType).Interface()

	assert.Equal(t, map[string]uint{}, result)

	stringMap := map[string]string{}
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&stringMap))
	mapType, err = typeInfo.makeMapType()

	assert.Nil(t, err)
	result = reflect.MakeMap(mapType).Interface()

	assert.Equal(t, map[string]string{}, result)

	sliceMap := map[string][]int{}
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&sliceMap))
	mapType, err = typeInfo.makeMapType()

	assert.Nil(t, err)
	result = reflect.MakeMap(mapType).Interface()

	assert.Equal(t, map[string][]int{}, result)

	mapOfMap := map[string]map[string]any{}
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&mapOfMap))
	mapType, err = typeInfo.makeMapType()

	assert.Nil(t, err)
	result = reflect.MakeMap(mapType).Interface()

	assert.Equal(t, map[string]map[string]any{}, result)
}

func TestArgumentTypeInfo_MakeSliceType(t *testing.T) {
	var anySlice []any
	typeInfo, _ := ArgumentTypeInfoFromType(reflect.TypeOf(&anySlice))
	sliceType, err := typeInfo.makeSliceType()

	assert.Nil(t, err)
	result := reflect.MakeSlice(sliceType, 0, 0).Interface()

	assert.Equal(t, []any{}, result)

	var intSlice []int
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&intSlice))
	sliceType, err = typeInfo.makeSliceType()

	assert.Nil(t, err)
	result = reflect.MakeSlice(sliceType, 0, 0).Interface()

	assert.Equal(t, []int{}, result)

	var boolSlice []bool
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&boolSlice))
	sliceType, err = typeInfo.makeSliceType()

	assert.Nil(t, err)
	result = reflect.MakeSlice(sliceType, 0, 0).Interface()

	assert.Equal(t, []bool{}, result)

	var uintSlice []uint
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&uintSlice))
	sliceType, err = typeInfo.makeSliceType()

	assert.Nil(t, err)
	result = reflect.MakeSlice(sliceType, 0, 0).Interface()

	assert.Equal(t, []uint{}, result)

	var stringSlice []string
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&stringSlice))
	sliceType, err = typeInfo.makeSliceType()

	assert.Nil(t, err)
	result = reflect.MakeSlice(sliceType, 0, 0).Interface()

	assert.Equal(t, []string{}, result)

	var sliceOfSlice [][]string
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&sliceOfSlice))
	sliceType, err = typeInfo.makeSliceType()

	assert.Nil(t, err)
	result = reflect.MakeSlice(sliceType, 0, 0).Interface()

	assert.Equal(t, [][]string{}, result)

	var sliceOfMap []map[string]any
	typeInfo, _ = ArgumentTypeInfoFromType(reflect.TypeOf(&sliceOfMap))
	sliceType, err = typeInfo.makeSliceType()

	assert.Nil(t, err)
	result = reflect.MakeSlice(sliceType, 0, 0).Interface()

	assert.Equal(t, []map[string]any{}, result)
}
