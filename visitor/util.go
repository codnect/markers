package visitor

func IsInterfaceType(t Type) bool {
	_, ok := t.(*Interface)
	return ok
}

func IsStructType(t Type) bool {
	_, ok := t.(*Struct)
	return ok
}

func IsErrorType(t Type) bool {
	interfaceType, ok := t.(*Interface)
	if !ok {
		return false
	}

	if interfaceType.file == nil || interfaceType.file.pkg == nil {
		return false
	}

	return interfaceType.file.pkg.ID == "builtin"
}
