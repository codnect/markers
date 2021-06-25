package marker

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

type Parser struct {
	markerComment      []byte
	tokenStartPosition int
	tokenEndPosition   int
	searchIndex        int
	character          rune
}

func NewParser(markerComment string) *Parser {
	return &Parser{
		markerComment:      []byte(markerComment),
		character:          Identifier,
		tokenStartPosition: -1,
		tokenEndPosition:   0,
		searchIndex:        -1,
	}
}

func (parser *Parser) Peek() rune {
	if parser.character == Identifier {
		parser.character = parser.Next()
	}

	return parser.character
}

func (parser *Parser) Expect(expected rune) bool {
	token := parser.Scan()

	if token != expected {
		return false
	}

	return true
}

func (parser *Parser) Next() rune {
	parser.searchIndex++

	if parser.searchIndex >= len(parser.markerComment) {
		return EOF
	}

	return rune(parser.markerComment[parser.searchIndex])
}

func (parser *Parser) SkipWhitespaces() rune {
	character := parser.Peek()

	for Whitespace&(1<<uint(character)) != 0 {
		character = parser.Next()
	}

	return character
}

func (parser *Parser) Scan() rune {
	character := parser.SkipWhitespaces()

	token := character

	parser.tokenStartPosition = parser.searchIndex

	if IsIdentifier(character, 0) {
		token = Identifier
		character = parser.ScanIdentifier()
	} else if IsDecimal(character) {
		token = Integer
		character = parser.ScanNumber()
	} else if character == EOF {
		return EOF
	} else if character == '"' {
		token = String
		parser.ScanString('"')
		character = parser.Next()
	} else {
		character = parser.Next()
	}

	parser.tokenEndPosition = parser.searchIndex
	parser.character = character
	return token
}

func (parser *Parser) ScanNumber() rune {
	character := parser.Next()

	for IsDecimal(character) {
		character = parser.Next()
	}

	return character
}

func (parser *Parser) ScanIdentifier() rune {
	character := parser.Next()

	for index := 1; IsIdentifier(character, index); index++ {
		character = parser.Next()
	}

	return character
}

func (parser *Parser) ScanString(quote rune) (len int) {
	character := parser.Next()

	for character != quote {
		if character == '\n' || character < 0 {
			return
		}

		character = parser.Next()
		len++
	}

	return
}

func (parser *Parser) PeekWithoutSpace() rune {
	character := parser.Peek()

	for ; character <= ' ' && ((1<<uint64(character))&Whitespace) != 0; character = parser.Peek() {
		parser.Next()
	}

	return character
}

func (parser *Parser) Token() string {
	if parser.tokenStartPosition < 0 {
		return ""
	}

	return string(parser.markerComment[parser.tokenStartPosition:parser.tokenEndPosition])
}
