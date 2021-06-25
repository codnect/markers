package marker

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Type int

const (
	InvalidType Type = iota
	RawType
	AnyType
	BoolType
	IntegerType
	StringType
	SliceType
	MapType
)

var (
	interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
	rawType       = reflect.TypeOf((*[]byte)(nil)).Elem()
)

type TypeInfo struct {
	ActualType Type
	ItemType   *TypeInfo
}

func GetTypeInfo(typ reflect.Type) (TypeInfo, error) {
	typeInfo := &TypeInfo{}

	if typ == rawType {
		typeInfo.ActualType = RawType
		return *typeInfo, nil
	}

	if typ == interfaceType {
		typeInfo.ActualType = AnyType
		return *typeInfo, nil
	}

	if typ.Kind() == reflect.Ptr {
		rawType = typ.Elem()
	}

	switch typ.Kind() {
	case reflect.String:
		typeInfo.ActualType = StringType
	case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
		typeInfo.ActualType = IntegerType
	case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
		typeInfo.ActualType = IntegerType
	case reflect.Bool:
		typeInfo.ActualType = BoolType
	case reflect.Slice:
		typeInfo.ActualType = SliceType
		itemType, err := GetTypeInfo(typ.Elem())

		if err != nil {
			return TypeInfo{}, fmt.Errorf("bad slice item type: %w", err)
		}

		typeInfo.ItemType = &itemType
	case reflect.Map:
		if typ.Key().Kind() != reflect.String {
			return TypeInfo{}, fmt.Errorf("map key must be string")
		}

		typeInfo.ActualType = MapType
		itemType, err := GetTypeInfo(typ.Elem())

		if err != nil {
			return TypeInfo{}, fmt.Errorf("bad map item type: %w", err)
		}

		typeInfo.ItemType = &itemType
	default:
		return TypeInfo{}, fmt.Errorf("type has unsupported kind %s", typ.Kind())
	}

	return *typeInfo, nil
}

func (typeInfo TypeInfo) Parse(scanner *Scanner, out reflect.Value) error {
	switch typeInfo.ActualType {
	case BoolType:
		return typeInfo.parseBoolean(scanner, out)
	case IntegerType:
		return typeInfo.parseInteger(scanner, out)
	case StringType:
		return typeInfo.parseString(scanner, out)
	case SliceType:
		return typeInfo.parseSlice(scanner, out)
	case MapType:
		return typeInfo.parseMap(scanner, out)
	}

	return nil
}

func (typeInfo TypeInfo) setValue(out, value reflect.Value) {
	outType := out.Type()

	if outType != value.Type() {
		value = value.Convert(outType)
	}

	out.Set(value)
}

func (typeInfo TypeInfo) parseBoolean(scanner *Scanner, out reflect.Value) error {
	if scanner == nil {
		return errors.New("scanner cannot be nil")
	}

	if !scanner.Expect(Identifier) {
		return nil
	}

	switch scanner.Token() {
	case "false":
		typeInfo.setValue(out, reflect.ValueOf(false))
	case "true":
		typeInfo.setValue(out, reflect.ValueOf(true))
	}

	return fmt.Errorf("expected true or false, got %q", scanner.Token())
}

func (typeInfo TypeInfo) parseInteger(scanner *Scanner, out reflect.Value) error {
	if scanner == nil {
		return errors.New("scanner cannot be nil")
	}

	nextCharacter := scanner.Peek()

	isNegative := false

	if nextCharacter == '-' {
		isNegative = true
		scanner.Scan()
	}

	if !scanner.Expect(Integer) {
		return nil
	}

	text := scanner.Token()

	if isNegative {
		text = "-" + text
	}

	intValue, err := strconv.Atoi(text)

	typeInfo.setValue(out, reflect.ValueOf(intValue))

	if err != nil {
		return fmt.Errorf("unable to parse integer: %v", err)
	}

	return nil
}

func (typeInfo TypeInfo) parseString(scanner *Scanner, out reflect.Value) error {
	if scanner == nil {
		return errors.New("scanner cannot be nil")
	}

	startPosition := scanner.searchIndex

	token := scanner.Scan()

	if token == String {

		value, err := strconv.Unquote(scanner.Token())

		if err != nil {
			return nil
		}

		typeInfo.setValue(out, reflect.ValueOf(value))
		return nil
	}

	for character := scanner.SkipWhitespaces(); character != ',' && character != ';' && character != ':' && character != '}' && character != EOF; character = scanner.SkipWhitespaces() {
		scanner.Scan()
	}

	endPosition := scanner.searchIndex

	value := string(scanner.markerComment[startPosition:endPosition])
	typeInfo.setValue(out, reflect.ValueOf(value))

	return nil
}

func (typeInfo TypeInfo) parseSlice(scanner *Scanner, out reflect.Value) error {
	if scanner == nil {
		return errors.New("scanner cannot be nil")
	}

	sliceType := reflect.Zero(out.Type())
	sliceItemType := reflect.Indirect(reflect.New(out.Type().Elem()))

	if scanner.SkipWhitespaces() == '{' {

		scanner.Scan()

		for character := scanner.SkipWhitespaces(); character != '}' && character != EOF; character = scanner.SkipWhitespaces() {
			err := typeInfo.ItemType.Parse(scanner, sliceItemType)

			if err != nil {
				return err
			}

			sliceType = reflect.Append(sliceType, sliceItemType)

			token := scanner.SkipWhitespaces()

			if token == '}' {
				break
			}

			if !scanner.Expect(',') {
				return nil
			}
		}

		if !scanner.Expect('}') {
			return nil
		}

		return nil
	}

	for character := scanner.SkipWhitespaces(); character != ',' && character != '}' && character != EOF; character = scanner.SkipWhitespaces() {
		err := typeInfo.ItemType.Parse(scanner, sliceItemType)

		if err != nil {
			return err
		}

		sliceType = reflect.Append(sliceType, sliceItemType)

		token := scanner.SkipWhitespaces()

		if token == ',' || token == '}' || token == EOF {
			break
		}

		scanner.Scan()

		if token != ';' {
			return nil
		}
	}

	typeInfo.setValue(out, sliceType)
	return nil
}

func (typeInfo TypeInfo) parseMap(scanner *Scanner, out reflect.Value) error {
	if scanner == nil {
		return errors.New("scanner cannot be nil")
	}

	mapType := reflect.MakeMap(out.Type())
	key := reflect.Indirect(reflect.New(out.Type().Key()))
	value := reflect.Indirect(reflect.New(out.Type().Elem()))

	if !scanner.Expect('{') {
		return nil
	}

	for character := scanner.SkipWhitespaces(); character != '}' && character != EOF; character = scanner.SkipWhitespaces() {
		err := typeInfo.parseString(scanner, key)

		if err != nil {
			return err
		}

		if !scanner.Expect(':') {
			return nil
		}

		err = typeInfo.ItemType.Parse(scanner, value)

		if err != nil {
			return err
		}

		mapType.SetMapIndex(key, value)

		if scanner.SkipWhitespaces() == '}' {
			break
		}

		if !scanner.Expect(',') {
			return nil
		}
	}

	if !scanner.Expect('}') {
		return nil
	}

	typeInfo.setValue(out, mapType)

	return nil
}
