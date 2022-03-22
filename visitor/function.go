package visitor

type Variable struct {
	name string
	typ  Type
}

func (v *Variable) Name() string {
	return v.name
}

func (v *Variable) Type() Type {
	return v.typ
}

type Tuple struct {
	variables []*Variable
}

func (t *Tuple) Len() int {
	return len(t.variables)
}

func (t *Tuple) At(index int) *Variable {
	if index >= 0 && index < len(t.variables) {
		return t.variables[index]
	}

	return nil
}

type Function struct {
	name       string
	isExported bool
	position   Position
	receiver   *Variable
	params     *Tuple
	results    *Tuple
	variadic   bool

	file *File
}

func (f *Function) File() *File {
	return f.file
}

func (f *Function) Position() Position {
	return f.position
}

func (f *Function) Underlying() Type {
	return f
}

func (f *Function) String() string {
	return ""
}

func (f *Function) Receiver() *Variable {
	return f.receiver
}

func (f *Function) Params() *Tuple {
	return f.params
}

func (f *Function) Results() *Tuple {
	return f.results
}

func (f *Function) IsVariadic() bool {
	return f.variadic
}

type Functions struct {
	elements []*Function
}

func (f *Functions) Len() int {
	return len(f.elements)
}

func (f *Functions) At(index int) *Function {
	if index >= 0 && index < len(f.elements) {
		return f.elements[index]
	}

	return nil
}
