package markers

import (
	"errors"
	"reflect"
	"strings"
)

type Argument struct {
	Name       string
	TypeInfo   ArgumentTypeInfo
	Required   bool
	Deprecated bool
	Default    any
}

func extractArgument(structField reflect.StructField) (Argument, error) {
	parameterName := upperCamelCase(structField.Name)
	parameterTag, parameterTagExists := structField.Tag.Lookup("parameter")

	if parameterTagExists && parameterTag != "" {
		parameterName = parameterTag
	}

	fieldType := structField.Type
	argumentTypeInfo, err := ArgumentTypeInfoFromType(fieldType)

	if err != nil {
		return Argument{}, err
	}

	isRequired := false
	requiredTag, requiredTagExists := structField.Tag.Lookup("required")
	if requiredTagExists && requiredTag != "" {
		if requiredTag == "true" {
			isRequired = true
		}
	}

	isDeprecated := false
	_, deprecatedTagExists := structField.Tag.Lookup("deprecated")
	if deprecatedTagExists {
		isDeprecated = true
	}

	defaultValue := ""
	defaultTag, defaultTagExists := structField.Tag.Lookup("default")
	if defaultTagExists && defaultTag != "" {
		defaultValue = defaultTag
	}

	enumTag, enumTagExists := structField.Tag.Lookup("enum")
	if enumTagExists && enumTag != "" {
		if argumentTypeInfo.ActualType != StringType && (argumentTypeInfo.ActualType != SliceType || argumentTypeInfo.ItemType.ActualType != StringType) {
			return Argument{}, errors.New("string and string-slice can have enum tag")
		}

		enumValues := strings.Split(enumTag, ",")
		for _, enumValue := range enumValues {
			enumKeyValueParts := strings.SplitN(enumValue, "=", 2)
			if len(enumKeyValueParts) == 2 {
				argumentTypeInfo.Enum[enumKeyValueParts[0]] = enumKeyValueParts[1]
			} else {
				argumentTypeInfo.Enum[enumKeyValueParts[0]] = enumKeyValueParts[0]
			}
		}
	}

	return Argument{
		Name:       parameterName,
		TypeInfo:   argumentTypeInfo,
		Required:   isRequired,
		Deprecated: isDeprecated,
		Default:    defaultValue,
	}, nil
}
