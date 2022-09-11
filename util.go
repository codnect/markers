package marker

import (
	"strings"
	"unicode"
)

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func IsDecimal(character rune) bool {
	return '0' <= character && character <= '9'
}

func IsIdentifier(character rune, index int) bool {
	return character == '_' || unicode.IsLetter(character) || unicode.IsDigit(character) && index > 0
}

func LowerCamelCase(str string) string {
	isFirst := true

	return strings.Map(func(r rune) rune {
		if isFirst {
			isFirst = false
			return unicode.ToLower(r)
		}

		return r
	}, str)

}

func UpperCamelCase(str string) string {
	isFirst := true

	return strings.Map(func(r rune) rune {
		if isFirst {
			isFirst = false
			return unicode.ToUpper(r)
		}

		return r
	}, str)

}
