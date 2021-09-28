package marker

import (
	"fmt"
	"go/ast"
	"go/token"
)

type Error struct {
	FileName string
	Position Position
	error
}

func NewError(err error, fileName string, position Position) error {
	return Error{
		fileName,
		position,
		err,
	}
}

type ScannerError struct {
	Message string
}

func (err ScannerError) Error() string {
	return err.Message
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
		return ParserError{
			FileName: position.Filename,
			Position: Position{
				Line:   position.Line,
				Column: position.Column,
			},
			error: err,
		}
	}

	errors := make(ErrorList, len(errorList))

	for index, errorElement := range errorList {
		errors[index] = toParseError(errorElement, node, position)
	}

	return errors
}
