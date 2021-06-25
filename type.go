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

func (typeInfo TypeInfo) Parse(parser *Parser, out reflect.Value) error {
	switch typeInfo.ActualType {
	case BoolType:
		return typeInfo.parseBoolean(parser, out)
	case IntegerType:
		return typeInfo.parseInteger(parser, out)
	case StringType:
		return typeInfo.parseString(parser, out)
	case SliceType:
		return typeInfo.parseSlice(parser, out)
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

func (typeInfo TypeInfo) parseBoolean(parser *Parser, out reflect.Value) error {
	if parser == nil {
		return errors.New("parser cannot be nil")
	}

	if !parser.Expect(Identifier) {
		return nil
	}

	switch parser.Token() {
	case "false":
		typeInfo.setValue(out, reflect.ValueOf(false))
	case "true":
		typeInfo.setValue(out, reflect.ValueOf(true))
	}

	return fmt.Errorf("expected true or false, got %q", parser.Token())
}

func (typeInfo TypeInfo) parseInteger(parser *Parser, out reflect.Value) error {
	if parser == nil {
		return errors.New("parser cannot be nil")
	}

	nextCharacter := parser.Peek()

	isNegative := false

	if nextCharacter == '-' {
		isNegative = true
		parser.Scan()
	}

	if !parser.Expect(Integer) {
		return nil
	}

	text := parser.Token()

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

func (typeInfo TypeInfo) parseString(parser *Parser, out reflect.Value) error {
	if parser == nil {
		return errors.New("parser cannot be nil")
	}

	startPosition := parser.searchIndex

	token := parser.Scan()

	if token == String {

		value, err := strconv.Unquote(parser.Token())

		if err != nil {
			return nil
		}

		typeInfo.setValue(out, reflect.ValueOf(value))
		return nil
	}

	for character := parser.PeekWithoutSpace(); character != ',' && character != ';' && character != ':' && character != '}' && character != EOF; character = parser.PeekWithoutSpace() {
		parser.Scan()
	}

	endPosition := parser.searchIndex

	value := string(parser.markerComment[startPosition:endPosition])
	typeInfo.setValue(out, reflect.ValueOf(value))

	return nil
}

func (typeInfo TypeInfo) parseSlice(parser *Parser, out reflect.Value) error {
	if parser == nil {
		return errors.New("parser cannot be nil")
	}

	sliceType := reflect.Zero(out.Type())
	sliceItemType := reflect.Indirect(reflect.New(out.Type().Elem()))

	if parser.PeekWithoutSpace() == '{' {

		parser.Scan()
		for character := parser.PeekWithoutSpace(); character != '}' && character != EOF; character = parser.PeekWithoutSpace() {
			err := typeInfo.ItemType.Parse(parser, sliceItemType)

			if err != nil {
				return err
			}

			sliceType = reflect.Append(sliceType, sliceItemType)

			token := parser.PeekWithoutSpace()

			if token == '}' {
				break
			}

			if !parser.Expect(',') {
				return nil
			}
		}

		if !parser.Expect('}') {
			return nil
		}

		return nil
	}

	for character := parser.PeekWithoutSpace(); character != ',' && character != '}' && character != EOF; character = parser.PeekWithoutSpace() {
		err := typeInfo.ItemType.Parse(parser, sliceItemType)

		if err != nil {
			return err
		}

		sliceType = reflect.Append(sliceType, sliceItemType)

		token := parser.PeekWithoutSpace()

		if token == ',' || token == '}' || token == EOF {
			break
		}

		parser.Scan()

		if token != ';' {
			return nil
		}
	}

	typeInfo.setValue(out, sliceType)
	return nil
}
