package visitor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
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

type Variables []*Variable

func (v Variables) Len() int {
	return len(v)
}

func (v Variables) At(index int) *Variable {
	if index >= 0 && index < len(v) {
		return v[index]
	}

	return nil
}

type Function struct {
	name       string
	isExported bool
	markers    marker.MarkerValues
	position   Position
	receiver   *Variable
	typeParams *TypeParams
	params     Variables
	results    Variables
	variadic   bool

	file *File

	funcDecl  *ast.FuncDecl
	funcField *ast.Field
	funcType  *ast.FuncType

	pkg     *packages.Package
	visitor *packageVisitor

	loadedTypeParams   bool
	loadedParams       bool
	loadedReturnValues bool
}

func newFunction(funcDecl *ast.FuncDecl, funcField *ast.Field, file *File, pkg *packages.Package, visitor *packageVisitor, markers marker.MarkerValues) *Function {
	function := &Function{
		file:       file,
		typeParams: &TypeParams{},
		params:     Variables{},
		results:    Variables{},
		markers:    markers,
		funcDecl:   funcDecl,
		funcField:  funcField,
		pkg:        pkg,
		visitor:    visitor,
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
		function.isExported = ast.IsExported(function.name)
		function.position = getPosition(file.pkg, function.funcType.Pos())
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
			candidateType = newStruct(receiverTypeSpec, nil, f.file, f.pkg, f.visitor, nil)
		}

		structType := candidateType.(*Struct)
		structType.methods = append(structType.methods, f)
	} else {
		if !ok {
			candidateType = newCustomType(receiverTypeSpec, f.file, f.pkg, f.visitor, nil)
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

func (f *Function) getTypeParams(fieldList []*ast.Field) *TypeParams {
	typeParams := &TypeParams{
		params: make([]*TypeParam, 0),
	}

	for _, field := range fieldList {

		typ := getTypeFromExpression(field.Type, f.visitor)

		if field.Names == nil {
			typeParams.params = append(typeParams.params, &TypeParam{
				typ: typ,
			})
		}

		for _, fieldName := range field.Names {
			typeParams.params = append(typeParams.params, &TypeParam{
				name: fieldName.Name,
				typ:  typ,
			})
		}

	}

	return typeParams
}

func (f *Function) getTypeParameterByName(name string) *TypeParam {
	f.loadTypeParams()
	/*for _, typeParam := range f.typeParams.variables {
		if typeParam.name == name {
			return typeParam
		}
	}*/

	return nil
}

func (f *Function) getGenericTypeFromExpression(exp ast.Expr) Type {
	var typeParam *TypeParam

	switch t := exp.(type) {
	case *ast.Ident:
		typeParam = f.getTypeParameterByName(t.Name)
	case *ast.SelectorExpr:
	}

	if typeParam == nil {
		return nil
	}

	return &Generic{
		typeParam,
	}
}

func (f *Function) getVariables(fieldList []*ast.Field) Variables {
	variables := Variables{}

	for _, field := range fieldList {
		typ := f.getGenericTypeFromExpression(field.Type)

		if typ == nil {
			typ = getTypeFromExpression(field.Type, f.visitor)
		}

		if field.Names == nil {
			variables = append(variables, &Variable{
				typ: typ,
			})
		}

		for _, fieldName := range field.Names {
			variables = append(variables, &Variable{
				name: fieldName.Name,
				typ:  typ,
			})
		}

	}

	return variables
}

func (f *Function) loadTypeParams() {

	if f.loadedTypeParams {
		return
	}

	if f.funcType.TypeParams != nil {
		f.typeParams.params = append(f.typeParams.params, f.getTypeParams(f.funcType.TypeParams.List).params...)
	}

	f.loadedTypeParams = true
}

func (f *Function) loadParams() {
	if f.loadedParams {
		return
	}

	if f.funcType.Params != nil {
		f.params = append(f.params, f.getVariables(f.funcType.Params.List)...)
	}

	if f.params.Len() != 0 {
		_, f.variadic = f.params.At(f.params.Len() - 1).Type().(*Variadic)
	}

	f.loadedParams = true
}

func (f *Function) loadResultValues() {
	if f.loadedReturnValues {
		return
	}

	if f.funcType.Results != nil {
		f.results = append(f.results, f.getVariables(f.funcType.Results.List)...)
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

func (f *Function) TypeParams() *TypeParams {
	f.loadTypeParams()
	return f.typeParams
}

func (f *Function) Params() Variables {
	f.loadParams()
	return f.params
}

func (f *Function) Results() Variables {
	f.loadResultValues()
	return f.results
}

func (f *Function) IsVariadic() bool {
	f.loadParams()
	return f.variadic
}

func (f *Function) Markers() marker.MarkerValues {
	return f.markers
}

type Functions struct {
	elements []*Function
}

func (f *Functions) ToSlice() []*Function {
	return f.elements
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

func (f *Functions) FindByName(name string) (*Function, bool) {
	for _, function := range f.elements {
		if function.name == name {
			return function, true
		}
	}

	return nil, false
}
