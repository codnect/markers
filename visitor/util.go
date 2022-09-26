package visitor

func IsInterfaceType(t Type) bool {
	_, ok := t.(*Interface)
	return ok
}

func IsStructType(t Type) bool {
	_, ok := t.(*Struct)
	return ok
}
