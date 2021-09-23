package marker

import (
	"fmt"
	"go/ast"
	"go/token"
)

type ScannerError struct {
	Position int
	Message  string
}

func (err ScannerError) Error() string {
	return fmt.Sprintf("%s (at %d)", err.Message, err.Position)
}

type ImportError struct {
	Marker string
}

func (err ImportError) Error() string {
	return fmt.Sprintf("the marker '%s' cannot be resolved", err.Marker)
}

type ParserError struct {
	FileName string
	Position Position
	error
}

type ErrorList []error

func NewErrorList(errors []error) error {
	if len(errors) == 0 {
		return nil
	}

	return ErrorList(errors)
}

func (errorList ErrorList) Error() string {
	return fmt.Sprintf("%v", []error(errorList))
}

func toParseError(err error, node ast.Node, position token.Position) error {

	errorList, ok := err.(ErrorList)

	if !ok {

		errorPosition := Position{
			Line:   position.Line,
			Column: position.Column,
		}

		return ParserError{
			FileName: position.Filename,
			Position: errorPosition,
			error:    err,
		}
	}

	errors := make(ErrorList, len(errorList))

	for index, errorElement := range errorList {
		errors[index] = toParseError(errorElement, node, position)
	}

	return errors
}
