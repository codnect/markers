package marker

import (
	"errors"
	"go/ast"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"
)

type Import struct {
	name          string
	path          string
	sideEffect    bool
	position      Position
	rawImportSpec *ast.ImportSpec
}

func (i *Import) Name() string {
	return i.name
}

func (i *Import) Path() string {
	return i.path
}

func (i *Import) SideEffect() bool {
	return i.sideEffect
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
		if importItem.name == name || strings.HasSuffix(importItem.path, "/"+name) {
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

type Structs struct {
	elements []*Struct
}

func (s *Structs) Len() int {
	return len(s.elements)
}

func (s *Structs) At(index int) *Struct {
	if index >= 0 && index < len(s.elements) {
		return s.elements[index]
	}

	return nil
}

func (s *Structs) FindByName(name string) (*Struct, bool) {
	for _, structType := range s.elements {
		if structType.name == name {
			return structType, true
		}
	}

	return nil, false
}

type Interfaces struct {
	elements []*Interface
}

func (i *Interfaces) Len() int {
	return len(i.elements)
}

func (i *Interfaces) At(index int) *Interface {
	if index >= 0 && index < len(i.elements) {
		return i.elements[index]
	}

	return nil
}

func (i *Interfaces) FindByName(name string) (*Interface, bool) {
	for _, interfaceType := range i.elements {
		if interfaceType.name == name {
			return interfaceType, true
		}
	}

	return nil, false
}

type CustomTypes struct {
	elements []*CustomType
}

func (c *CustomTypes) Len() int {
	return len(c.elements)
}

func (c *CustomTypes) At(index int) *CustomType {
	if index >= 0 && index < len(c.elements) {
		return c.elements[index]
	}

	return nil
}

func (c *CustomTypes) FindByName(name string) (*CustomType, bool) {
	for _, customType := range c.elements {
		if customType.name == name {
			return customType, true
		}
	}

	return nil, false
}

type SourceFiles struct {
	elements []*SourceFile
}

func (s *SourceFiles) FindByName(name string) (*SourceFile, bool) {
	for _, file := range s.elements {
		if file.name == name {
			return file, true
		}
	}

	return nil, false
}

func (s *SourceFiles) Len() int {
	return len(s.elements)
}

func (s *SourceFiles) At(index int) *SourceFile {
	if index >= 0 && index < len(s.elements) {
		return s.elements[index]
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

	functions   *Functions
	structs     *Structs
	interfaces  *Interfaces
	customTypes *CustomTypes

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
	"bool":       basicTypes[Bool],
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

type ImportedType struct {
	typ T
}

func (i *ImportedType) Underlying() T {
	return i.typ
}

func (i *ImportedType) String() string {
	return ""
}

func (i *ImportedType) Name() string {
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

type CustomType struct {
	name       string
	aliasType  T
	isExported bool
	position   Position
	markers    MarkerValues
	methods    []*Function
	file       *SourceFile

	rawFile     *ast.File
	rawGenDecl  *ast.GenDecl
	rawTypeSpec *ast.TypeSpec
}

func (c *CustomType) Name() string {
	return c.name
}

func (c *CustomType) AliasType() T {
	return c.aliasType
}

func (c *CustomType) Underlying() T {
	return c
}

func (c *CustomType) String() string {
	return ""
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

	//var fileMap = make(map[*ast.File]*SourceFile)
	var errs []error

	packageMarkers := make(map[string]map[ast.Node]MarkerValues)

	for _, pkg := range pkgs {
		markers, err := collector.Collect(pkg)

		if err != nil {
			errs = append(errs, err.(ErrorList)...)
			continue
		}

		packageMarkers[pkg.ID] = markers
	}

	VisitPackages(pkgs, packageMarkers)

	/*fileNodeMap := EachPackage(cache, pkg, markers)

	for fileNode, file := range fileNodeMap {
		fileMap[fileNode] = file
	}*/
}

type packageCollector struct {
	hasSeen map[string]bool
	files   map[string]*SourceFiles
}

func newPackageCollector() *packageCollector {
	return &packageCollector{
		hasSeen: make(map[string]bool),
		files:   make(map[string]*SourceFiles),
	}
}

func (collector *packageCollector) markAsSeen(pkgId string) {
	collector.hasSeen[pkgId] = true
}

func (collector *packageCollector) isVisited(pkgId string) bool {
	visited, ok := collector.hasSeen[pkgId]

	if !ok {
		return false
	}

	return visited
}

func (collector *packageCollector) addFile(pkgId string, file *SourceFile) {
	if _, ok := collector.files[pkgId]; !ok {
		collector.files[pkgId] = &SourceFiles{
			elements: make([]*SourceFile, 0),
		}
	}

	if _, ok := collector.files[pkgId].FindByName(file.name); ok {
		return
	}

	collector.files[pkgId].elements = append(collector.files[pkgId].elements, file)
}

func (collector *packageCollector) findTypeByPkgIdAndName(pkgId, typeName string) (T, bool) {
	if files, ok := collector.files[pkgId]; ok {

		for i := 0; i < files.Len(); i++ {
			file := files.At(i)

			if structType, ok := file.structs.FindByName(typeName); ok {
				return structType, true
			}

			if interfaceType, ok := file.interfaces.FindByName(typeName); ok {
				return interfaceType, true
			}

			if customType, ok := file.customTypes.FindByName(typeName); ok {
				return customType, true
			}
		}

	}

	return nil, false
}

type PackageVisitor struct {
	collector *packageCollector

	pkg               *Package
	packageMarkers    map[ast.Node]MarkerValues
	allPackageMarkers map[string]map[ast.Node]MarkerValues

	currentFile *SourceFile

	genDecl  *ast.GenDecl
	funcDecl *ast.FuncDecl
	file     *ast.File
	typeSpec *ast.TypeSpec
}

func (visitor *PackageVisitor) VisitPackage() {
	visitor.packageMarkers = visitor.allPackageMarkers[visitor.pkg.ID]
	visitor.collector.markAsSeen(visitor.pkg.ID)

	for _, file := range visitor.pkg.Syntax {
		ast.Walk(visitor, file)
	}

	if visitor != nil {

	}
}

func (visitor *PackageVisitor) Visit(node ast.Node) ast.Visitor {
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

func (visitor *PackageVisitor) getPosition(tokenPosition token.Pos) Position {
	position := visitor.pkg.Fset.Position(tokenPosition)
	return Position{
		Line:   position.Line,
		Column: position.Column,
	}
}

func (visitor *PackageVisitor) createSourceFile() *SourceFile {
	position := visitor.pkg.Fset.Position(visitor.file.Pos())
	fileFullPath := position.Filename

	file := &SourceFile{
		name:          filepath.Base(fileFullPath),
		fullPath:      fileFullPath,
		allMarkers:    visitor.packageMarkers[visitor.file],
		pkg:           visitor.pkg,
		imports:       visitor.getFileImports(),
		fileMarkers:   make(MarkerValues, 0),
		importMarkers: make([]ImportMarker, 0),
		functions:     &Functions{},
		structs:       &Structs{},
		interfaces:    &Interfaces{},
		customTypes:   &CustomTypes{},
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

	visitor.currentFile = file
	visitor.collector.addFile(visitor.pkg.ID, file)
	return file
}

func (visitor *PackageVisitor) getFileImports() *Imports {
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

func (visitor *PackageVisitor) collectFunction() {
	function := &Function{
		name:        visitor.funcDecl.Name.Name,
		isExported:  ast.IsExported(visitor.funcDecl.Name.Name),
		file:        visitor.currentFile,
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
		visitor.currentFile.functions.elements = append(visitor.currentFile.functions.elements, function)
	} else {
		receiverVariable := &Variable{
			name: visitor.funcDecl.Recv.List[0].Names[0].Name,
		}

		var receiverTypeSpec *ast.TypeSpec
		receiver := visitor.funcDecl.Recv.List[0].Type

		receiverTypeName := ""
		isPointerReceiver := false
		isStructMethod := false

		switch typedReceiver := receiver.(type) {
		case *ast.Ident:
			receiverTypeSpec = typedReceiver.Obj.Decl.(*ast.TypeSpec)
			receiverTypeName = receiverTypeSpec.Name.Name
			_, isStructMethod = receiverTypeSpec.Type.(*ast.StructType)
		case *ast.StarExpr:
			receiverTypeSpec = typedReceiver.X.(*ast.Ident).Obj.Decl.(*ast.TypeSpec)
			receiverTypeName = receiverTypeSpec.Name.Name
			isPointerReceiver = true
			_, isStructMethod = receiverTypeSpec.Type.(*ast.StructType)
		}

		candidateType, ok := visitor.collector.findTypeByPkgIdAndName(visitor.pkg.ID, receiverTypeName)

		if isStructMethod {
			var structType *Struct

			if !ok {
				structType = visitor.getStruct(receiverTypeSpec)
				visitor.currentFile.structs.elements = append(visitor.currentFile.structs.elements, structType)
				candidateType = structType
			} else {
				structType = candidateType.(*Struct)
			}

			structType.methods = append(structType.methods, function)
		} else {
			var customType *CustomType

			if !ok {
				customType = visitor.getCustomType(receiverTypeSpec)
				visitor.currentFile.customTypes.elements = append(visitor.currentFile.customTypes.elements, customType)
				candidateType = customType
			} else {
				customType = candidateType.(*CustomType)
			}

			customType.methods = append(customType.methods, function)
		}

		if isPointerReceiver {
			receiverVariable.typ = &Pointer{
				base: candidateType,
			}
		} else {
			receiverVariable.typ = candidateType
		}

		function.receiver = receiverVariable
	}
}

func (visitor *PackageVisitor) getTypeFromTypeSpec() T {
	typeName := visitor.typeSpec.Name.Name

	switch visitor.typeSpec.Type.(type) {
	case *ast.InterfaceType:
		interfaceType := visitor.getInterface(visitor.typeSpec)
		visitor.currentFile.interfaces.elements = append(visitor.currentFile.interfaces.elements, interfaceType)
		return interfaceType
	case *ast.StructType:
		structCandidate, ok := visitor.collector.findTypeByPkgIdAndName(visitor.pkg.ID, typeName)

		var structType *Struct
		if ok {
			structType = structCandidate.(*Struct)
			structType.rawGenDecl = visitor.genDecl
		} else {
			structType = visitor.getStruct(visitor.typeSpec)
			visitor.currentFile.structs.elements = append(visitor.currentFile.structs.elements, structType)
		}

		return structType
	}

	customTypeCandidate, ok := visitor.collector.findTypeByPkgIdAndName(visitor.pkg.ID, typeName)

	var customType *CustomType
	if ok {
		customType = customTypeCandidate.(*CustomType)
		customType.rawGenDecl = visitor.genDecl
	} else {
		customType := visitor.getCustomType(visitor.typeSpec)
		visitor.currentFile.customTypes.elements = append(visitor.currentFile.customTypes.elements, customType)
	}

	return customType
}

func (visitor *PackageVisitor) getInterface(specType *ast.TypeSpec) *Interface {
	interfaceType := &Interface{
		name:        specType.Name.Name,
		isExported:  ast.IsExported(specType.Name.Name),
		methods:     visitor.getInterfaceMethods(specType.Type.(*ast.InterfaceType).Methods.List),
		position:    visitor.getPosition(specType.Pos()),
		markers:     visitor.packageMarkers[specType],
		file:        visitor.currentFile,
		rawFile:     visitor.file,
		rawGenDecl:  visitor.genDecl,
		rawTypeSpec: specType,
	}

	return interfaceType
}

func (visitor *PackageVisitor) getStruct(specType *ast.TypeSpec) *Struct {

	structType := &Struct{
		name:        specType.Name.Name,
		isExported:  ast.IsExported(specType.Name.Name),
		position:    visitor.getPosition(specType.Pos()),
		markers:     visitor.packageMarkers[specType],
		file:        visitor.currentFile,
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

func (visitor *PackageVisitor) getCustomType(specType *ast.TypeSpec) *CustomType {

	customType := &CustomType{
		name:        specType.Name.Name,
		aliasType:   visitor.getTypeFromExpression(specType.Type),
		isExported:  ast.IsExported(specType.Name.Name),
		position:    visitor.getPosition(specType.Pos()),
		markers:     visitor.packageMarkers[specType],
		methods:     make([]*Function, 0),
		file:        visitor.currentFile,
		rawFile:     visitor.file,
		rawGenDecl:  nil,
		rawTypeSpec: specType,
	}

	return customType
}

func (visitor *PackageVisitor) getTypeFromExpression(expr ast.Expr) T {

	switch typed := expr.(type) {
	case *ast.Ident:
		var typ T
		var ok bool
		typ, ok = basicTypesMap[typed.Name]

		if ok {
			return typ
		}

		if typed.Name == "error" {
			errorType, exists := visitor.collector.findTypeByPkgIdAndName("builtin", "error")

			if !exists {
				visitor.loadPackageAndVisit("builtin")
			}

			errorType, _ = visitor.collector.findTypeByPkgIdAndName("builtin", "error")

			return errorType

		}

		typ, ok = visitor.collector.findTypeByPkgIdAndName(visitor.pkg.ID, typed.Name)

		if ok {
			return typ
		}

		visitor.typeSpec = typed.Obj.Decl.(*ast.TypeSpec)
		return visitor.getTypeFromTypeSpec()
		/*
			visitor.typeSpec = typed.Obj.Decl.(*ast.TypeSpec)

			x := visitor.getTypeFromTypeSpec()

			if x == nil {
				return nil
			}

			return x*/
	case *ast.SelectorExpr:
		importName := typed.X.(*ast.Ident).Name
		typeName := typed.Sel.Name
		return visitor.findTypeByImportAndTypeName(importName, typeName)
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

func (visitor *PackageVisitor) loadPackageAndVisit(pkgId string) {
	if visitor.collector.isVisited(pkgId) {
		return
	}

	loadResult, err := LoadPackages(pkgId)

	if err != nil {
		panic(err)
	}

	pkg, _ := loadResult.Lookup(pkgId)

	visitPackage(pkg, visitor.collector, visitor.allPackageMarkers)
}

func (visitor *PackageVisitor) findTypeByImportAndTypeName(importName, typeName string) *ImportedType {
	packageImport, _ := visitor.currentFile.imports.FindByName(importName)
	typ, exists := visitor.collector.findTypeByPkgIdAndName(packageImport.path, typeName)

	if exists {
		return &ImportedType{
			typ,
		}
	}

	visitor.loadPackageAndVisit(packageImport.path)

	typ, exists = visitor.collector.findTypeByPkgIdAndName(packageImport.path, typeName)

	if exists {
		return &ImportedType{
			typ,
		}
	}

	return nil
}

func (visitor *PackageVisitor) getVariables(fieldList []*ast.Field) *Tuple {
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

func (visitor *PackageVisitor) getInterfaceMethods(fieldList []*ast.Field) []*Function {
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

func (visitor *PackageVisitor) getDefinedTypeFromExpression(expr ast.Expr) *DefinedType {
	return nil
}

func (visitor *PackageVisitor) getFieldsFromFieldList(fieldList []*ast.Field) []*Field {
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
				markers:    visitor.packageMarkers[rawField],
				file:       visitor.currentFile,
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
				markers:    visitor.packageMarkers[rawField],
				file:       visitor.currentFile,
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

func VisitPackages(pkgList []*Package, allPackageMarkers map[string]map[ast.Node]MarkerValues) {
	pkgCollector := newPackageCollector()

	for _, pkg := range pkgList {
		if !pkgCollector.isVisited(pkg.ID) {
			visitPackage(pkg, pkgCollector, allPackageMarkers)
		}
	}

	if pkgCollector == nil {

	}
}

func visitPackage(pkg *Package, collector *packageCollector, allPackageMarkers map[string]map[ast.Node]MarkerValues) {
	pkgVisitor := &PackageVisitor{
		collector:         collector,
		pkg:               pkg,
		allPackageMarkers: allPackageMarkers,
	}

	pkgVisitor.VisitPackage()
}

func IsInterfaceType(t T) bool {
	_, ok := t.(*Interface)
	return ok
}

func IsStructType(t T) bool {
	_, ok := t.(*Struct)
	return ok
}

func IsErrorType(t T) bool {
	interfaceType, ok := t.(*Interface)
	if !ok {
		return false
	}

	if interfaceType.file == nil || interfaceType.file.pkg == nil {
		return false
	}

	return interfaceType.file.pkg.ID == "builtin"
}
