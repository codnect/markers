package marker

const Whitespace = 1<<'\t' | 1<<'\r' | 1<<' '

const (
	EOF = -(iota + 1)
	Identifier
	Integer
	String
)

type Scanner struct {
	markerComment      []byte
	tokenStartPosition int
	tokenEndPosition   int
	searchIndex        int
	character          rune
}

func NewScanner(markerComment string) *Scanner {
	return &Scanner{
		markerComment:      []byte(markerComment),
		character:          Identifier,
		tokenStartPosition: -1,
		tokenEndPosition:   0,
		searchIndex:        -1,
	}
}

func (scanner *Scanner) Peek() rune {
	if scanner.character == Identifier {
		scanner.character = scanner.Next()
	}

	return scanner.character
}

func (scanner *Scanner) Expect(expected rune) bool {
	token := scanner.Scan()

	if token != expected {
		return false
	}

	return true
}

func (scanner *Scanner) SetSearchIndex(searchIndex int) {
	if searchIndex >= len(scanner.markerComment) {
		searchIndex = len(scanner.markerComment)
		return
	}

	scanner.searchIndex = searchIndex
	scanner.character = rune(scanner.markerComment[searchIndex])
}

func (scanner *Scanner) Next() rune {
	scanner.searchIndex++

	if scanner.searchIndex >= len(scanner.markerComment) {
		return EOF
	}

	return rune(scanner.markerComment[scanner.searchIndex])
}

func (scanner *Scanner) SkipWhitespaces() rune {
	character := scanner.Peek()

	for Whitespace&(1<<uint(character)) != 0 {
		character = scanner.Next()
	}

	return character
}

func (scanner *Scanner) Scan() rune {
	character := scanner.SkipWhitespaces()

	token := character

	scanner.tokenStartPosition = scanner.searchIndex

	if IsIdentifier(character, 0) {
		token = Identifier
		character = scanner.ScanIdentifier()
	} else if IsDecimal(character) {
		token = Integer
		character = scanner.ScanNumber()
	} else if character == EOF {
		return EOF
	} else if character == '"' {
		token = String
		scanner.ScanString('"')
		character = scanner.Next()
	} else if character == '`' {
		token = String
		scanner.ScanString('`')
		character = scanner.Next()
	} else {
		character = scanner.Next()
	}

	scanner.tokenEndPosition = scanner.searchIndex
	scanner.character = character
	return token
}

func (scanner *Scanner) ScanNumber() rune {
	character := scanner.Next()

	for IsDecimal(character) {
		character = scanner.Next()
	}

	return character
}

func (scanner *Scanner) ScanIdentifier() rune {
	character := scanner.Next()

	for index := 1; IsIdentifier(character, index); index++ {
		character = scanner.Next()
	}

	return character
}

func (scanner *Scanner) ScanString(quote rune) (len int) {
	character := scanner.Next()

	for character != quote {
		if character == '\n' || character < 0 {
			return
		}

		character = scanner.Next()
		len++
	}

	return
}

func (scanner *Scanner) Token() string {
	if scanner.tokenStartPosition < 0 {
		return ""
	}

	return string(scanner.markerComment[scanner.tokenStartPosition:scanner.tokenEndPosition])
}
