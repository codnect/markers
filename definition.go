package marker

import (
	"errors"
	"fmt"
	"reflect"
)

type Output struct {
	Type              reflect.Type
	IsAnonymous       bool
	AnonymousType     Type
	AnonymousItemType *Type
	Fields            map[string]Argument
	FieldNames        map[string]string
}

type Definition struct {
	Name   string
	Level  TargetLevel
	Output Output
}

func MakeDefinition(name string, level TargetLevel, output interface{}) (*Definition, error) {
	outputType := reflect.TypeOf(output)

	if outputType.Kind() == reflect.Ptr {
		outputType = outputType.Elem()
	}

	definition := &Definition{
		Name:  name,
		Level: level,
		Output: Output{
			Type:          nil,
			IsAnonymous:   false,
			AnonymousType: InvalidType,
			Fields:        nil,
			FieldNames:    nil,
		},
	}

	err := definition.extract()

	if err != nil {
		return nil, err
	}

	return definition, nil
}

func (definition *Definition) extract() error {

	if definition.Output.Type.Kind() != reflect.Struct {
		argumentType, err := GetType(definition.Output.Type)

		if err != nil {
			return err
		}

		definition.Output.IsAnonymous = true
		definition.Output.AnonymousType = argumentType

		if argumentType == SliceType || argumentType == MapType {
			itemType, err := GetType(definition.Output.Type.Elem())

			if err != nil && argumentType == SliceType {
				return fmt.Errorf("bad slice item type: %w", err)
			} else if err != nil && argumentType == MapType {
				return fmt.Errorf("bad map item type: %w", err)
			}

			definition.Output.AnonymousItemType = &itemType
		}

		return nil
	}

	for index := 0; index < definition.Output.Type.NumField(); index++ {
		field := definition.Output.Type.Field(index)

		if field.PkgPath != "" {
			continue
		}

		argumentInfo, err := ExtractArgument(field)

		if err != nil {
			return err
		}

		if argumentInfo.Type == RawType {
			return errors.New("RawArgument cannot be a field")
		}

		definition.Output.Fields[argumentInfo.Name] = argumentInfo
		definition.Output.FieldNames[argumentInfo.Name] = field.Name
	}

	return nil
}

func (definition *Definition) Parse(marker string) (interface{}, error) {
	output := reflect.Indirect(reflect.New(definition.Output.Type))

	splitMarker(marker)

	parser := NewParser(marker)

	if parser.Peek() != EOF {
		for {
			if !parser.Expect(Identifier) {
				break
			}

			argumentName := parser.Token()

			if !parser.Expect('=') {
				break
			}

			fieldName, exists := definition.Output.FieldNames[argumentName]

			if !exists {
				break
			}

			argument, exists := definition.Output.Fields[argumentName]

			if !exists {
				break
			}

			fieldValue := output.FieldByName(fieldName)

			if !fieldValue.CanSet() {
				break
			}

			//argument.parseArgument(parser, fieldValue, false)
		}
	}

	return nil, nil
}
