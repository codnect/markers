package any

type errorList []error

func (e errorList) Print() {

}

// +deprecated any deprecation message

// ToErrors returns an array of errors
func (e errorList) ToErrors() []error {
	return e
}
