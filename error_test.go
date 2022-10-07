package marker

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go/token"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(errors.New("anyError"), "anyFileName", Position{
		Line:   2,
		Column: 4,
	})
	assert.Equal(t, Error{
		FileName: "anyFileName",
		Position: Position{
			Line:   2,
			Column: 4,
		},
		error: errors.New("anyError"),
	}, err)
}

func TestScannerError_Error(t *testing.T) {
	scannerError := &ScannerError{
		Message: "anyErrorMessage",
	}
	assert.Equal(t, "anyErrorMessage", scannerError.Error())
}

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

func TestErrorList_ToErrors(t *testing.T) {
	errorSlice := []error{errors.New("anyError1"), errors.New("anyError2")}
	anyErrorList := NewErrorList(errorSlice)
	assert.Equal(t, errorSlice, anyErrorList.(ErrorList).ToErrors())
}

func TestErrorList_Error(t *testing.T) {
	anyErrorList := NewErrorList([]error{errors.New("anyError1"), errors.New("anyError2")})
	assert.Equal(t, "[anyError1 anyError2]", anyErrorList.Error())
}

func TestToParseErrorWithErrorList(t *testing.T) {
	anyErrorList := NewErrorList([]error{errors.New("anyError1"), errors.New("anyError2")})

	convertedError := toParseError(anyErrorList, nil, token.Position{
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
