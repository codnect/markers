package marker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScanner_Scan(t *testing.T) {
	scanner := NewScanner("key1=123,key2=`test`,key3=\"hello\"")

	current := scanner.Scan()
	assert.Equal(t, Identifier, int(current))
	assert.Equal(t, "key1", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, '=', current)

	current = scanner.Scan()
	assert.Equal(t, Integer, int(current))
	assert.Equal(t, "123", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, ',', current)

	current = scanner.Scan()
	assert.Equal(t, Identifier, int(current))
	assert.Equal(t, "key2", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, '=', current)

	current = scanner.Scan()
	assert.Equal(t, String, int(current))
	assert.Equal(t, "`test`", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, ',', current)

	current = scanner.Scan()
	assert.Equal(t, Identifier, int(current))
	assert.Equal(t, "key3", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, '=', current)

	current = scanner.Scan()
	assert.Equal(t, String, int(current))
	assert.Equal(t, "\"hello\"", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, EOF, int(current))
}

func TestScanner_ScanIdentifier(t *testing.T) {
	scanner := NewScanner("key=123")

	current := scanner.ScanIdentifier()
	assert.Equal(t, '=', current)
	assert.Equal(t, "key", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, '=', current)

	current = scanner.Scan()
	assert.Equal(t, Integer, int(current))
	assert.Equal(t, "123", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, EOF, int(current))
}

func TestScanner_ScanNumber(t *testing.T) {
	scanner := NewScanner("key=123")

	current := scanner.Scan()
	assert.Equal(t, Identifier, int(current))
	assert.Equal(t, "key", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, '=', current)

	current = scanner.ScanNumber()
	assert.Equal(t, EOF, int(current))
	assert.Equal(t, "123", scanner.Token())
}

func TestScanner_ScanString(t *testing.T) {
	scanner := NewScanner("key=\"hello\"")

	current := scanner.Scan()
	assert.Equal(t, Identifier, int(current))
	assert.Equal(t, "key", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, '=', current)

	length := scanner.ScanString('"')
	assert.Equal(t, len("hello"), length)
	assert.Equal(t, "\"hello\"", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, EOF, int(current))
}

func TestScanner_AllScans(t *testing.T) {
	scanner := NewScanner("key1=`test`,key2=123,key3=\"hello\"")

	current := scanner.ScanIdentifier()
	assert.Equal(t, '=', current)
	assert.Equal(t, "key1", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, '=', current)

	length := scanner.ScanString('`')
	assert.Equal(t, len("test"), length)
	assert.Equal(t, "`test`", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, ',', current)

	current = scanner.ScanIdentifier()
	assert.Equal(t, '=', current)
	assert.Equal(t, "key2", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, '=', current)

	current = scanner.ScanNumber()
	assert.Equal(t, ',', current)
	assert.Equal(t, "123", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, ',', current)

	current = scanner.ScanIdentifier()
	assert.Equal(t, '=', current)
	assert.Equal(t, "key3", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, '=', current)

	length = scanner.ScanString('"')
	assert.Equal(t, len("hello"), length)
	assert.Equal(t, "\"hello\"", scanner.Token())

	current = scanner.Scan()
	assert.Equal(t, EOF, int(current))
}
