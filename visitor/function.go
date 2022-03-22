package visitor

import (
	"go/ast"
	"strings"
)

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

func (v *Variable) String() string {
	return ""
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

	funcDecl *ast.FuncDecl

	loadedParams       bool
	loadedReturnValues bool
}

func newFunction(funcDecl *ast.FuncDecl, file *File) *Function {
	function := &Function{
		name:       funcDecl.Name.Name,
		isExported: ast.IsExported(funcDecl.Name.Name),
		file:       file,
		position:   getPosition(file.pkg, funcDecl.Pos()),
		params:     &Tuple{},
		results:    &Tuple{},
	}

	return function.initialize()
}

func (f *Function) initialize() *Function {
	if f.funcDecl.Recv == nil {
		f.file.functions.elements = append(f.file.functions.elements, f)
	}

	return f
}

func (f *Function) loadParams() {
	if f.loadedParams {
		return
	}

	funcType := f.funcDecl.Type

	if funcType.Params != nil {
		f.params.variables = append(f.params.variables, getVariables(funcType.Params.List).variables...)
	}

	f.loadedParams = true
}

func (f *Function) loadResultValues() {
	if f.loadedReturnValues {
		return
	}

	funcType := f.funcDecl.Type

	if funcType.Results != nil {
		f.results.variables = append(f.results.variables, getVariables(funcType.Results.List).variables...)
	}

	f.loadedReturnValues = true
}

func (f *Function) Name() string {
	return f.name
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
	f.loadParams()
	f.loadResultValues()

	var builder strings.Builder
	builder.WriteString("func ")

	if f.receiver != nil {
		builder.WriteString("(")
		builder.WriteString(f.receiver.Name())
		builder.WriteString(" ")
		builder.WriteString(f.receiver.Type().String())
		builder.WriteString(") ")
	}

	builder.WriteString(f.name)
	builder.WriteString("(")

	if f.Params().Len() != 0 {
		for i := 0; i < f.Params().Len(); i++ {
			param := f.Params().At(i)
			builder.WriteString(param.String())

			if i != f.Params().Len()-1 {
				builder.WriteString(",")
			}
		}
	}

	builder.WriteString(") ")

	if f.Results().Len() > 1 {
		builder.WriteString("(")
	}

	if f.Results().Len() != 0 {
		for i := 0; i < f.Results().Len(); i++ {
			result := f.Results().At(i)
			builder.WriteString(result.String())

			if i != f.Params().Len()-1 {
				builder.WriteString(",")
			}
		}
	}

	if f.Results().Len() > 1 {
		builder.WriteString(")")
	}

	return builder.String()
}

func (f *Function) Receiver() *Variable {
	return f.receiver
}

func (f *Function) Params() *Tuple {
	f.loadParams()
	return f.params
}

func (f *Function) Results() *Tuple {
	f.loadResultValues()
	return f.results
}

func (f *Function) IsVariadic() bool {
	f.loadParams()
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
