package marker

import (
	"fmt"
	"reflect"
	"strings"
)

type Argument struct {
	Name     string
	Type     Type
	Pointer  bool
	Required bool

	ItemType *Type
}

func ExtractArgument(structField reflect.StructField) (Argument, error) {
	fieldName := LowerCamelCase(structField.Name)

	markerTag, tagExists := structField.Tag.Lookup("marker")
	markerTagValues := strings.Split(markerTag, ",")

	if tagExists && markerTagValues[0] != "" {
		fieldName = markerTagValues[0]
	}

	optionalOption := false

	for _, tagOption := range markerTagValues[1:] {

		if tagOption == "optional" {
			optionalOption = true
		}

	}

	fieldType := structField.Type
	argumentType, err := GetType(fieldType)

	if err != nil {
		return Argument{}, err
	}

	var argumentItemType Type

	if argumentType == SliceType || argumentType == MapType {
		itemType, err := GetType(fieldType.Elem())

		if err != nil && argumentType == SliceType {
			return Argument{}, fmt.Errorf("bad slice item type: %w", err)
		} else if err != nil && argumentType == MapType {
			return Argument{}, fmt.Errorf("bad map item type: %w", err)
		}

		argumentItemType = itemType
	}

	isPointer := false
	isOptional := false

	if fieldType.Kind() == reflect.Ptr {
		isPointer = true
		isOptional = true
	}

	optionalOption = optionalOption || isOptional

	return Argument{
		Name:     fieldName,
		Type:     argumentType,
		Pointer:  isPointer,
		Required: !optionalOption,
		ItemType: &argumentItemType,
	}, nil
}

/*
func (argument *Argument) parseArgument(parser *Parser, value reflect.Value, itemType bool) {
	typ := argument.Type

	if itemType {
		typ = *argument.ItemType
	}

	switch typ {
	case BoolType:
		if !parser.Expect(Identifier) {
			return
		}

		switch parser.Token() {
		case "true":
		case "false":
		default:
			return
		}
	case IntegerType:
		nextCharacter := parser.Peek()

		isNegative := false

		if nextCharacter == '-' {
			isNegative = true
			parser.Scan()
		}

		if !parser.Expect(Integer) {
			return
		}

		text := parser.Token()

		if isNegative {
			text = "-" + text
		}

		x, _ := strconv.Atoi(text)

		if x == 0 {

		}
	case StringType:
		startPosition := parser.searchIndex

		token := parser.Scan()

		if token == String {

			val, err := strconv.Unquote(parser.Token())

			if err != nil {

				return
			}

			if val == "" {

			}
			return
		}

		for character := scanForString(parser); character != ',' && character != ';' && character != ':' && character != '}' && character != EOF; character = scanForString(parser) {
			parser.Scan()
		}

		endPosition := parser.searchIndex

		value := string(parser.markerComment[startPosition:endPosition])

		if value == "" {

		}
	case SliceType:
		argument.parseSlice(parser, value)
	}
}

func (argument *Argument) parseSlice(parser *Parser, value reflect.Value) {

	resSlice := reflect.Zero(value.Type())
	elem := reflect.Indirect(reflect.New(value.Type().Elem()))

	if scanForString(parser) == '{' {

		parser.Scan()
		for hint := scanForString(parser); hint != '}' && hint != EOF; hint = scanForString(parser) {
			argument.parseArgument(parser, elem,true)

			resSlice = reflect.Append(resSlice, elem)

			tok := scanForString(parser)
			if tok == '}' {
				break
			}
			if !parser.Expect(',') {
				return
			}
		}

		if !parser.Expect('}') {
			return
		}

		return
	}

	for hint := scanForString(parser); hint != ',' && hint != '}' && hint != EOF; hint = scanForString(parser) {
		argument.parseArgument(parser, elem,true)

		resSlice = reflect.Append(resSlice, elem)

		token :=scanForString(parser)

		if token == ',' || token == '}' || token == EOF {
			break
		}

		parser.Scan()

		if token != ';' {
			return
		}
	}

	val :=resSlice.Interface()

	if val == nil {

	}
}

func scanForString(parser *Parser) rune {
	character := parser.Peek()

	for ; character <= ' ' && ((1<<uint64(character))&Whitespace) != 0; character = parser.Peek() {
		parser.Next()
	}

	return character
}
*/
