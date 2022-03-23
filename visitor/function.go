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

	funcDecl  *ast.FuncDecl
	funcField *ast.Field
	funcType  *ast.FuncType

	visitor            *packageVisitor
	loadedParams       bool
	loadedReturnValues bool
}

func newFunction(funcDecl *ast.FuncDecl, funcField *ast.Field, file *File) *Function {
	function := &Function{
		file:      file,
		params:    &Tuple{},
		results:   &Tuple{},
		funcDecl:  funcDecl,
		funcField: funcField,
	}

	if funcDecl != nil {
		function.name = funcDecl.Name.Name
		function.isExported = ast.IsExported(funcDecl.Name.Name)
		function.position = getPosition(file.pkg, funcDecl.Pos())
		function.funcType = funcDecl.Type
	} else {
		if funcField.Names != nil {
			function.name = funcField.Names[0].Name
		}
		function.funcType = funcField.Type.(*ast.FuncType)
	}

	return function.initialize()
}

func (f *Function) initialize() *Function {
	if f.funcDecl != nil {
		if f.funcDecl.Recv == nil {
			f.file.functions.elements = append(f.file.functions.elements, f)
		} else {
			f.receiver = &Variable{}

			if f.funcDecl.Recv.List[0].Names != nil {
				f.receiver.name = f.funcDecl.Recv.List[0].Names[0].Name
			}

			f.receiver.typ = f.receiverType(f.funcDecl.Recv.List[0].Type)
		}
	}

	return f
}

func (f *Function) receiverType(receiverExpr ast.Expr) Type {
	var receiverTypeSpec *ast.TypeSpec

	receiverTypeName := ""
	isPointerReceiver := false
	isStructMethod := false

	switch typedReceiver := receiverExpr.(type) {
	case *ast.Ident:
		if typedReceiver.Obj == nil {
			receiverTypeName = typedReceiver.Name
			unprocessedype := getTypeFromScope(receiverTypeName, f.visitor)
			_, isStructMethod = unprocessedype.(*Struct)
		} else {
			receiverTypeSpec = typedReceiver.Obj.Decl.(*ast.TypeSpec)
			receiverTypeName = receiverTypeSpec.Name.Name
			_, isStructMethod = receiverTypeSpec.Type.(*ast.StructType)
		}
	case *ast.StarExpr:
		if typedReceiver.X.(*ast.Ident).Obj == nil {
			receiverTypeName = typedReceiver.X.(*ast.Ident).Name
			unprocessedype := getTypeFromScope(receiverTypeName, f.visitor)
			_, isStructMethod = unprocessedype.(*Struct)
		} else {
			receiverTypeSpec = typedReceiver.X.(*ast.Ident).Obj.Decl.(*ast.TypeSpec)
			receiverTypeName = receiverTypeSpec.Name.Name
			_, isStructMethod = receiverTypeSpec.Type.(*ast.StructType)
		}
		isPointerReceiver = true
	}

	candidateType, ok := f.visitor.collector.findTypeByPkgIdAndName(f.file.pkg.ID, receiverTypeName)

	if isStructMethod {
		if !ok {
			candidateType = newStruct(receiverTypeSpec, nil, f.file, f.file.pkg, nil)
		}

		structType := candidateType.(*Struct)
		structType.methods = append(structType.methods, f)
	} else {
		if !ok {
			candidateType = newCustomType(receiverTypeSpec, f.file, nil, nil, nil)
		}

		customType := candidateType.(*CustomType)
		customType.methods = append(customType.methods, f)
	}

	if isPointerReceiver {
		return &Pointer{
			base: candidateType,
		}
	}

	return candidateType
}

func (f *Function) getVariables(fieldList []*ast.Field) *Tuple {
	tuple := &Tuple{
		variables: make([]*Variable, 0),
	}

	for _, field := range fieldList {

		typ := getTypeFromExpression(field.Type, f.visitor)

		if field.Names == nil {
			tuple.variables = append(tuple.variables, &Variable{
				typ: typ,
			})
		}

		for _, fieldName := range field.Names {
			tuple.variables = append(tuple.variables, &Variable{
				name: fieldName.Name,
				typ:  typ,
			})
		}

	}

	return tuple
}

func (f *Function) loadParams() {
	if f.loadedParams {
		return
	}

	if f.funcType.Params != nil {
		f.params.variables = append(f.params.variables, f.getVariables(f.funcType.Params.List).variables...)
	}

	f.loadedParams = true
}

func (f *Function) loadResultValues() {
	if f.loadedReturnValues {
		return
	}

	funcType := f.funcDecl.Type

	if funcType.Results != nil {
		f.results.variables = append(f.results.variables, f.getVariables(f.funcType.Results.List).variables...)
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
