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
	Name        string
	Package     string
	TargetLevel TargetLevel
	Repeatable  bool
	Deprecated  bool
	Output      Output
}

func MakeDefinition(name, pkg string, level TargetLevel, output any) (*Definition, error) {
	if len(strings.TrimSpace(name)) == 0 {
		return nil, errors.New("marker name cannot be empty")
	}

	outputType := reflect.TypeOf(output)

	if outputType.Kind() == reflect.Ptr {
		outputType = outputType.Elem()
	}

	definition := &Definition{
		Name:        strings.TrimSpace(name),
		Package:     strings.TrimSpace(pkg),
		TargetLevel: level,
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

	err = definition.validate()

	if err != nil {
		return nil, err
	}

	return definition, nil
}

func (definition *Definition) validate() error {
	if definition.TargetLevel == 0 {
		return fmt.Errorf("specify target levels for the definition: %v", definition.Name)
	}

	if !IsLower(definition.Name) {
		return fmt.Errorf("marker '%s' should only contain lower case characters", definition.Name)
	}

	if strings.ContainsAny(definition.Name, " \t") {
		return fmt.Errorf("marker '%s' cannot contain any whitespace", definition.Name)
	}

	return nil
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

		/*
			if argumentInfo.SyntaxFree {
				definition.Output.SyntaxFree = argumentInfo.SyntaxFree
			}

			if argumentInfo.UseValueSyntax {
				definition.Output.UseValueSyntax = argumentInfo.UseValueSyntax
			}*/

		definition.Output.Fields[argumentInfo.Name] = argumentInfo
		definition.Output.FieldNames[argumentInfo.Name] = field.Name
	}

	if len(definition.Output.Fields) > 1 && definition.Output.SyntaxFree {
		return fmt.Errorf("output can only have 'Value' field since syntaxFree option is used")
	}

	return nil
}

// TODO fix parse method implementation
func (definition *Definition) Parse(comment string) (interface{}, error) {
	if definition.Output.SyntaxFree {
		return definition.parseSyntaxFree(comment), nil
	}

	output := reflect.Indirect(reflect.New(definition.Output.Type))
	comment = strings.TrimLeft(comment, " \t")
	_, _, fields := splitMarker(comment)

	isValueSyntax := true
	if strings.HasPrefix(comment, "+"+definition.Name+":") {
		isValueSyntax = false
	} else {
		tempComment := strings.Replace(comment, fields, "", -1)
		tempComment = strings.Replace(tempComment, "+"+definition.Name, "", -1)
		fields = tempComment + fields
	}

	var errs []error
	scanner := NewScanner(fields)
	scanner.ErrorCallback = func(scanner *Scanner, message string) {
		errs = append(errs, ScannerError{
			Message: message,
		})
	}
	seen := make(map[string]struct{}, len(definition.Output.Fields))

	if len(definition.Output.Fields) != 0 && scanner.Peek() != EOF {
		for {
			var argument Argument
			argumentName := ""
			argumentExists := false
			fieldName := ""
			fieldExists := false
			scanner.SkipWhitespaces()

			if !scanner.Expect(Identifier, "Value or Argument value") {
				if isValueSyntax {
					if "=" != scanner.Token() && len(seen) == 0 {
						break
					} else if "," == scanner.Token() {
						continue
					} else {
						argumentName = "Value"
						fieldName, fieldExists = definition.Output.FieldNames[argumentName]
						argument, argumentExists = definition.Output.Fields[argumentName]
						if !fieldExists || !argumentExists {
							var anyValue interface{}
							(&ArgumentTypeInfo{ActualType: AnyType}).Parse(scanner, reflect.ValueOf(&anyValue))
							//goto nextAttribute
						}
					}
				} else {
					if "," == scanner.Token() {
						continue
					}

					break
				}
			} else {
				argumentName = scanner.Token()
				fieldName, fieldExists = definition.Output.FieldNames[argumentName]
				argument, argumentExists = definition.Output.Fields[argumentName]
				if !fieldExists || !argumentExists {
					var anyValue interface{}
					(&ArgumentTypeInfo{ActualType: AnyType}).Parse(scanner, reflect.ValueOf(&anyValue))
					//goto nextAttribute
				}

				scanner.SkipWhitespaces()

				if !scanner.Expect('=', "Equal sign") {
					break
				}
			}

			seen[argumentName] = struct{}{}

			fieldValue := output.FieldByName(fieldName)

			if !fieldValue.CanSet() {
				break
			}

			err := argument.TypeInfo.Parse(scanner, fieldValue)

			if err != nil {
				break
			}

			//nextAttribute:
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

func (definition *Definition) parseSyntaxFree(marker string) any {
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

	if argument.TypeInfo.IsPointer {
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
