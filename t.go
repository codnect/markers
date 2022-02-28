package marker

import (
	"errors"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
	"path/filepath"
	"strconv"
)

type Import struct {
	name          string
	path          string
	position      Position
	rawImportSpec *ast.ImportSpec
}

func (i *Import) Name() string {
	return i.name
}

func (i *Import) Path() string {
	return i.path
}

func (i *Import) Position() Position {
	return i.position
}

func (i *Import) RawImportSpec() *ast.ImportSpec {
	return i.rawImportSpec
}

type Imports struct {
	imports []*Import
}

func (i *Imports) Len() int {
	return len(i.imports)
}

func (i *Imports) At(index int) *Import {
	if index >= 0 && index < len(i.imports) {
		return i.imports[index]
	}

	return nil
}

func (i *Imports) FindByName(name string) (*Import, bool) {
	for _, importItem := range i.imports {
		if importItem.name == name {
			return importItem, true
		}
	}

	return nil, false
}

func (i *Imports) FindByPath(path string) (*Import, bool) {
	for _, importItem := range i.imports {
		if importItem.path == path {
			return importItem, true
		}
	}

	return nil, false
}

type Functions struct {
	functions []*Function
}

func (f *Functions) Len() int {
	return len(f.functions)
}

func (f *Functions) At(index int) *Function {
	if index >= 0 && index < len(f.functions) {
		return f.functions[index]
	}

	return nil
}

type Structs struct {
	strutcs []*Struct
}

func (s *Structs) Len() int {
	return len(s.strutcs)
}

func (s *Structs) At(index int) *Struct {
	if index >= 0 && index < len(s.strutcs) {
		return s.strutcs[index]
	}

	return nil
}

type Interfaces struct {
	interfaces []*Interface
}

func (i *Interfaces) Len() int {
	return len(i.interfaces)
}

func (i *Interfaces) At(index int) *Interface {
	if index >= 0 && index < len(i.interfaces) {
		return i.interfaces[index]
	}

	return nil
}

type SourceFile struct {
	name     string
	fullPath string
	pkg      *Package

	allMarkers  MarkerValues
	fileMarkers MarkerValues

	imports       *Imports
	importMarkers []ImportMarker

	functions  *Functions
	structs    *Structs
	interfaces *Interfaces

	rawFile *ast.File
}

func (s *SourceFile) Name() string {
	return s.name
}

func (s *SourceFile) FullPath() string {
	return s.name
}

func (s *SourceFile) Markers() MarkerValues {
	return s.fileMarkers
}

func (s *SourceFile) Package() *Package {
	return s.pkg
}

func (s *SourceFile) Imports() *Imports {
	return s.imports
}

func (s *SourceFile) ImportMarkers() []ImportMarker {
	return s.importMarkers
}

func (s *SourceFile) Functions() *Functions {
	return s.functions
}

func (s *SourceFile) Structs() *Structs {
	return s.structs
}

func (s *SourceFile) Interfaces() *Interfaces {
	return s.interfaces
}

func (s *SourceFile) RawFile() *ast.File {
	return s.rawFile
}

type Position struct {
	Line   int
	Column int
}

type T interface {
	Underlying() T
	String() string
}

type BasicKind int

const (
	Invalid BasicKind = iota

	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String
	UnsafePointer

	UntypedBool
	UntypedInt
	UntypedRune
	UntypedFloat
	UntypedComplex
	UntypedString
	UntypedNil

	Byte = Uint8
	Rune = Int32
)

// BasicInfo is a set of flags describing properties of a basic type.
type BasicInfo int

// Properties of basic types.
const (
	IsBoolean BasicInfo = 1 << iota
	IsInteger
	IsUnsigned
	IsFloat
	IsComplex
	IsString
	IsUntyped

	IsOrdered   = IsInteger | IsFloat | IsString
	IsNumeric   = IsInteger | IsFloat | IsComplex
	IsConstType = IsBoolean | IsNumeric | IsString
)

var basicTypesMap = map[string]*Basic{
	"int":        basicTypes[Int],
	"int8":       basicTypes[Int8],
	"int16":      basicTypes[Int16],
	"int32":      basicTypes[Int32],
	"int64":      basicTypes[Int64],
	"uint":       basicTypes[Uint],
	"uint8":      basicTypes[Uint8],
	"uint16":     basicTypes[Uint16],
	"uint32":     basicTypes[Uint32],
	"uint64":     basicTypes[Uint64],
	"uintptr":    basicTypes[Uintptr],
	"float32":    basicTypes[Float32],
	"float64":    basicTypes[Float64],
	"complex64":  basicTypes[Complex64],
	"complex128": basicTypes[Complex128],
	"string":     basicTypes[String],
	"byte":       basicTypes[Byte],
	"rune":       basicTypes[Rune],
}

var basicTypes = []*Basic{
	Invalid: {Invalid, 0, "invalid type"},

	Bool:          {Bool, IsBoolean, "bool"},
	Int:           {Int, IsInteger, "int"},
	Int8:          {Int8, IsInteger, "int8"},
	Int16:         {Int16, IsInteger, "int16"},
	Int32:         {Int32, IsInteger, "int32"},
	Int64:         {Int64, IsInteger, "int64"},
	Uint:          {Uint, IsInteger | IsUnsigned, "uint"},
	Uint8:         {Uint8, IsInteger | IsUnsigned, "uint8"},
	Uint16:        {Uint16, IsInteger | IsUnsigned, "uint16"},
	Uint32:        {Uint32, IsInteger | IsUnsigned, "uint32"},
	Uint64:        {Uint64, IsInteger | IsUnsigned, "uint64"},
	Uintptr:       {Uintptr, IsInteger | IsUnsigned, "uintptr"},
	Float32:       {Float32, IsFloat, "float32"},
	Float64:       {Float64, IsFloat, "float64"},
	Complex64:     {Complex64, IsComplex, "complex64"},
	Complex128:    {Complex128, IsComplex, "complex128"},
	String:        {String, IsString, "string"},
	UnsafePointer: {UnsafePointer, 0, "Pointer"},

	UntypedBool:    {UntypedBool, IsBoolean | IsUntyped, "untyped bool"},
	UntypedInt:     {UntypedInt, IsInteger | IsUntyped, "untyped int"},
	UntypedRune:    {UntypedRune, IsInteger | IsUntyped, "untyped rune"},
	UntypedFloat:   {UntypedFloat, IsFloat | IsUntyped, "untyped float"},
	UntypedComplex: {UntypedComplex, IsComplex | IsUntyped, "untyped complex"},
	UntypedString:  {UntypedString, IsString | IsUntyped, "untyped string"},
	UntypedNil:     {UntypedNil, IsUntyped, "untyped nil"},

	{Byte, IsInteger | IsUnsigned, "byte"},
	{Rune, IsInteger, "rune"},
}

type Basic struct {
	kind BasicKind
	info BasicInfo
	name string
}

func (b *Basic) Kind() BasicKind {
	return b.kind
}

func (b *Basic) Name() string {
	return b.name
}

func (b *Basic) Underlying() T {
	return b
}

func (b *Basic) String() string {
	return b.name
}

type Variable struct {
	name string
	typ  T
}

func (v *Variable) Name() string {
	return v.name
}

func (v *Variable) Type() T {
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

type Array struct {
	len  int64
	elem T
}

func (a *Array) Len() int64 {
	return a.len
}

func (a *Array) Elem() T {
	return a.elem
}

func (a *Array) Underlying() T {
	return a
}

func (a *Array) String() string {
	return ""
}

type Slice struct {
	elem T
}

func (s *Slice) Elem() T {
	return s.elem
}

func (s *Slice) Underlying() T {
	return s
}

func (s *Slice) String() string {
	return ""
}

type Map struct {
	key  T
	elem T
}

func (m *Map) Key() T {
	return m.key
}

func (m *Map) Elem() T {
	return m.elem
}

func (m *Map) Underlying() T {
	return nil
}

func (m *Map) String() string {
	return ""
}

type ChanDirection int

const (
	SEND ChanDirection = 1 << iota
	RECEIVE
)

type Chan struct {
	direction ChanDirection
	elem      T
}

func (c *Chan) Direction() ChanDirection {
	return c.direction
}

func (c *Chan) Elem() T {
	return c.elem
}

func (c *Chan) Underlying() T {
	return c
}

func (c *Chan) String() string {
	return ""
}

type Pointer struct {
	base T
}

func (p *Pointer) Elem() T {
	return p.base
}

func (p *Pointer) Underlying() T {
	return nil
}

func (p *Pointer) String() string {
	return ""
}

type Function struct {
	name       string
	isExported bool
	position   Position
	receiver   *Variable
	params     *Tuple
	results    *Tuple
	variadic   bool

	file *SourceFile

	rawFile     *ast.File
	rawField    *ast.Field
	rawFuncDecl *ast.FuncDecl
}

func (f *Function) File() *SourceFile {
	return f.file
}

func (f *Function) Position() Position {
	return f.position
}

func (f *Function) Underlying() T {
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

func (f *Function) RawFile() *ast.File {
	return f.rawFile
}

func (f *Function) RawField() *ast.Field {
	return f.rawField
}

func (f *Function) RawFuncDecl() *ast.FuncDecl {
	return f.rawFuncDecl
}

type Field struct {
	name       string
	isExported bool
	tags       string
	typ        T
	position   Position
	markers    MarkerValues
	file       *SourceFile
	rawFile    *ast.File
	rawField   *ast.Field
}

func (f *Field) Name() string {
	return f.name
}

func (f *Field) Type() T {
	return f.typ
}

func (f *Field) IsExported() bool {
	return f.isExported
}

func (f *Field) Tags() string {
	return f.tags
}

type DefinedType struct {
	underlying T
}

func (d *DefinedType) Underlying() T {
	return d.underlying
}

func (d *DefinedType) String() string {
	return ""
}

type Struct struct {
	name       string
	isExported bool
	position   Position
	markers    MarkerValues
	fields     []*Field
	methods    []*Function
	file       *SourceFile

	rawFile     *ast.File
	rawGenDecl  *ast.GenDecl
	rawTypeSpec *ast.TypeSpec
}

func (s *Struct) File() *SourceFile {
	return s.file
}

func (s *Struct) Position() Position {
	return s.position
}

func (s *Struct) Underlying() T {
	return s
}

func (s *Struct) String() string {
	return ""
}

func (s *Struct) Name() string {
	return s.name
}

func (s *Struct) IsExported() bool {
	return s.isExported
}

func (s *Struct) Markers() MarkerValues {
	return s.markers
}

func (s *Struct) RawFile() *ast.File {
	return s.rawFile
}

func (s *Struct) RawGenDecl() *ast.GenDecl {
	return s.rawGenDecl
}

func (s *Struct) RawTypeSpec() *ast.TypeSpec {
	return s.rawTypeSpec
}

func (s *Struct) Implements(i *Interface) bool {
	if i == nil {
		return false
	}

	return false
}

type Interface struct {
	name       string
	isError    bool
	isExported bool
	position   Position
	markers    MarkerValues
	embeddeds  []Type
	allMethods []*Function
	methods    []*Function
	file       *SourceFile

	rawFile     *ast.File
	rawGenDecl  *ast.GenDecl
	rawTypeSpec *ast.TypeSpec
}

func (i *Interface) IsError() bool {
	return i.isError
}

func (i *Interface) IsEmptyInterface() bool {
	return false
}

func (i *Interface) File() *SourceFile {
	return i.file
}

func (i *Interface) Position() Position {
	return i.position
}

func (i *Interface) Underlying() T {
	return i
}

func (i *Interface) String() string {
	return ""
}

func (i *Interface) Name() string {
	return i.name
}

func (i *Interface) IsExported() bool {
	return i.isExported
}

func (i *Interface) Markers() MarkerValues {
	return i.markers
}

func (i *Interface) NumExplicitMethods() int {
	return len(i.methods)
}

func (i *Interface) ExplicitMethod(index int) *Function {
	if index >= 0 && index < len(i.methods) {
		return i.methods[index]
	}

	return nil
}

func (i *Interface) NumMethods() int {
	return len(i.allMethods)
}

func (i *Interface) Method(index int) *Function {
	if index >= 0 && index < len(i.allMethods) {
		return i.allMethods[index]
	}

	return nil
}

func (i *Interface) RawFile() *ast.File {
	return i.rawFile
}

func (i *Interface) RawGenDecl() *ast.GenDecl {
	return i.rawGenDecl
}

func (i *Interface) RawTypeSpec() *ast.TypeSpec {
	return i.rawTypeSpec
}

type SourceFileCallback func(file *SourceFile, err error)

func eachFile(collector *Collector, pkgs []*Package, callback SourceFileCallback) {
	if collector == nil {
		callback(nil, errors.New("collector cannot be nil"))
	}

	if pkgs == nil {
		callback(nil, errors.New("pkgs(packages) cannot be nil"))
	}

	var fileMap = make(map[*ast.File]*SourceFile)
	var errs []error

	pkgList := make([]*packages.Package, 0)

	for _, pkg := range pkgs {
		pkgList = append(pkgList, pkg.Package)
	}

	//pkgMap := make(map[string]*packages.Package)

	//detectAllImports(pkgList, pkgMap)

	for _, pkg := range pkgs {

		if len(pkg.Imports) != 0 {
			//emptyMarker := make(map[ast.Node]MarkerValues)
			//for _, imp := range pkg.Imports {
			//	loadedPackage, _ := LoadPackages(imp.ID)
			// 	if p, err := loadedPackage.Lookup("strings"); err == nil {
			//		EachPackage(p, emptyMarker)
			//	}
			//}
		}

		markers, err := collector.Collect(pkg)

		if err != nil {
			errs = append(errs, err.(ErrorList)...)
			continue
		}

		fileNodeMap := EachPackage(pkg, markers)

		for fileNode, file := range fileNodeMap {
			fileMap[fileNode] = file
		}
	}
}

func detectAllImports(pkgList []*packages.Package, m map[string]*packages.Package) {

	for _, pkg := range pkgList {

		if len(pkg.Imports) == 0 {
			continue
		}

		for pkgName, imp := range pkg.Imports {
			m[pkgName] = imp
			detectImports(imp, m)
		}

	}

}

func detectImports(p *packages.Package, m map[string]*packages.Package) {
	if len(p.Imports) == 0 {
		return
	}

	for pkgName, imp := range p.Imports {
		m[pkgName] = imp
		detectImports(imp, m)
	}
}

type FileElementVisitor struct {
	pkg         *Package
	markers     map[ast.Node]MarkerValues
	fileNodeMap map[*ast.File]*SourceFile

	files   []*SourceFile
	structs map[string]*Struct

	genDecl  *ast.GenDecl
	funcDecl *ast.FuncDecl
	file     *ast.File
	typeSpec *ast.TypeSpec
}

func (visitor *FileElementVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return visitor
	}

	switch typedNode := node.(type) {
	case *ast.File:
		visitor.file = typedNode
		visitor.createSourceFile()
		return visitor
	case *ast.GenDecl:
		visitor.genDecl = typedNode
		if typedNode.Tok == token.CONST {
			//visitor.constCallback(visitor.file, visitor.genDecl)
		}
		return visitor
	case *ast.FuncDecl:
		visitor.funcDecl = typedNode
		visitor.collectFunction()
		return nil
	case *ast.TypeSpec:
		visitor.typeSpec = typedNode
		visitor.getTypeFromTypeSpec()
		return nil
	default:
		return nil
	}
}

func (visitor *FileElementVisitor) createSourceFile() {
	_, exists := visitor.fileNodeMap[visitor.file]

	if exists {
		return
	}

	position := visitor.pkg.Fset.Position(visitor.file.Pos())
	fileFullPath := position.Filename

	file := &SourceFile{
		name:          filepath.Base(fileFullPath),
		fullPath:      fileFullPath,
		allMarkers:    visitor.markers[visitor.file],
		pkg:           visitor.pkg,
		imports:       visitor.getFileImports(),
		fileMarkers:   make(MarkerValues, 0),
		importMarkers: make([]ImportMarker, 0),
		functions:     &Functions{},
		structs:       &Structs{},
		interfaces:    &Interfaces{},
		rawFile:       visitor.file,
	}

	for markerName, markers := range file.allMarkers {
		if ImportMarkerName == markerName {
			for _, importMarker := range markers {
				file.importMarkers = append(file.importMarkers, importMarker.(ImportMarker))
			}
		} else {
			file.fileMarkers[markerName] = append(file.fileMarkers[markerName], markers...)
		}
	}

	visitor.fileNodeMap[visitor.file] = file
}

func (visitor *FileElementVisitor) collectFunction() {
	file, exists := visitor.fileNodeMap[visitor.file]

	if !exists {
		return
	}

	function := &Function{
		name:        visitor.funcDecl.Name.Name,
		isExported:  ast.IsExported(visitor.funcDecl.Name.Name),
		file:        file,
		position:    visitor.getPosition(visitor.funcDecl.Pos()),
		params:      &Tuple{},
		results:     &Tuple{},
		rawFile:     visitor.file,
		rawField:    nil,
		rawFuncDecl: visitor.funcDecl,
	}

	funcType := visitor.funcDecl.Type

	if funcType.Params != nil {
		function.params.variables = append(function.params.variables, visitor.getVariables(funcType.Params.List).variables...)
	}

	if funcType.Results != nil {
		function.results.variables = append(function.results.variables, visitor.getVariables(funcType.Results.List).variables...)
	}

	if visitor.funcDecl.Recv == nil {
		file.functions.functions = append(file.functions.functions, function)
	} else {
		receiverVariable := &Variable{
			name: visitor.funcDecl.Recv.List[0].Names[0].Name,
		}

		var receiverTypeSpec *ast.TypeSpec
		receiver := visitor.funcDecl.Recv.List[0].Type
		receiverTypeName := ""
		isPointerReceiver := false

		switch typedReceiver := receiver.(type) {
		case *ast.Ident:
			receiverTypeSpec = typedReceiver.Obj.Decl.(*ast.TypeSpec)
			receiverTypeName = receiverTypeSpec.Name.Name

		case *ast.StarExpr:
			receiverTypeSpec = typedReceiver.X.(*ast.Ident).Obj.Decl.(*ast.TypeSpec)
			receiverTypeName = receiverTypeSpec.Name.Name
			isPointerReceiver = true
		}

		structType, ok := visitor.structs[receiverTypeName]

		if !ok {
			structType = visitor.getStruct(receiverTypeSpec)
			visitor.structs[receiverTypeName] = structType
		}

		if isPointerReceiver {
			receiverVariable.typ = &Pointer{
				base: structType,
			}
		} else {
			receiverVariable.typ = structType
		}

		function.receiver = receiverVariable
		structType.methods = append(structType.methods, function)
	}
}

func (visitor *FileElementVisitor) getTypeFromTypeSpec() {
	typeName := visitor.typeSpec.Name.Name

	switch visitor.typeSpec.Type.(type) {
	case *ast.InterfaceType:
		interfaceType := visitor.getInterface(visitor.typeSpec)

		if interfaceType != nil {

		}
	case *ast.StructType:
		structType, ok := visitor.structs[typeName]

		if ok {
			structType.rawGenDecl = visitor.genDecl
		} else {
			structType = visitor.getStruct(visitor.typeSpec)
			visitor.structs[typeName] = structType
		}
	}
}

func (visitor *FileElementVisitor) getInterface(specType *ast.TypeSpec) *Interface {
	interfaceType := &Interface{
		name:        specType.Name.Name,
		isExported:  ast.IsExported(specType.Name.Name),
		methods:     visitor.getInterfaceMethods(specType.Type.(*ast.InterfaceType).Methods.List),
		position:    visitor.getPosition(specType.Pos()),
		markers:     visitor.markers[specType],
		rawFile:     visitor.file,
		rawGenDecl:  visitor.genDecl,
		rawTypeSpec: specType,
	}

	return interfaceType
}

func (visitor *FileElementVisitor) getStruct(specType *ast.TypeSpec) *Struct {

	structType := &Struct{
		name:        specType.Name.Name,
		isExported:  ast.IsExported(specType.Name.Name),
		position:    visitor.getPosition(specType.Pos()),
		markers:     visitor.markers[specType],
		fields:      make([]*Field, 0),
		methods:     make([]*Function, 0),
		rawFile:     visitor.file,
		rawGenDecl:  nil,
		rawTypeSpec: specType,
	}

	fieldList := specType.Type.(*ast.StructType).Fields.List
	visitor.getFieldsFromFieldList(fieldList)
	return structType
}

func (visitor *FileElementVisitor) getTypeFromExpression(expr ast.Expr) T {

	switch typed := expr.(type) {
	case *ast.Ident:
		typ, ok := basicTypesMap[typed.Name]

		if ok {
			return typ
		}

		if typed.Name == "error" {
			return &Interface{
				name:    typed.Name,
				isError: true,
			}
		}

	case *ast.SelectorExpr:
		importName := typed.X.(*ast.Ident).Name
		typeName := typed.Sel.Name

		if importName == "" {

		}

		if typeName == "" {

		}
	case *ast.StarExpr:
		return &Pointer{
			base: visitor.getTypeFromExpression(typed.X),
		}
	case *ast.ArrayType:

		if typed.Len == nil {
			return &Slice{
				elem: visitor.getTypeFromExpression(typed.Elt),
			}
		} else {
			length, _ := strconv.ParseInt(typed.Len.(*ast.BasicLit).Value, 10, 64)
			return &Array{
				elem: visitor.getTypeFromExpression(typed.Elt),
				len:  length,
			}
		}
	case *ast.ChanType:
		chanType := &Chan{
			elem: visitor.getTypeFromExpression(typed.Value),
		}

		if typed.Dir&ast.SEND == ast.SEND {
			chanType.direction |= SEND
		}

		if typed.Dir&ast.RECV == ast.RECV {
			chanType.direction |= RECEIVE
		}

		return chanType
	case *ast.FuncType:
		return &Function{}
	case *ast.MapType:
		return &Map{
			key:  visitor.getTypeFromExpression(typed.Key),
			elem: visitor.getTypeFromExpression(typed.Value),
		}
	case *ast.InterfaceType:
		interfaceType := &Interface{
			position: visitor.getPosition(typed.Pos()),
			rawFile:  visitor.file,
		}
		visitor.getInterfaceMethods(typed.Methods.List)
		return interfaceType
	case *ast.StructType:
		structType := &Struct{
			position: visitor.getPosition(typed.Pos()),
			rawFile:  visitor.file,
		}

		return structType
	}

	return nil
}

func (visitor *FileElementVisitor) getVariables(fieldList []*ast.Field) *Tuple {
	tuple := &Tuple{
		variables: make([]*Variable, 0),
	}

	for _, field := range fieldList {

		typ := visitor.getTypeFromExpression(field.Type)

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

func (visitor *FileElementVisitor) getInterfaceMethods(fieldList []*ast.Field) []*Function {
	methods := make([]*Function, 0)

	for _, rawMethod := range fieldList {
		funcType, ok := rawMethod.Type.(*ast.FuncType)

		if ok {
			method := &Function{
				params:  &Tuple{},
				results: &Tuple{},
			}

			if rawMethod.Names != nil {
				method.name = rawMethod.Names[0].Name
			}

			if funcType.Params != nil {
				method.params.variables = append(method.params.variables, visitor.getVariables(funcType.Params.List).variables...)
			}

			if funcType.Results != nil {
				method.results.variables = append(method.results.variables, visitor.getVariables(funcType.Results.List).variables...)
			}

			methods = append(methods, method)
		}
	}

	return methods
}

func (visitor *FileElementVisitor) getDefinedTypeFromExpression(expr ast.Expr) *DefinedType {
	return nil
}

func (visitor *FileElementVisitor) getFieldsFromFieldList(fieldList []*ast.Field) []*Field {
	file, _ := visitor.fileNodeMap[visitor.file]
	fields := make([]*Field, 0)

	for _, rawField := range fieldList {
		tags := ""

		if rawField.Tag != nil {
			tags = rawField.Tag.Value
		}

		if rawField.Names == nil {
			embeddedType := visitor.getDefinedTypeFromExpression(rawField.Type)

			field := &Field{
				name:       "",
				isExported: false,
				position:   Position{},
				markers:    visitor.markers[rawField],
				file:       file,
				rawFile:    visitor.file,
				rawField:   rawField,
				tags:       tags,
				typ:        embeddedType,
			}

			fields = append(fields, field)
			continue
		}

		for _, fieldName := range rawField.Names {
			typ := visitor.getTypeFromExpression(rawField.Type)

			field := &Field{
				name:       fieldName.Name,
				isExported: ast.IsExported(fieldName.Name),
				position:   visitor.getPosition(fieldName.Pos()),
				markers:    visitor.markers[rawField],
				file:       file,
				rawFile:    visitor.file,
				rawField:   rawField,
				tags:       tags,
				typ:        typ,
			}

			fields = append(fields, field)
		}

	}

	return fields
}

func (visitor *FileElementVisitor) getPosition(tokenPosition token.Pos) Position {
	position := visitor.pkg.Fset.Position(tokenPosition)
	return Position{
		Line:   position.Line,
		Column: position.Column,
	}
}

func (visitor *FileElementVisitor) getFileImports() *Imports {
	imports := &Imports{}

	for _, importPackage := range visitor.file.Imports {
		importPosition := visitor.getPosition(importPackage.Pos())
		importName := ""

		if importPackage.Name != nil {
			importName = importPackage.Name.Name
		}

		imports.imports = append(imports.imports, &Import{
			name: importName,
			path: importPackage.Path.Value[1 : len(importPackage.Path.Value)-1],
			position: Position{
				importPosition.Line,
				importPosition.Column,
			},
			rawImportSpec: importPackage,
		})
	}

	return imports
}

func EachPackage(pkg *Package, markers map[ast.Node]MarkerValues) map[*ast.File]*SourceFile {
	visitor := visitPackage(pkg, markers)

	if visitor != nil {

	}

	return nil
}

func visitPackage(pkg *Package, markers map[ast.Node]MarkerValues) *FileElementVisitor {
	visitor := &FileElementVisitor{
		pkg:         pkg,
		markers:     markers,
		fileNodeMap: make(map[*ast.File]*SourceFile),
		structs:     make(map[string]*Struct),
	}

	for _, file := range pkg.Syntax {
		ast.Walk(visitor, file)
	}

	return visitor
}
