package marker

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Output struct {
	Type              reflect.Type
	IsAnonymous       bool
	AnonymousTypeInfo ArgumentTypeInfo
	Fields            map[string]Argument
	FieldNames        map[string]string
}

type Definition struct {
	Name           string
	Level          TargetLevel
	Output         Output
	UseValueSyntax bool
}

func MakeDefinition(name string, level TargetLevel, output interface{}, useValueSyntax ...bool) (*Definition, error) {
	if len(strings.TrimSpace(name)) == 0 {
		return nil, errors.New("marker name cannot be empty")
	}

	nameParts := strings.Split(name, ":")
	if ImportMarkerName == nameParts[0] {
		return nil, errors.New("import is reserved for marker project, please select another marker name")
	}

	outputType := reflect.TypeOf(output)

	if outputType.Kind() == reflect.Ptr {
		outputType = outputType.Elem()
	}

	definition := &Definition{
		Name:  strings.TrimSpace(name),
		Level: level,
		Output: Output{
			Type:        outputType,
			IsAnonymous: false,
			Fields:      make(map[string]Argument),
			FieldNames:  make(map[string]string),
		},
	}

	if useValueSyntax != nil {
		definition.UseValueSyntax = useValueSyntax[0]
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

	name, anonymousName, fields := splitMarker(marker)

	if !definition.UseValueSyntax && len(anonymousName) >= len(name)+1 {
		fields = anonymousName[len(name)+1:] + "=" + fields
	}

	var errs []error

	if strings.ContainsAny(anonymousName, ".,;=") {
		errs = append(errs, ScannerError{
			Position: 0,
			Message:  fmt.Sprintf("Marker format is not valid : %s", marker),
		})
		return nil, NewErrorList(errs)
	}

	scanner := NewScanner(fields)
	scanner.ErrorCallback = func(scanner *Scanner, message string) {
		errs = append(errs, ScannerError{
			Position: scanner.SearchIndex(),
			Message:  message,
		})
	}

	valueArgumentProcessed := false
	canBeValueArgument := false

	if scanner.Peek() != EOF {
		for {
			var argumentName string
			currentCharacter := scanner.SkipWhitespaces()

			if definition.UseValueSyntax && !valueArgumentProcessed && currentCharacter == '{' {
				canBeValueArgument = true
			} else if !scanner.Expect(Identifier, "Argument Name") {
				break
			}

			argumentName = scanner.Token()
			currentCharacter = scanner.SkipWhitespaces()

			if definition.UseValueSyntax && !valueArgumentProcessed && (currentCharacter == ',' || currentCharacter == ';') {
				canBeValueArgument = true
			} else if (valueArgumentProcessed || !canBeValueArgument) && !scanner.Expect('=', "Equals Sign '='") {
				break
			}

			if canBeValueArgument && !valueArgumentProcessed {
				valueArgumentProcessed = true
				argumentName = ValueArgument
				scanner.Reset()
			}

			fieldName, exists := definition.Output.FieldNames[argumentName]

			var err error
			var fieldValue reflect.Value
			var argument Argument

			// if the argument name does not exist in field names, parse its value to skip
			if !exists {
				var anyValue interface{}
				(&ArgumentTypeInfo{ActualType: AnyType}).Parse(scanner, reflect.ValueOf(&anyValue))
				goto nextAttribute
			}

			argument, exists = definition.Output.Fields[argumentName]

			// if the argument name does not exist in fields, parse its value to skip
			if !exists {
				var anyValue interface{}
				(&ArgumentTypeInfo{ActualType: AnyType}).Parse(scanner, reflect.ValueOf(&anyValue))
				goto nextAttribute
			}

			fieldValue = output.FieldByName(fieldName)

			if !fieldValue.CanSet() {
				break
			}

			err = argument.TypeInfo.Parse(scanner, fieldValue)

			if err != nil {
				break
			}

		nextAttribute:
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
