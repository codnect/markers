package marker

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Argument names
const (
	ValueArgument = "Value"
)

type Output struct {
	Type              reflect.Type
	IsAnonymous       bool
	SyntaxFree        bool
	UseValueSyntax    bool
	AnonymousTypeInfo ArgumentTypeInfo
	Fields            map[string]Argument
	FieldNames        map[string]string
}

type Definition struct {
	Name   string
	Level  TargetLevel
	Output Output
	PkgId  string
}

func MakeDefinition(name string, pkgId string, level TargetLevel, output interface{}) (*Definition, error) {
	if len(strings.TrimSpace(name)) == 0 {
		return nil, errors.New("marker name cannot be empty")
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
		PkgId: strings.TrimSpace(pkgId),
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

		if argumentInfo.SyntaxFree {
			definition.Output.SyntaxFree = argumentInfo.SyntaxFree
		}

		if argumentInfo.UseValueSyntax {
			definition.Output.UseValueSyntax = argumentInfo.UseValueSyntax
		}

		definition.Output.Fields[argumentInfo.Name] = argumentInfo
		definition.Output.FieldNames[argumentInfo.Name] = field.Name
	}

	if len(definition.Output.Fields) > 1 && definition.Output.SyntaxFree {
		return fmt.Errorf("output can only have 'Value' field since syntaxFree option is used")
	}

	return nil
}

func (definition *Definition) Parse(marker string) (interface{}, error) {
	if definition.Output.SyntaxFree {
		return definition.parseSyntaxFree(marker), nil
	}

	output := reflect.Indirect(reflect.New(definition.Output.Type))

	name, anonymousName, fields := splitMarker(marker)

	if !definition.Output.UseValueSyntax && len(anonymousName) >= len(name)+1 {
		fields = anonymousName[len(name)+1:] + "=" + fields
	}

	var errs []error

	if strings.ContainsAny(anonymousName, ".,;=") {
		errs = append(errs, ScannerError{
			Message: fmt.Sprintf("Marker format is not valid : %s", marker),
		})
		return nil, NewErrorList(errs)
	}

	scanner := NewScanner(fields)
	scanner.ErrorCallback = func(scanner *Scanner, message string) {
		errs = append(errs, ScannerError{
			Message: message,
		})
	}

	valueArgumentProcessed := false
	canBeValueArgument := false

	seen := make(map[string]struct{}, len(definition.Output.Fields))

	if scanner.Peek() != EOF {
		for {
			var argumentName string
			currentCharacter := scanner.SkipWhitespaces()

			if definition.Output.UseValueSyntax && !valueArgumentProcessed && currentCharacter == '{' {
				canBeValueArgument = true
			} else if !scanner.Expect(Identifier, "Argument Name") {
				break
			}

			argumentName = scanner.Token()
			currentCharacter = scanner.SkipWhitespaces()

			if definition.Output.UseValueSyntax && !valueArgumentProcessed && (currentCharacter == EOF || currentCharacter == ',' || currentCharacter == ';') {
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

			seen[argumentName] = struct{}{}

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

	for argumentName, argument := range definition.Output.Fields {
		if _, wasSeen := seen[argumentName]; !wasSeen && argument.Required {
			scanner.AddError(fmt.Sprintf("missing argument %q", argumentName))
		}
	}

	return output.Interface(), NewErrorList(errs)
}

func (definition *Definition) parseSyntaxFree(marker string) interface{} {
	output := reflect.Indirect(reflect.New(definition.Output.Type))

	fieldName, exists := definition.Output.FieldNames[ValueArgument]

	if !exists {
		return output.Interface()
	}

	var argument Argument
	argument, exists = definition.Output.Fields[ValueArgument]

	if !exists {
		return output.Interface()
	}

	fieldValue := output.FieldByName(fieldName)

	if !fieldValue.CanSet() {
		return output.Interface()
	}

	fieldOutType := fieldValue.Type()

	if argument.Pointer {
		fieldOutType = fieldOutType.Elem()
		fieldValue = fieldValue.Elem()
	}

	name, _, _ := splitMarker(marker)
	// markers can be syntax free such as +build
	name = strings.Split(name, " ")[0]

	value := reflect.ValueOf(strings.Replace(marker, fmt.Sprintf("+%s", name), "", 1))

	if fieldOutType != value.Type() {
		value = value.Convert(fieldOutType)
	}

	fieldValue.Set(value)

	return output.Interface()
}
