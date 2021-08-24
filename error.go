package marker

import "fmt"

type ParserError struct {
	FileName string
	Line     int
	Message  string
}

func (parserError *ParserError) Error() string {
	return ""
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
