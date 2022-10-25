package markers

import (
	"strings"
	"unicode"
)

func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func isLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func isDecimal(character rune) bool {
	return '0' <= character && character <= '9'
}

func isIdentifier(character rune, index int) bool {
	return character == '_' || unicode.IsLetter(character) || unicode.IsDigit(character) && index > 0
}

func lowerCamelCase(str string) string {
	isFirst := true

	return strings.Map(func(r rune) rune {
		if isFirst {
			isFirst = false
			return unicode.ToLower(r)
		}

		return r
	}, str)

}

func upperCamelCase(str string) string {
	isFirst := true

	return strings.Map(func(r rune) rune {
		if isFirst {
			isFirst = false
			return unicode.ToUpper(r)
		}

		return r
	}, str)

}
