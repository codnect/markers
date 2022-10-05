package any

type errorList []error

func (e errorList) Print() {

}

func (e errorList) ToErrors() []error {
	return e
}
