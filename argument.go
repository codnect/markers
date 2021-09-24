package marker

import (
	"fmt"
	"reflect"
	"strings"
)

type Argument struct {
	Name           string
	TypeInfo       ArgumentTypeInfo
	Pointer        bool
	Required       bool
	SyntaxFree     bool
	UseValueSyntax bool
}

func ExtractArgument(structField reflect.StructField) (Argument, error) {
	fieldName := LowerCamelCase(structField.Name)

	markerTag, tagExists := structField.Tag.Lookup("marker")
	markerTagValues := strings.Split(markerTag, ",")

	if tagExists && markerTagValues[0] != "" {
		fieldName = markerTagValues[0]
	}

	optionalOption := false
	syntaxFree := false
	useValueSyntax := false

	for _, tagOption := range markerTagValues[1:] {

		if tagOption == "optional" {
			optionalOption = true
		}

		if tagOption == "syntaxFree" {
			syntaxFree = true
		}

		if tagOption == "useValueSyntax" {
			useValueSyntax = true
		}
	}

	if ValueArgument != fieldName && syntaxFree {
		return Argument{}, fmt.Errorf("'Value' field can only have syntaxFree option")
	}

	if ValueArgument != fieldName && useValueSyntax {
		return Argument{}, fmt.Errorf("'Value' field can only have useValueSyntax option")
	}

	if ValueArgument == fieldName && syntaxFree && useValueSyntax {
		return Argument{}, fmt.Errorf("'Value' cannot have both syntaxFree and useValueSyntax options at the same time")
	}

	fieldType := structField.Type
	argumentTypeInfo, err := GetArgumentTypeInfo(fieldType)

	if err != nil {
		return Argument{}, err
	}

	if syntaxFree && argumentTypeInfo.ActualType != StringType {
		return Argument{}, fmt.Errorf("'Value' field with syntaxFree option can be only string")
	}

	isPointer := false
	isOptional := false

	if fieldType.Kind() == reflect.Ptr {
		isPointer = true
		isOptional = true
	}

	optionalOption = optionalOption || isOptional

	return Argument{
		Name:           fieldName,
		TypeInfo:       argumentTypeInfo,
		Pointer:        isPointer,
		Required:       !optionalOption,
		SyntaxFree:     syntaxFree,
		UseValueSyntax: useValueSyntax,
	}, nil
}
