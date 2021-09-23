package marker

import (
	"errors"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
	"path/filepath"
)

type Kind int

const (
	AnyKind Kind = iota
	Object
	Array
	Chan
	Map
	Ptr
	Variadic
	Function
	Interface
	Struct
	UserDefined
)

type Type interface {
	Kind() Kind
}

type TypeInfo struct {
	Name     string
	Type     Type
	Markers  MarkerValues
	RawField *ast.Field
}

type FileCallback func(file *File, error error)

type Position struct {
	Line   int
	Column int
}

type PackageInfo struct {
	Id         string
	Name       string
	Path       string
	ModuleInfo *ModuleInfo
	RawPackage *packages.Package
}

type ModuleInfo struct {
	Path      string
	Version   string
	Main      bool
	Indirect  bool
	Dir       string
	GoMod     string
	GoVersion string
	RawModule *packages.Module
}

type Import struct {
	Name          string
	Path          string
	Position      Position
	RawImportSpec *ast.ImportSpec
}

type File struct {
	Name          string
	FullPath      string
	Package       PackageInfo
	Imports       []Import
	Markers       MarkerValues
	ImportMarkers []MarkerValues

	FunctionTypes    []FunctionType
	StructTypes      []StructType
	InterfaceTypes   []InterfaceType
	UserDefinedTypes []UserDefinedType
	RawFile          *ast.File
}

type AnyKindType struct {
}

func (typ AnyKindType) Kind() Kind {
	return AnyKind
}

type ObjectType struct {
	ImportName string
	Name       string
}

func (typ ObjectType) Kind() Kind {
	return Object
}

type PointerType struct {
	Typ Type
}

func (typ PointerType) Kind() Kind {
	return Ptr
}

type ArrayType struct {
	ItemType Type
}

func (typ ArrayType) Kind() Kind {
	return Array
}

type DictionaryType struct {
	KeyType   Type
	ValueType Type
}

func (typ DictionaryType) Kind() Kind {
	return Map
}

type ChanDirection int

const (
	SEND ChanDirection = 1 << iota
	RECEIVE
)

type ChanType struct {
	Typ       Type
	Direction ChanDirection
}

func (typ ChanType) Kind() Kind {
	return Chan
}

type FunctionType struct {
	Name         string
	Position     Position
	Markers      MarkerValues
	Parameters   []TypeInfo
	ReturnValues []TypeInfo
	File         *File
	RawFile      *ast.File
	RawFuncDecl  *ast.FuncDecl
	RawFuncType  *ast.FuncType
}

func (function FunctionType) Kind() Kind {
	return Function
}

type Field struct {
	Name     string
	Position Position
	Markers  MarkerValues
	Type     Type
	RawFile  *ast.File
	RawField *ast.Field
}

type Method struct {
	Name         string
	Position     Position
	Markers      MarkerValues
	Receiver     *TypeInfo
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
	Position    Position
	Markers     MarkerValues
	Fields      []Field
	Methods     []Method
	File        *File
	RawFile     *ast.File
	RawGenDecl  *ast.GenDecl
	RawTypeSpec *ast.TypeSpec
}

func (typ StructType) Kind() Kind {
	return Struct
}

type UserDefinedType struct {
	Name        string
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
	Name        string
	Position    Position
	Markers     MarkerValues
	Methods     []Method
	File        *File
	RawFile     *ast.File
	RawGenDecl  *ast.GenDecl
	RawTypeSpec *ast.TypeSpec
}

func (typ InterfaceType) Kind() Kind {
	return Interface
}

type AnonymousStructType struct {
	Markers MarkerValues
	Fields  []Field
}

func (typ AnonymousStructType) Kind() Kind {
	return Struct
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
		return
	}

	if pkgs == nil {
		callback(nil, errors.New("pkgs(packages) cannot be nil"))
		return
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
		return
	}

	for _, file := range fileMap {
		callback(file, nil)
	}
}

func eachPackage(pkg *Package, markers map[ast.Node]MarkerValues) map[*ast.File]*File {
	var fileNodeMap = make(map[*ast.File]*File)
	var methods = make([]Method, 0)

	visitFiles(pkg, func(file *ast.File) {
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
	}, func(file *ast.File, decl *ast.GenDecl, spec *ast.TypeSpec) {
		fileInfo, ok := fileNodeMap[file]

		if !ok {
			return
		}

		typ := getType(pkg.Fset, fileInfo, file, decl, spec, markers)

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
	case *PointerType:
		objectType := typed.Typ.(*ObjectType)

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

	packageInfo := PackageInfo{
		Id:         pkg.ID,
		Name:       file.Name.Name,
		Path:       pkg.PkgPath,
		RawPackage: pkg.Package,
	}

	if pkg.Module != nil {
		packageInfo.ModuleInfo = &ModuleInfo{
			Path:      pkg.Module.Path,
			Version:   pkg.Module.Version,
			Main:      pkg.Module.Main,
			Indirect:  pkg.Module.Indirect,
			Dir:       pkg.Module.Dir,
			GoMod:     pkg.Module.GoMod,
			GoVersion: pkg.Module.GoVersion,
			RawModule: pkg.Module,
		}
	}

	return &File{
		Name:           filepath.Base(fileFullPath),
		FullPath:       fileFullPath,
		Package:        packageInfo,
		Imports:        getFileImports(pkg.Fset, file),
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
			Name: importName,
			Path: importInfo.Path.Value[1 : len(importInfo.Path.Value)-1],
			Position: Position{
				importPosition.Line,
				importPosition.Column,
			},
			RawImportSpec: importInfo,
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

func getType(fileSet *token.FileSet,
	fileInfo *File,
	file *ast.File,
	decl *ast.GenDecl,
	spec *ast.TypeSpec,
	markers map[ast.Node]MarkerValues) Type {

	var typ Type

	switch specType := spec.Type.(type) {
	case *ast.InterfaceType:
		interfaceType := InterfaceType{
			Name:        spec.Name.Name,
			Position:    getPosition(fileSet, spec.Pos()),
			Markers:     markers[spec],
			File:        fileInfo,
			RawFile:     file,
			RawGenDecl:  decl,
			RawTypeSpec: spec,
		}

		interfaceType.Methods = getInterfaceMethods(fileSet, fileInfo, file, specType, markers)
		typ = interfaceType
	case *ast.StructType:
		structType := StructType{
			Name:        spec.Name.Name,
			Position:    getPosition(fileSet, spec.Pos()),
			Markers:     markers[spec],
			File:        fileInfo,
			RawFile:     file,
			RawGenDecl:  decl,
			RawTypeSpec: spec,
		}

		structType.Fields = getStructFields(fileSet, fileInfo, file, specType, markers)
		typ = structType
	case *ast.Ident:
		typ = UserDefinedType{
			Name:        spec.Name.Name,
			Position:    getPosition(fileSet, spec.Pos()),
			ActualType:  getTypeFromExpression(fileSet, fileInfo, file, spec.Type, markers),
			Markers:     markers[spec],
			File:        fileInfo,
			RawFile:     file,
			RawGenDecl:  decl,
			RawTypeSpec: spec,
		}

	case *ast.SelectorExpr:
		typ = UserDefinedType{
			Name:        spec.Name.Name,
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

func getInterfaceMethods(fileSet *token.FileSet,
	fileInfo *File,
	file *ast.File,
	specType *ast.InterfaceType,
	markers map[ast.Node]MarkerValues) []Method {

	methods := make([]Method, 0)

	for _, methodInfo := range specType.Methods.List {

		method := &Method{
			Name:        methodInfo.Names[0].Name,
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
	markers map[ast.Node]MarkerValues) []Field {

	fields := make([]Field, 0)

	for _, fieldTypeInfo := range specType.Fields.List {

		for _, fieldName := range fieldTypeInfo.Names {
			field := &Field{
				Name:     fieldName.Name,
				Position: getPosition(fileSet, fieldName.Pos()),
				Markers:  markers[fieldTypeInfo],
				Type:     getTypeFromExpression(fileSet, fileInfo, file, fieldTypeInfo.Type, markers),
				RawFile:  file,
				RawField: fieldTypeInfo,
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
		Position:    getPosition(fileSet, funcType.Pos()),
		Markers:     markers[decl],
		Receiver:    &TypeInfo{},
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

	method.Receiver.Type = receiverType
	method.Receiver.Name = receiver.Names[0].Name

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
		return &PointerType{
			Typ: getTypeFromExpression(tokenFileSet, fileInfo, file, result.X, markers),
		}
	case *ast.ArrayType:
		return &ArrayType{
			ItemType: getTypeFromExpression(tokenFileSet, fileInfo, file, result.Elt, markers),
		}
	case *ast.MapType:
		return &DictionaryType{
			KeyType:   getTypeFromExpression(tokenFileSet, fileInfo, file, result.Key, markers),
			ValueType: getTypeFromExpression(tokenFileSet, fileInfo, file, result.Value, markers),
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

			field := &Field{
				Name:     fieldTypeInfo.Name,
				Position: getPosition(tokenFileSet, fieldTypeInfo.RawField.Pos()),
				Markers:  fieldTypeInfo.Markers,
				Type:     fieldTypeInfo.Type,
				RawField: fieldTypeInfo.RawField,
			}

			anonymousStructType.Fields = append(anonymousStructType.Fields, *field)

		}

		return anonymousStructType
	case *ast.InterfaceType:
		return &AnyKindType{}
	case *ast.FuncType:
		return getFunctionType(tokenFileSet, fileInfo, file, nil, result, markers)
	case *ast.Ellipsis:
		return &VariadicType{
			ItemType: getTypeFromExpression(tokenFileSet, fileInfo, file, result.Elt, markers),
		}
	}

	panic("Unreachable code!")
}
