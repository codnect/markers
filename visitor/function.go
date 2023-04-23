package visitor

import (
	"fmt"
	"github.com/procyon-projects/markers"
	"github.com/procyon-projects/markers/packages"
	"go/ast"
	"strings"
	"sync"
)

type Parameter struct {
	name string
	typ  Type
}

func (p *Parameter) Name() string {
	return p.name
}

func (p *Parameter) Type() Type {
	return p.typ
}

func (p *Parameter) String() string {
	_, isTypeParameter := p.Type().(*TypeParameter)

	if p.name == "" {
		if isTypeParameter {
			return p.typ.Name()
		}

		return p.typ.String()
	}

	if isTypeParameter {
		return fmt.Sprintf("%s %s", p.name, p.typ.Name())
	}

	return fmt.Sprintf("%s %s", p.name, p.typ.String())
}

type Parameters struct {
	elements []*Parameter
}

func (p *Parameters) Len() int {
	return len(p.elements)
}

func (p *Parameters) At(index int) *Parameter {
	if index >= 0 && index < len(p.elements) {
		return p.elements[index]
	}

	return nil
}

func (p *Parameters) FindByName(name string) (*Parameter, bool) {
	for _, parameter := range p.elements {
		if parameter.name == name {
			return parameter, true
		}
	}

	return nil, false
}

type Result struct {
	name string
	typ  Type
}

func (r *Result) Name() string {
	return r.name
}

func (r *Result) Type() Type {
	return r.typ
}

func (r *Result) String() string {
	_, isTypeParameter := r.Type().(*TypeParameter)

	if r.name == "" {
		if isTypeParameter {
			return r.typ.Name()
		}

		return r.typ.String()
	}

	if isTypeParameter {
		return fmt.Sprintf("%s %s", r.name, r.typ.Name())
	}

	return fmt.Sprintf("%s %s", r.name, r.typ.String())
}

type Results struct {
	elements []*Result
}

func (r *Results) Len() int {
	return len(r.elements)
}

func (r *Results) At(index int) *Result {
	if index >= 0 && index < len(r.elements) {
		return r.elements[index]
	}

	return nil
}

func (r *Results) FindByName(name string) (*Result, bool) {
	for _, result := range r.elements {
		if result.name == name {
			return result, true
		}
	}

	return nil, false
}

type Receiver struct {
	name string
	typ  Type
}

func (r *Receiver) Name() string {
	return r.name
}

func (r *Receiver) Type() Type {
	return r.typ
}

func (r *Receiver) String() string {
	if r.name == "" {
		return r.typ.Name()
	}

	return fmt.Sprintf("%s %s", r.name, r.typ.Name())
}

type Function struct {
	name               string
	isExported         bool
	markers            markers.Values
	position           Position
	receiver           *Receiver
	typeParams         *TypeParameters
	receiverTypeParams *TypeParameters
	typeParamAliases   []string
	params             *Parameters
	results            *Results
	variadic           bool
	ownerType          Type

	file *File

	funcDecl  *ast.FuncDecl
	funcField *ast.Field
	funcType  *ast.FuncType

	pkg     *packages.Package
	visitor *packageVisitor

	typeParamsOnce sync.Once
	paramsOnce     sync.Once
	resultsOnce    sync.Once
}

func newFunction(funcDecl *ast.FuncDecl, funcType *ast.FuncType, funcField *ast.Field, ownerType Type, file *File, pkg *packages.Package, visitor *packageVisitor, markers markers.Values) *Function {
	function := &Function{
		file: file,
		typeParams: &TypeParameters{
			[]*TypeParameter{},
		},
		receiverTypeParams: &TypeParameters{
			[]*TypeParameter{},
		},
		params: &Parameters{
			[]*Parameter{},
		},
		results: &Results{
			[]*Result{},
		},
		markers:   markers,
		funcDecl:  funcDecl,
		funcField: funcField,
		funcType:  funcType,
		pkg:       pkg,
		visitor:   visitor,
		ownerType: ownerType,
	}

	if funcDecl != nil {
		function.name = funcDecl.Name.Name
		function.isExported = ast.IsExported(funcDecl.Name.Name)
		function.position = getPosition(file.pkg, funcDecl.Pos())
		function.funcType = funcDecl.Type
	} else if funcType != nil {
		function.position = getPosition(file.pkg, function.funcType.Pos())
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
			f.receiver = &Receiver{}

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
	isStructMethod := false

	switch typedReceiver := receiverExpr.(type) {
	case *ast.Ident:
		if typedReceiver.Obj == nil {
			receiverTypeName = typedReceiver.Name
			unprocessedType := getTypeFromScope(receiverTypeName, f.visitor)
			_, isStructMethod = unprocessedType.(*Struct)
		} else {
			receiverTypeSpec = typedReceiver.Obj.Decl.(*ast.TypeSpec)
			receiverTypeName = receiverTypeSpec.Name.Name
			_, isStructMethod = receiverTypeSpec.Type.(*ast.StructType)
		}
	case *ast.IndexExpr:
		f.typeParamAliases = append(f.typeParamAliases, typedReceiver.Index.(*ast.Ident).Name)
		return f.receiverType(typedReceiver.X)
	case *ast.IndexListExpr:
		for _, typeParamAlias := range typedReceiver.Indices {
			f.typeParamAliases = append(f.typeParamAliases, typeParamAlias.(*ast.Ident).Name)
		}
		return f.receiverType(typedReceiver.X)
	case *ast.StarExpr:
		return &Pointer{
			base: f.receiverType(typedReceiver.X),
		}
	}

	candidateType, ok := f.visitor.collector.findTypeByPkgIdAndName(f.file.pkg.ID, receiverTypeName)

	if isStructMethod {
		if !ok {
			candidateType = newStruct(receiverTypeSpec, nil, f.file, f.pkg, f.visitor, nil)
		}

		structType := candidateType.(*Struct)
		f.ownerType = structType
		structType.methods = append(structType.methods, f)
		f.file.functions.elements = append(f.file.functions.elements, f)
	} else {
		if !ok {
			candidateType = newCustomType(receiverTypeSpec, f.file, f.pkg, f.visitor, nil)
		}

		customType := candidateType.(*CustomType)
		f.ownerType = customType
		customType.methods = append(customType.methods, f)
		f.file.functions.elements = append(f.file.functions.elements, f)
	}

	return candidateType
}

func (f *Function) Name() string {
	if f.name == "" {
		return f.String()
	}

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

func (f *Function) Receiver() *Receiver {
	return f.receiver
}

func (f *Function) TypeParameters() *TypeParameters {
	f.loadTypeParams()
	if f.ownerType != nil {
		_, isStruct := f.ownerType.(*Struct)
		_, isCustomType := f.ownerType.(*CustomType)
		if !isStruct && !isCustomType {
			return &TypeParameters{}
		}

		return f.receiverTypeParams
	}

	return f.typeParams
}

func (f *Function) Parameters() *Parameters {
	f.loadParams()
	return f.params
}

func (f *Function) Results() *Results {
	f.loadResultValues()
	return f.results
}

func (f *Function) IsVariadic() bool {
	f.loadParams()
	return f.variadic
}

func (f *Function) Markers() markers.Values {
	return f.markers
}

func (f *Function) getResults(fieldList []*ast.Field) []*Result {
	variables := make([]*Result, 0)

	for _, field := range fieldList {
		typ := getTypeFromExpression(field.Type, f.file, f.visitor, nil, f.typeParams)

		if field.Names == nil {
			variables = append(variables, &Result{
				typ: typ,
			})
		}

		for _, fieldName := range field.Names {
			variables = append(variables, &Result{
				name: fieldName.Name,
				typ:  typ,
			})
		}

	}

	return variables
}

func (f *Function) getParameters(fieldList []*ast.Field) []*Parameter {
	variables := make([]*Parameter, 0)

	for index, field := range fieldList {
		typ := getTypeFromExpression(field.Type, f.file, f.visitor, nil, f.typeParams)

		if field.Names == nil {
			variables = append(variables, &Parameter{
				typ: typ,
			})
		}

		for _, fieldName := range field.Names {
			variables = append(variables, &Parameter{
				name: fieldName.Name,
				typ:  typ,
			})
		}

		if index == len(fieldList)-1 {
			f.variadic = true
		}
	}

	if len(variables) != 0 {
		_, f.variadic = variables[len(variables)-1].Type().(*Variadic)
	}

	return variables
}

func (f *Function) String() string {
	f.loadParams()
	f.loadResultValues()

	var builder strings.Builder
	builder.WriteString("func ")

	if f.receiver != nil {
		builder.WriteString("(")
		if f.receiver.Name() != "" {
			builder.WriteString(f.receiver.Name())
			builder.WriteString(" ")
		}

		builder.WriteString(f.receiver.Type().Name())

		if f.TypeParameters().Len() != 0 {
			builder.WriteString("[")
			for i := 0; i < f.TypeParameters().Len(); i++ {
				typeParam := f.TypeParameters().At(i)
				builder.WriteString(typeParam.Name())

				if i != f.TypeParameters().Len()-1 {
					builder.WriteString(",")
				}
			}
			builder.WriteString("]")
		}

		if f.name != "" {
			builder.WriteString(") ")
		} else {
			builder.WriteString(")")
		}
	}

	builder.WriteString(f.name)

	if f.ownerType == nil && f.TypeParameters().Len() != 0 {
		builder.WriteString("[")
		for i := 0; i < f.TypeParameters().Len(); i++ {
			typeParam := f.TypeParameters().At(i)
			builder.WriteString(typeParam.String())

			if i != f.TypeParameters().Len()-1 {
				builder.WriteString(",")
			}
		}
		builder.WriteString("]")
	}

	builder.WriteString("(")

	if f.Parameters().Len() != 0 {
		for i := 0; i < f.Parameters().Len(); i++ {
			param := f.Parameters().At(i)
			builder.WriteString(param.String())

			if i != f.Parameters().Len()-1 {
				builder.WriteString(",")
			}
		}
	}

	if f.Results().Len() == 0 {
		builder.WriteString(")")
	} else {
		builder.WriteString(") ")
	}

	if f.Results().Len() > 1 {
		builder.WriteString("(")
	}

	if f.Results().Len() != 0 {
		for i := 0; i < f.Results().Len(); i++ {
			result := f.Results().At(i)
			builder.WriteString(result.String())

			if i != f.Results().Len()-1 {
				builder.WriteString(",")
			}
		}
	}

	if f.Results().Len() > 1 {
		builder.WriteString(")")
	}

	return builder.String()
}

func (f *Function) loadTypeParams() {
	f.typeParamsOnce.Do(func() {
		if f.ownerType != nil {

			switch typedOwner := f.ownerType.(type) {
			case *CustomType:
				f.typeParams.elements = append(f.typeParams.elements, typedOwner.TypeParameters().elements...)

				for index, typeParamAlias := range f.typeParamAliases {
					if typeParameter, exists := f.typeParams.FindByName(typeParamAlias); exists {
						f.receiverTypeParams.elements = append(f.receiverTypeParams.elements, typeParameter)
						continue
					}

					typeParam := f.typeParams.At(index)

					if typeParam != nil && typeParam.Name() != typeParamAlias {
						typeParameter := &TypeParameter{
							typeParamAlias,
							typeParam.TypeConstraints(),
						}
						f.typeParams.elements = append(f.typeParams.elements, typeParameter)
						f.receiverTypeParams.elements = append(f.receiverTypeParams.elements, typeParameter)
					}
				}
			case *Interface:
				f.typeParams.elements = append(f.typeParams.elements, typedOwner.TypeParameters().elements...)
			case *Struct:
				f.typeParams.elements = append(f.typeParams.elements, typedOwner.TypeParameters().elements...)

				for index, typeParamAlias := range f.typeParamAliases {
					if typeParameter, exists := f.typeParams.FindByName(typeParamAlias); exists {
						f.receiverTypeParams.elements = append(f.receiverTypeParams.elements, typeParameter)
						continue
					}

					typeParam := f.typeParams.At(index)

					if typeParam != nil && typeParam.Name() != typeParamAlias {
						typeParameter := &TypeParameter{
							typeParamAlias,
							typeParam.TypeConstraints(),
						}
						f.typeParams.elements = append(f.typeParams.elements, typeParameter)
						f.receiverTypeParams.elements = append(f.receiverTypeParams.elements, typeParameter)
					}
				}
			}

		}

		if f.funcType == nil || f.funcType.TypeParams == nil {
			return
		}

		for _, field := range f.funcType.TypeParams.List {
			for _, fieldName := range field.Names {
				typeParameter := &TypeParameter{
					name: fieldName.Name,
					constraints: &TypeConstraints{
						[]*TypeConstraint{},
					},
				}
				f.typeParams.elements = append(f.typeParams.elements, typeParameter)
			}
		}

		for _, field := range f.funcType.TypeParams.List {
			constraints := make([]*TypeConstraint, 0)
			typ := getTypeFromExpression(field.Type, f.file, f.visitor, nil, f.typeParams)

			if typeSets, isTypeSets := typ.(TypeSets); isTypeSets {
				for _, item := range typeSets {
					if constraint, isConstraint := item.(*TypeConstraint); isConstraint {
						constraints = append(constraints, constraint)
					}
				}
			} else {
				if constraint, isConstraint := typ.(*TypeConstraint); isConstraint {
					constraints = append(constraints, constraint)
				} else {
					constraints = append(constraints, &TypeConstraint{typ: typ})
				}
			}

			for _, fieldName := range field.Names {
				typeParam, exists := f.typeParams.FindByName(fieldName.Name)

				if exists {
					typeParam.constraints.elements = append(typeParam.constraints.elements, constraints...)
				}
			}
		}
	})
}

func (f *Function) loadParams() {
	f.paramsOnce.Do(func() {
		f.loadTypeParams()

		if f.funcType != nil && f.funcType.Params != nil {
			f.params.elements = append(f.params.elements, f.getParameters(f.funcType.Params.List)...)
		}
	})
}

func (f *Function) loadResultValues() {
	f.resultsOnce.Do(func() {
		f.loadTypeParams()

		if f.funcType.Results != nil {
			f.results.elements = append(f.results.elements, f.getResults(f.funcType.Results.List)...)
		}
	})
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
