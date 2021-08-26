package marker

import (
	"errors"
	"reflect"
)

type Output struct {
	Type              reflect.Type
	IsAnonymous       bool
	AnonymousTypeInfo ArgumentTypeInfo
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
			Type:        outputType,
			IsAnonymous: false,
			Fields:      make(map[string]Argument),
			FieldNames:  make(map[string]string),
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
		argumentTypeInfo, err := GetArgumentTypeInfo(definition.Output.Type)

		if err != nil {
			return err
		}

		definition.Output.IsAnonymous = true
		definition.Output.AnonymousTypeInfo = argumentTypeInfo
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

		if argumentInfo.TypeInfo.ActualType == RawType {
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

	var errs []error

	scanner := NewScanner(marker)
	scanner.ErrorCallback = func(scanner *Scanner, message string) {
		errs = append(errs, ScannerError{
			Position: scanner.SearchIndex(),
			Message:  message,
		})
	}

	if scanner.Peek() != EOF {
		for {
			if !scanner.Expect(Identifier, "Argument Name") {
				break
			}

			argumentName := scanner.Token()

			if !scanner.Expect('=', "Equals Sign '='") {
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

			err := argument.TypeInfo.Parse(scanner, fieldValue)

			if err != nil {
				break
			}

			if scanner.Peek() == EOF {
				break
			}
			if !scanner.Expect(',', "Comma ','") {
				break
			}
		}
	}

	return output.Interface(), NewErrorList(errs)
}
