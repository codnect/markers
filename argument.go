package marker

import (
	"reflect"
	"strings"
)

type Argument struct {
	Name     string
	TypeInfo ArgumentTypeInfo
	Pointer  bool
	Required bool
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
	argumentTypeInfo, err := GetArgumentTypeInfo(fieldType)

	if err != nil {
		return Argument{}, err
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
		TypeInfo: argumentTypeInfo,
		Pointer:  isPointer,
		Required: !optionalOption,
	}, nil
}
