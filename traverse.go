package marker

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"
)

type Kind int

const (
	Anyx Kind = iota
	Objectx
	Arrayx
	Chanx
	Mapx
	Ptr
	Variadic
	Functionx
	Interfacex
	Structx
	UserDefined
)

type Type interface {
	Kind() Kind
}

type ValueType struct {
	ImportName string
	Name       string
}

type ConstValue struct {
	Name         string
	Value        string
	IsExported   bool
	Type         *ValueType
	Position     Position
	RawFile      *ast.File
	RawValueSpec *ast.ValueSpec
}

type TypeInfo struct {
	Name     string
	Type     Type
	Markers  MarkerValues
	RawField *ast.Field
}

type FileCallback func(file *File, err error)

type File struct {
	Name          string
	FullPath      string
	Package       *Package
	Imports       []Import
	Consts        []ConstValue
	Markers       MarkerValues
	ImportMarkers []MarkerValues

	FunctionTypes    []FunctionType
	StructTypes      []StructType
	InterfaceTypes   []InterfaceType
	UserDefinedTypes []UserDefinedType
	RawFile          *ast.File
}

type AnyObject struct {
}

func (typ AnyObject) Kind() Kind {
	return Anyx
}

type ObjectType struct {
	ImportName string
	Name       string
}

func (typ ObjectType) Kind() Kind {
	return Objectx
}

type Pointerx struct {
	Base Type
}

func (typ Pointerx) Kind() Kind {
	return Ptr
}

type ArrayType struct {
	Elem Type
}

func (typ ArrayType) Kind() Kind {
	return Arrayx
}

type DictionaryType struct {
	Key  Type
	Elem Type
}

func (typ DictionaryType) Kind() Kind {
	return Mapx
}

type ChanType struct {
	Typ       Type
	Direction ChanDirection
}

func (typ ChanType) Kind() Kind {
	return Chanx
}

type FunctionType struct {
	Name         string
	IsExported   bool
	Position     Position
	Markers      MarkerValues
	Parameters   []TypeInfo
	ReturnValues []TypeInfo
	File         *File
	RawFile      *ast.File
	RawFuncDecl  *ast.FuncDecl
	RawFuncType  *ast.FuncType
}

func (typ FunctionType) Kind() Kind {
	return Functionx
}

type Fieldx struct {
	Name       string
	IsExported bool
	IsEmbedded bool
	Position   Position
	Markers    MarkerValues
	Type       Type
	File       *File
	RawFile    *ast.File
	RawField   *ast.Field
}

type Method struct {
	Name         string
	IsExported   bool
	Position     Position
	Markers      MarkerValues
	Receiver     TypeInfo
	Parameters   []TypeInfo
	ReturnValues []TypeInfo
	File         *File
	RawFile      *ast.File
	RawField     *ast.Field
	RawFuncDecl  *ast.FuncDecl
	RawFuncType  *ast.FuncType
}

type StructType struct {
	Name        string
	IsExported  bool
	Position    Position
	Markers     MarkerValues
	Fields      []Fieldx
	Methods     []Method
	File        *File
	RawFile     *ast.File
	RawGenDecl  *ast.GenDecl
	RawTypeSpec *ast.TypeSpec
	namedType   *types.Named
}

func (structType StructType) Kind() Kind {
	return Structx
}

func (structType StructType) Embeddeds() []Type {
	return nil
}

func (structType StructType) Implements(typ InterfaceType) bool {

	if types.Implements(structType.namedType, typ.interfaceType) {
		return true
	}

	pointerType := types.NewPointer(structType.namedType)

	if types.Implements(pointerType, typ.interfaceType) {
		return true
	}

	return false
}

type UserDefinedType struct {
	Name        string
	IsExported  bool
	ActualType  Type
	Position    Position
	Markers     MarkerValues
	Methods     []Method
	File        *File
	RawFile     *ast.File
	RawGenDecl  *ast.GenDecl
	RawTypeSpec *ast.TypeSpec
}

func (typ UserDefinedType) Kind() Kind {
	return UserDefined
}

type InterfaceType struct {
	Name          string
	IsExported    bool
	Position      Position
	Markers       MarkerValues
	AllMethods    []Method
	Methods       []Method
	Embeddeds     []Type
	File          *File
	RawFile       *ast.File
	RawGenDecl    *ast.GenDecl
	RawTypeSpec   *ast.TypeSpec
	interfaceType *types.Interface
}

func (typ InterfaceType) Kind() Kind {
	return Interfacex
}

type AnonymousStructType struct {
	Fields []Fieldx
}

func (typ AnonymousStructType) Kind() Kind {
	return Structx
}

type VariadicType struct {
	ItemType Type
}

func (typ VariadicType) Kind() Kind {
	return Variadic
}

func EachFile(collector *Collector, pkgs []*Package, callback FileCallback) {
	if collector == nil {
		callback(nil, errors.New("collector cannot be nil"))
	}

	if pkgs == nil {
		callback(nil, errors.New("pkgs(packages) cannot be nil"))
	}

	var fileMap = make(map[*ast.File]*File)
	var errs []error

	for _, pkg := range pkgs {
		markers, err := collector.Collect(pkg)

		if err != nil {
			errs = append(errs, err.(ErrorList)...)
			continue
		}

		fileNodeMap := eachPackage(pkg, markers)

		for fileNode, file := range fileNodeMap {
			fileMap[fileNode] = file
		}
	}

	if errs != nil {
		callback(nil, NewErrorList(errs))
	}

	for _, file := range fileMap {
		callback(file, nil)
	}
}

func eachPackage(pkg *Package, markers map[ast.Node]MarkerValues) map[*ast.File]*File {
	var fileNodeMap = make(map[*ast.File]*File)
	var methods = make([]Method, 0)

	visitPackageFiles(pkg, func(file *ast.File) {
		_, ok := fileNodeMap[file]

		if ok {
			return
		}

		fileNodeMap[file] = getFile(pkg, file, markers)
	}, func(file *ast.File, decl *ast.GenDecl) {
		fileInfo, ok := fileNodeMap[file]

		if !ok {
			return
		}

		if markerValues, ok := markers[decl]; ok {
			fileInfo.ImportMarkers = append(fileInfo.ImportMarkers, markerValues)
		}
	}, func(file *ast.File, decl *ast.GenDecl) {
		fileInfo, ok := fileNodeMap[file]

		if !ok {
			return
		}

		constValues := getConstValues(pkg.Fset, file, decl.Specs)

		if constValues != nil {
			fileInfo.Consts = append(fileInfo.Consts, constValues...)
		}
	}, func(file *ast.File, decl *ast.GenDecl, spec *ast.TypeSpec) {
		fileInfo, ok := fileNodeMap[file]

		if !ok {
			return
		}

		typ := getType(pkg, pkg.Fset, fileInfo, file, decl, spec, markers)

		switch spec.Type.(type) {
		case *ast.InterfaceType:
			fileInfo.InterfaceTypes = append(fileInfo.InterfaceTypes, typ.(InterfaceType))
		case *ast.StructType:
			fileInfo.StructTypes = append(fileInfo.StructTypes, typ.(StructType))
		default:
			fileInfo.UserDefinedTypes = append(fileInfo.UserDefinedTypes, typ.(UserDefinedType))
		}
	}, func(file *ast.File, decl *ast.FuncDecl, funcType *ast.FuncType) {
		fileInfo, ok := fileNodeMap[file]

		if !ok {
			return
		}

		// If Recv is nil, it is a function, not a method
		if decl.Recv == nil {
			functionType := getFunctionType(pkg.Fset, fileInfo, file, decl, funcType, markers)
			fileInfo.FunctionTypes = append(fileInfo.FunctionTypes, functionType)
		} else {
			method := getMethod(pkg.Fset, fileInfo, file, decl, funcType, markers)
			methods = append(methods, method)
		}
	})

	resolveMethods(fileNodeMap, methods)

	return fileNodeMap
}

func resolveMethods(fileInfoMap map[*ast.File]*File, methods []Method) {

	for _, method := range methods {
		resolveMethod(fileInfoMap, method)
	}

}

func resolveMethod(fileInfoMap map[*ast.File]*File, method Method) {

	receiverType := method.Receiver.Type
	var receiverTypeName string

	switch typed := receiverType.(type) {
	case *Pointerx:
		objectType := typed.Base.(*ObjectType)

		if objectType.ImportName != "" {
			return
		}

		receiverTypeName = objectType.Name

	case *ObjectType:
		if typed.ImportName != "" {
			return
		}

		receiverTypeName = typed.Name
	}

	for file, fileInfo := range fileInfoMap {

		for index, structType := range fileInfo.StructTypes {

			if file.Name.Name == fileInfo.Package.Name && structType.Name == receiverTypeName {
				fileInfo.StructTypes[index].Methods = append(fileInfo.StructTypes[index].Methods, method)
				return
			}

		}

		for index, userDefinedType := range fileInfo.UserDefinedTypes {

			if file.Name.Name == fileInfo.Package.Name && userDefinedType.Name == receiverTypeName {
				fileInfo.UserDefinedTypes[index].Methods = append(fileInfo.UserDefinedTypes[index].Methods, method)
				return
			}

		}
	}
}

func getFile(pkg *Package, file *ast.File, markers map[ast.Node]MarkerValues) *File {
	position := pkg.Fset.Position(file.Pos())
	fileFullPath := position.Filename

	return &File{
		Name:           filepath.Base(fileFullPath),
		FullPath:       fileFullPath,
		Package:        pkg,
		Imports:        getFileImports(pkg.Fset, file),
		Consts:         make([]ConstValue, 0),
		Markers:        markers[file],
		ImportMarkers:  make([]MarkerValues, 0),
		FunctionTypes:  make([]FunctionType, 0),
		StructTypes:    make([]StructType, 0),
		InterfaceTypes: make([]InterfaceType, 0),
		RawFile:        file,
	}
}

func getFileImports(fileSet *token.FileSet, file *ast.File) []Import {
	imports := make([]Import, 0)

	for _, importInfo := range file.Imports {
		importPosition := fileSet.Position(importInfo.Pos())
		importName := ""

		if importInfo.Name != nil {
			importName = importInfo.Name.Name
		}

		imports = append(imports, Import{
			name: importName,
			path: importInfo.Path.Value[1 : len(importInfo.Path.Value)-1],
			position: Position{
				importPosition.Line,
				importPosition.Column,
			},
			rawImportSpec: importInfo,
		})
	}

	return imports
}

func getFunctionType(fileSet *token.FileSet,
	fileInfo *File,
	file *ast.File,
	decl *ast.FuncDecl,
	funcType *ast.FuncType,
	markers map[ast.Node]MarkerValues) FunctionType {

	function := &FunctionType{
		Position:    getPosition(fileSet, funcType.Pos()),
		File:        fileInfo,
		RawFile:     file,
		RawFuncDecl: decl,
		RawFuncType: funcType,
	}

	if decl != nil {
		function.Name = decl.Name.Name
		function.IsExported = ast.IsExported(decl.Name.Name)
		function.Markers = markers[decl]
	}

	if funcType.Params != nil {
		function.Parameters = getFieldTypesInfo(fileSet, fileInfo, file, funcType.Params.List, markers)
	}

	if funcType.Results != nil {
		function.ReturnValues = getFieldTypesInfo(fileSet, fileInfo, file, funcType.Results.List, markers)
	}

	return *function
}

func getConstValues(fileSet *token.FileSet,
	file *ast.File,
	specs []ast.Spec) []ConstValue {
	if specs == nil {
		return nil
	}

	constValues := make([]ConstValue, 0)

	var previousValueType *ValueType
	for _, spec := range specs {
		valueSpec := spec.(*ast.ValueSpec)
		constValue, inferredTypeName := getConstValue(fileSet, file, valueSpec)

		if valueSpec.Type == nil && constValue.Type == nil && inferredTypeName != "" {
			constValue.Type = &ValueType{
				Name: inferredTypeName,
			}
		} else if constValue.Type != nil {
			previousValueType = constValue.Type
		} else {
			constValue.Type = previousValueType
		}

		constValues = append(constValues, *constValue)

	}

	return constValues
}

func getConstValue(fileSet *token.FileSet, file *ast.File, spec *ast.ValueSpec) (*ConstValue, string) {
	constValue := &ConstValue{
		Name:         spec.Names[0].Name,
		Type:         getConstValueType(spec.Type),
		IsExported:   ast.IsExported(spec.Names[0].Name),
		Position:     getPosition(fileSet, spec.Pos()),
		RawFile:      file,
		RawValueSpec: spec,
	}

	inferredTypeName := ""

	if spec.Values != nil && len(spec.Values) > 0 {
		switch typedValue := spec.Values[0].(type) {
		case *ast.Ident:
			value := typedValue.Name

			if !strings.HasPrefix(value, "\"") && !strings.HasSuffix(value, "\"") {
				if value == "true" || value == "false" {
					inferredTypeName = "bool"
					constValue.Value = value
				}
			}

		case *ast.BasicLit:
			constValue.Value = typedValue.Value

			switch typedValue.Kind {
			case token.STRING:
				inferredTypeName = "string"
			case token.INT:
				inferredTypeName = "int"
			case token.FLOAT:
				inferredTypeName = "float"
			case token.CHAR:
				inferredTypeName = "char"
			}
		}
	}

	return constValue, inferredTypeName
}

func getConstValueType(typ ast.Expr) *ValueType {
	switch typedExpr := typ.(type) {
	case *ast.Ident:
		return &ValueType{
			Name: typedExpr.Name,
		}
	case *ast.SelectorExpr:
		return &ValueType{
			ImportName: typedExpr.X.(*ast.Ident).Name,
			Name:       typedExpr.Sel.Name,
		}
	}

	return nil
}

func getType(pkg *Package, fileSet *token.FileSet,
	fileInfo *File,
	file *ast.File,
	decl *ast.GenDecl,
	spec *ast.TypeSpec,
	markers map[ast.Node]MarkerValues) Type {

	var typ Type

	switch specType := spec.Type.(type) {
	case *ast.InterfaceType:
		y := pkg.Types.Scope().Lookup(spec.Name.Name)

		if y != nil {

		}

		interfaceType := InterfaceType{
			Name:          spec.Name.Name,
			IsExported:    ast.IsExported(spec.Name.Name),
			Position:      getPosition(fileSet, spec.Pos()),
			Markers:       markers[spec],
			File:          fileInfo,
			RawFile:       file,
			RawGenDecl:    decl,
			RawTypeSpec:   spec,
			interfaceType: pkg.Types.Scope().Lookup(spec.Name.Name).Type().Underlying().(*types.Interface),
		}

		interfaceType.Methods = getInterfaceMethods(pkg, fileSet, fileInfo, file, specType, markers)
		typ = interfaceType
	case *ast.StructType:
		y := pkg.Types.Scope().Lookup(spec.Name.Name)

		if y != nil {

		}

		structType := StructType{
			Name:        spec.Name.Name,
			IsExported:  ast.IsExported(spec.Name.Name),
			Position:    getPosition(fileSet, spec.Pos()),
			Markers:     markers[spec],
			File:        fileInfo,
			RawFile:     file,
			RawGenDecl:  decl,
			RawTypeSpec: spec,
			namedType:   pkg.Types.Scope().Lookup(spec.Name.Name).Type().(*types.Named),
		}

		structType.Fields = getStructFields(fileSet, fileInfo, file, specType, markers)
		typ = structType
	case *ast.Ident, *ast.StarExpr, *ast.SelectorExpr, *ast.MapType, *ast.ArrayType, *ast.ChanType, *ast.FuncType:
		typ = UserDefinedType{
			Name:        spec.Name.Name,
			IsExported:  ast.IsExported(spec.Name.Name),
			Position:    getPosition(fileSet, spec.Pos()),
			ActualType:  getTypeFromExpression(fileSet, fileInfo, file, spec.Type, markers),
			Markers:     markers[spec],
			File:        fileInfo,
			RawFile:     file,
			RawGenDecl:  decl,
			RawTypeSpec: spec,
		}
	}

	return typ
}

func getInterfaceMethods(pkg *Package, fileSet *token.FileSet,
	fileInfo *File,
	file *ast.File,
	specType *ast.InterfaceType,
	markers map[ast.Node]MarkerValues) []Method {

	methods := make([]Method, 0)

	for _, methodInfo := range specType.Methods.List {
		y := pkg.Types.Scope().Lookup(methodInfo.Names[0].Name)

		if y != nil {

		}

		method := &Method{
			Name:        methodInfo.Names[0].Name,
			IsExported:  ast.IsExported(methodInfo.Names[0].Name),
			Position:    getPosition(fileSet, methodInfo.Pos()),
			Markers:     markers[methodInfo],
			File:        fileInfo,
			RawFile:     file,
			RawFuncType: methodInfo.Type.(*ast.FuncType),
			RawField:    methodInfo,
		}

		if methodInfo.Type.(*ast.FuncType).Params != nil {
			method.Parameters = getFieldTypesInfo(fileSet, fileInfo, file, methodInfo.Type.(*ast.FuncType).Params.List, markers)
		}

		if methodInfo.Type.(*ast.FuncType).Results != nil {
			method.ReturnValues = getFieldTypesInfo(fileSet, fileInfo, file, methodInfo.Type.(*ast.FuncType).Results.List, markers)
		}

		methods = append(methods, *method)

	}

	return methods
}

func getStructFields(fileSet *token.FileSet,
	fileInfo *File,
	file *ast.File,
	specType *ast.StructType,
	markers map[ast.Node]MarkerValues) []Fieldx {

	fields := make([]Fieldx, 0)

	for _, fieldTypeInfo := range specType.Fields.List {

		if fieldTypeInfo.Names == nil {
			field := &Fieldx{
				IsEmbedded: true,
				Position:   getPosition(fileSet, fieldTypeInfo.Type.Pos()),
				Markers:    markers[fieldTypeInfo],
				Type:       getTypeFromExpression(fileSet, fileInfo, file, fieldTypeInfo.Type, markers),
				File:       fileInfo,
				RawFile:    file,
				RawField:   fieldTypeInfo,
			}
			fields = append(fields, *field)
			continue
		}

		for _, fieldName := range fieldTypeInfo.Names {
			field := &Fieldx{
				Name:       fieldName.Name,
				IsExported: ast.IsExported(fieldName.Name),
				IsEmbedded: false,
				Position:   getPosition(fileSet, fieldName.Pos()),
				Markers:    markers[fieldTypeInfo],
				Type:       getTypeFromExpression(fileSet, fileInfo, file, fieldTypeInfo.Type, markers),
				File:       fileInfo,
				RawFile:    file,
				RawField:   fieldTypeInfo,
			}

			fields = append(fields, *field)
		}

	}

	return fields
}

func getMethod(fileSet *token.FileSet,
	fileInfo *File,
	file *ast.File,
	decl *ast.FuncDecl,
	funcType *ast.FuncType,
	markers map[ast.Node]MarkerValues) Method {

	method := &Method{
		Name:        decl.Name.Name,
		IsExported:  ast.IsExported(decl.Name.Name),
		Position:    getPosition(fileSet, funcType.Pos()),
		Markers:     markers[decl],
		File:        fileInfo,
		RawFile:     file,
		RawFuncDecl: decl,
		RawFuncType: funcType,
	}

	if funcType.Params != nil {
		method.Parameters = getFieldTypesInfo(fileSet, fileInfo, file, funcType.Params.List, markers)
	}

	if funcType.Results != nil {
		method.ReturnValues = getFieldTypesInfo(fileSet, fileInfo, file, funcType.Results.List, markers)
	}

	// Receiver
	receiver := decl.Recv.List[0]
	receiverType := getTypeFromExpression(fileSet, fileInfo, file, receiver.Type, markers)

	method.Receiver = TypeInfo{
		Name:     receiver.Names[0].Name,
		Type:     receiverType,
		Markers:  nil,
		RawField: nil,
	}

	return *method
}

func getPosition(tokenFileSet *token.FileSet, pos token.Pos) Position {
	position := tokenFileSet.Position(pos)
	return Position{
		Line:   position.Line,
		Column: position.Column,
	}
}

func getFieldTypesInfo(tokenFileSet *token.FileSet,
	fileInfo *File,
	file *ast.File,
	fieldList []*ast.Field, markers map[ast.Node]MarkerValues) []TypeInfo {
	types := make([]TypeInfo, 0)

	for _, field := range fieldList {

		if field.Names == nil {
			typeInfo := &TypeInfo{
				RawField: field,
				Type:     getTypeFromExpression(tokenFileSet, fileInfo, file, field.Type, markers),
				Markers:  markers[field],
			}

			types = append(types, *typeInfo)
			continue
		}

		for _, name := range field.Names {
			typeInfo := &TypeInfo{
				Name:     name.Name,
				RawField: field,
				Type:     getTypeFromExpression(tokenFileSet, fileInfo, file, field.Type, markers),
				Markers:  markers[field],
			}

			types = append(types, *typeInfo)
		}
	}

	return types
}

func getTypeFromExpression(tokenFileSet *token.FileSet,
	fileInfo *File,
	file *ast.File,
	expression ast.Expr,
	markers map[ast.Node]MarkerValues) Type {

	switch result := expression.(type) {
	case *ast.Ident:
		return &ObjectType{
			Name: result.Name,
		}
	case *ast.SelectorExpr:
		return &ObjectType{
			ImportName: result.X.(*ast.Ident).Name,
			Name:       result.Sel.Name,
		}
	case *ast.StarExpr:
		return &Pointerx{
			Base: getTypeFromExpression(tokenFileSet, fileInfo, file, result.X, markers),
		}
	case *ast.ArrayType:
		return &ArrayType{
			Elem: getTypeFromExpression(tokenFileSet, fileInfo, file, result.Elt, markers),
		}
	case *ast.MapType:
		return &DictionaryType{
			Key:  getTypeFromExpression(tokenFileSet, fileInfo, file, result.Key, markers),
			Elem: getTypeFromExpression(tokenFileSet, fileInfo, file, result.Value, markers),
		}
	case *ast.ChanType:
		chanTyp := &ChanType{
			Typ: getTypeFromExpression(tokenFileSet, fileInfo, file, result.Value, markers),
		}

		if result.Dir&ast.SEND == ast.SEND {
			chanTyp.Direction |= SEND
		}

		if result.Dir&ast.RECV == ast.RECV {
			chanTyp.Direction |= RECEIVE
		}

		return chanTyp
	case *ast.StructType:
		anonymousStructType := &AnonymousStructType{}

		fieldTypeInfoList := getFieldTypesInfo(tokenFileSet, fileInfo, file, result.Fields.List, markers)

		for _, fieldTypeInfo := range fieldTypeInfoList {

			field := &Fieldx{
				Name:       fieldTypeInfo.Name,
				IsExported: ast.IsExported(fieldTypeInfo.Name),
				Position:   getPosition(tokenFileSet, fieldTypeInfo.RawField.Pos()),
				Markers:    fieldTypeInfo.Markers,
				Type:       fieldTypeInfo.Type,
				File:       fileInfo,
				RawFile:    file,
				RawField:   fieldTypeInfo.RawField,
			}

			anonymousStructType.Fields = append(anonymousStructType.Fields, *field)

		}

		return anonymousStructType
	case *ast.InterfaceType:
		return &AnyObject{}
	case *ast.FuncType:
		return getFunctionType(tokenFileSet, fileInfo, file, nil, result, markers)
	case *ast.Ellipsis:
		return &VariadicType{
			ItemType: getTypeFromExpression(tokenFileSet, fileInfo, file, result.Elt, markers),
		}
	}

	panic("Unreachable code!")
}
