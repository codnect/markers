package visitor

type Array struct {
	len  int64
	elem Type
}

func (a *Array) Len() int64 {
	return a.len
}

func (a *Array) Elem() Type {
	return a.elem
}

func (a *Array) Name() string {
	return "[]"
}

func (a *Array) Underlying() Type {
	return a
}

func (a *Array) String() string {
	return ""
}
