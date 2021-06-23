package marker

import (
	"unicode"
)

const Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '

const (
	EOF = -(iota + 1)
	Identifier
	Integer
	Float
	Character
	String
	RawString
)

type parser struct {
	markerComment      []byte
	tokenStartPosition int
	tokenEndPosition   int
	searchIndex        int
	character          rune
}

func newParser(markerComment string) *parser {
	return &parser{
		markerComment:      []byte(markerComment),
		character:          Identifier,
		tokenStartPosition: -1,
		tokenEndPosition:   0,
		searchIndex:        -1,
	}
}

func (parser *parser) peek() rune {
	if parser.character == Identifier {
		parser.character = parser.next()
	}

	return parser.character
}

func (parser *parser) expect(expected rune) bool {
	token := parser.scan()

	if token != expected {
		return false
	}

	return true
}

func (parser *parser) next() rune {
	parser.searchIndex++

	if parser.searchIndex >= len(parser.markerComment) {
		return EOF
	}

	return rune(parser.markerComment[parser.searchIndex])
}

func (parser *parser) scan() rune {
	character := parser.peek()

	for Whitespace&(1<<uint(character)) != 0 {
		character = parser.next()
	}

	token := character

	parser.tokenStartPosition = parser.searchIndex

	if parser.isIdentifier(character, 0) {
		token = Identifier
		character = parser.scanIdentifier()
	} else if parser.isDecimal(character) {
		token = Integer
		character = parser.scanNumber()
	} else if character == EOF {
		return EOF
	} else if character == '"' {
		token = String
		parser.scanString('"')
		character = parser.next()
	} else {
		character = parser.next()
	}

	parser.tokenEndPosition = parser.searchIndex
	parser.character = character
	return token
}

func (parser *parser) isIdentifier(character rune, index int) bool {
	return character == '_' || unicode.IsLetter(character) || unicode.IsDigit(character) && index > 0
}

func (parser *parser) isDecimal(character rune) bool {
	return '0' <= character && character <= '9'
}

func (parser *parser) lower(character rune) rune {
	return ('a' - 'A') | character
}

func (parser *parser) isHex(ch rune) bool {
	return '0' <= ch && ch <= '9' || 'a' <= parser.lower(ch) && parser.lower(ch) <= 'f'
}

func (parser *parser) scanNumber() rune {
	character := parser.next()

	for parser.isDecimal(character) {
		character = parser.next()
	}

	return character
}

func (parser *parser) scanIdentifier() rune {
	character := parser.next()

	for index := 1; parser.isIdentifier(character, index); index++ {
		character = parser.next()
	}

	return character
}

func (parser *parser) scanString(quote rune) (n int) {
	character := parser.next()

	for character != quote {
		if character == '\n' || character < 0 {
			return
		}

		if character == '\\' {
			character = parser.next()
		} else {
			character = parser.next()
		}

		n++
	}

	return
}

func (parser *parser) token() string {
	if parser.tokenStartPosition < 0 {
		return ""
	}

	return string(parser.markerComment[parser.tokenStartPosition:parser.tokenEndPosition])
}
