package marker

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go/token"
	"testing"
)

func TestImportError_Error(t *testing.T) {
	importErr := &ImportError{
		Marker: "anyMarker",
	}
	assert.Equal(t, "the marker 'anyMarker' cannot be resolved", importErr.Error())
}

func TestToParseError(t *testing.T) {
	err := errors.New("anyError")
	convertedError := toParseError(err, nil, token.Position{
		Filename: "anyFileName",
		Offset:   0,
		Line:     10,
		Column:   13,
	})

	assert.NotNil(t, convertedError)
	parserError, isParserError := convertedError.(ParserError)
	assert.NotNil(t, convertedError)
	assert.True(t, isParserError)
	assert.Equal(t, err, parserError.error)
	assert.Equal(t, "anyFileName", parserError.FileName)
	assert.Equal(t, Position{Line: 10, Column: 13}, parserError.Position)
}

func TestToParseErrorWithErrorList(t *testing.T) {
	errors := NewErrorList([]error{errors.New("anyError1"), errors.New("anyError2")})

	convertedError := toParseError(errors, nil, token.Position{
		Filename: "anyFileName",
		Offset:   0,
		Line:     10,
		Column:   13,
	})

	assert.NotNil(t, convertedError)
	errorList, isErrorList := convertedError.(ErrorList)
	assert.NotNil(t, convertedError)
	assert.True(t, isErrorList)
	assert.Len(t, errorList, 2)

	parserError, isParserError := errorList[0].(ParserError)
	assert.NotNil(t, convertedError)
	assert.True(t, isParserError)
	assert.Equal(t, "anyError1", parserError.error.Error())
	assert.Equal(t, "anyFileName", parserError.FileName)
	assert.Equal(t, Position{Line: 10, Column: 13}, parserError.Position)

	parserError, isParserError = errorList[1].(ParserError)
	assert.NotNil(t, convertedError)
	assert.True(t, isParserError)
	assert.Equal(t, "anyError2", parserError.error.Error())
	assert.Equal(t, "anyFileName", parserError.FileName)
	assert.Equal(t, Position{Line: 10, Column: 13}, parserError.Position)
}
