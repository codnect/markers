package visitor

type Variadic struct {
	elem Type
}

func (v *Variadic) Name() string {
	return v.elem.Name()
}

func (v *Variadic) Elem() Type {
	return v.elem
}

func (v *Variadic) Underlying() Type {
	return v
}

func (v *Variadic) String() string {
	return ""
}
