package marker

import (
	"errors"
	"go/ast"
	"go/token"
	"path/filepath"
)

type Kind int

const (
	AnyObject Kind = iota
	Object
	Array
	Chan
	Map
	Ptr
	Interface
	Struct
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
	Id   string
	Name string
	Path string
}

type Import struct {
	Name          string
	Path          string
	Position      Position
	RawImportSpec *ast.ImportSpec
}

type File struct {
	Name     string
	FullPath string
	Package  PackageInfo
	Imports  []Import
	Markers  MarkerValues

	Functions      []Function
	StructTypes    []StructType
	InterfaceTypes []InterfaceType
	RawFile        *ast.File
}

type AnyObjectType struct {
}

func (typ AnyObjectType) Kind() Kind {
	return AnyObject
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

type Function struct {
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

func EachFile(collector *Collector, pkgs []*Package, callback FileCallback) {

	if collector == nil {
		callback(nil, errors.New("collector cannot be nil"))
		return
	}

	if pkgs == nil {
		callback(nil, errors.New("pkgs(packages) cannot be nil"))
		return
	}

	var filesInfoMap = make(map[*ast.File]*File)
	var errs []error

	for _, pkg := range pkgs {
		markers, err := collector.Collect(pkg)

		if err != nil {
			errs = append(errs, err.(ErrorList)...)
			continue
		}

		infoMap := eachPackage(pkg, markers)

		for file, fileInfo := range infoMap {
			filesInfoMap[file] = fileInfo
		}
	}

	if errs != nil {
		callback(nil, NewErrorList(errs))
		return
	}

	for _, fileInfo := range filesInfoMap {
		callback(fileInfo, nil)
	}
}

func eachPackage(pkg *Package, markers map[ast.Node]MarkerValues) map[*ast.File]*File {
	var fileInfoMap = make(map[*ast.File]*File)
	var structMethods = make([]Method, 0)

	visitFiles(pkg, func(file *ast.File) {
		fileInfo, ok := fileInfoMap[file]

		if ok {
			return
		}

		position := pkg.Fset.Position(file.Pos())
		fileFullPath := position.Filename

		fileInfo = &File{
			Name:     filepath.Base(fileFullPath),
			FullPath: fileFullPath,
			Package: PackageInfo{
				Id:   pkg.ID,
				Name: file.Name.Name,
				Path: pkg.PkgPath,
			},
			Imports:        getFileImports(pkg.Fset, file),
			Markers:        markers[file],
			Functions:      make([]Function, 0),
			StructTypes:    make([]StructType, 0),
			InterfaceTypes: make([]InterfaceType, 0),
			RawFile:        file,
		}

		fileInfoMap[file] = fileInfo
	})

	visitTypeElements(pkg, func(file *ast.File, decl *ast.GenDecl, spec *ast.TypeSpec) {

		fileInfo, ok := fileInfoMap[file]

		if !ok {
			return
		}

		typ := getType(pkg.Fset, fileInfo, file, decl, spec, markers)

		switch spec.Type.(type) {
		case *ast.InterfaceType:
			fileInfo.InterfaceTypes = append(fileInfo.InterfaceTypes, typ.(InterfaceType))
		case *ast.StructType:
			fileInfo.StructTypes = append(fileInfo.StructTypes, typ.(StructType))
		}

	})

	visitFunctions(pkg, func(file *ast.File, decl *ast.FuncDecl, funcType *ast.FuncType) {

		fileInfo, ok := fileInfoMap[file]

		if !ok {
			return
		}

		// If Recv is nil, it is a function, not a method
		if decl.Recv == nil {
			function := getFunction(pkg.Fset, fileInfo, file, decl, funcType, markers)
			fileInfo.Functions = append(fileInfo.Functions, function)
		} else {
			method := getStructMethod(pkg.Fset, fileInfo, file, decl, funcType, markers)
			structMethods = append(structMethods, method)
		}

	})

	resolveStructMethods(fileInfoMap, structMethods)

	return fileInfoMap
}

func resolveStructMethods(fileInfoMap map[*ast.File]*File, structMethods []Method) {

	for _, structMethod := range structMethods {
		resolveStructMethod(fileInfoMap, structMethod)
	}

}

func resolveStructMethod(fileInfoMap map[*ast.File]*File, structMethod Method) {
	for file, fileInfo := range fileInfoMap {

		for structIndex, structType := range fileInfo.StructTypes {

			receiverType := structMethod.Receiver.Type
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

			if file.Name.Name == fileInfo.Package.Name && structType.Name == receiverTypeName {
				fileInfo.StructTypes[structIndex].Methods = append(fileInfo.StructTypes[structIndex].Methods, structMethod)
				return
			}

		}
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

func getFunction(fileSet *token.FileSet,
	fileInfo *File,
	file *ast.File,
	decl *ast.FuncDecl,
	funcType *ast.FuncType,
	markers map[ast.Node]MarkerValues) Function {

	function := &Function{
		Name:        decl.Name.Name,
		Position:    getPosition(fileSet, funcType.Pos()),
		Markers:     markers[decl],
		File:        fileInfo,
		RawFile:     file,
		RawFuncDecl: decl,
		RawFuncType: funcType,
	}

	if funcType.Params != nil {
		function.Parameters = getFieldTypesInfo(fileSet, funcType.Params.List, markers)
	}

	if funcType.Results != nil {
		function.ReturnValues = getFieldTypesInfo(fileSet, funcType.Results.List, markers)
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

		structType.Fields = getStructFields(fileSet, file, specType, markers)
		typ = structType
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
			method.Parameters = getFieldTypesInfo(fileSet, methodInfo.Type.(*ast.FuncType).Params.List, markers)
		}

		if methodInfo.Type.(*ast.FuncType).Results != nil {
			method.ReturnValues = getFieldTypesInfo(fileSet, methodInfo.Type.(*ast.FuncType).Results.List, markers)
		}

		methods = append(methods, *method)

	}

	return methods
}

func getStructFields(fileSet *token.FileSet,
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
				Type:     getTypeFromExpression(fileSet, fieldTypeInfo.Type, markers),
				RawFile:  file,
				RawField: fieldTypeInfo,
			}

			fields = append(fields, *field)
		}

	}

	return fields
}

func getStructMethod(fileSet *token.FileSet,
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
		method.Parameters = getFieldTypesInfo(fileSet, funcType.Params.List, markers)
	}

	if funcType.Results != nil {
		method.ReturnValues = getFieldTypesInfo(fileSet, funcType.Results.List, markers)
	}

	// Receiver
	receiver := decl.Recv.List[0]
	receiverType := getTypeFromExpression(fileSet, receiver.Type, markers)

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

func getFieldTypesInfo(tokenFileSet *token.FileSet, fieldList []*ast.Field, markers map[ast.Node]MarkerValues) []TypeInfo {
	types := make([]TypeInfo, 0)

	for _, field := range fieldList {

		if field.Names == nil {
			typeInfo := &TypeInfo{
				RawField: field,
				Type:     getTypeFromExpression(tokenFileSet, field.Type, markers),
				Markers:  markers[field],
			}

			types = append(types, *typeInfo)
			continue
		}

		for _, name := range field.Names {
			typeInfo := &TypeInfo{
				Name:     name.Name,
				RawField: field,
				Type:     getTypeFromExpression(tokenFileSet, field.Type, markers),
				Markers:  markers[field],
			}

			types = append(types, *typeInfo)
		}
	}

	return types
}

func getTypeFromExpression(tokenFileSet *token.FileSet, expression ast.Expr, markers map[ast.Node]MarkerValues) Type {

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
			Typ: getTypeFromExpression(tokenFileSet, result.X, markers),
		}
	case *ast.ArrayType:
		return &ArrayType{
			ItemType: getTypeFromExpression(tokenFileSet, result.Elt, markers),
		}
	case *ast.MapType:
		return &DictionaryType{
			KeyType:   getTypeFromExpression(tokenFileSet, result.Key, markers),
			ValueType: getTypeFromExpression(tokenFileSet, result.Value, markers),
		}
	case *ast.ChanType:
		chanTyp := &ChanType{
			Typ: getTypeFromExpression(tokenFileSet, result.Value, markers),
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

		fieldTypeInfoList := getFieldTypesInfo(tokenFileSet, result.Fields.List, markers)

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
		return &AnyObjectType{}
	}

	panic("Unreachable code!")
}
