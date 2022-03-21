package marker

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"
)

type T interface {
	Underlying() T
	String() string
}

type ImportedType struct {
	pkg *Package
	typ T
}

func (i *ImportedType) Package() *Package {
	return i.pkg
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
	if a.len != -1 {
		return "[" + fmt.Sprintf("%d", a.len) + "]" + a.elem.String()
	}

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
	return "[]" + s.elem.String()
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
	return m
}

func (m *Map) String() string {
	return "map[" + m.key.String() + "]" + m.elem.String()
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
	return p
}

func (p *Pointer) String() string {
	return "*" + p.base.String()
}

type Position struct {
	Line   int
	Column int
}

type Import struct {
	name       string
	path       string
	sideEffect bool
	position   Position
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

type Constant struct {
	name       string
	isExported bool
	value      interface{}
	typ        T
	expression ast.Expr
	initType   ast.Expr

	iota                int
	expressionEvaluated bool

	pkg     *Package
	visitor *PackageVisitor
}

func (c *Constant) Name() string {
	return c.name
}

func (c *Constant) Value() interface{} {
	c.evaluateExpression()
	return c.value
}

func (c *Constant) evaluateExpression() {
	if c.expressionEvaluated {
		return
	}

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	params := make(map[string]interface{}, 0)
	params["iota"] = c.iota
	c.value, c.typ = c.evalConstantExpression(c.expression, params)

	if c.initType != nil {
		switch typed := c.initType.(type) {
		case *ast.Ident:
			c.typ, _ = c.visitor.collector.findTypeByPkgIdAndName(c.pkg.ID, typed.Name)
		case *ast.SelectorExpr:
			c.typ = c.visitor.findTypeByImportAndTypeName(typed.X.(*ast.Ident).Name, typed.Sel.Name)
		}
	}

	c.expressionEvaluated = true
}

func (c *Constant) Type() T {
	return c.typ
}

func (c *Constant) IsExported() string {
	return c.name
}

func (c *Constant) Underlying() T {
	return c
}

func (c *Constant) String() string {
	return ""
}

func (c *Constant) evalConstantExpression(exp ast.Expr, variableMap map[string]interface{}) (interface{}, T) {
	switch exp := exp.(type) {
	case *ast.Ident:
		if value, ok := variableMap[exp.Name]; ok {
			return value, basicTypes[UntypedInt]
		}

		candidateConstant, ok := c.visitor.collector.findTypeByPkgIdAndName(c.pkg.ID, exp.Name)
		if ok {
			constant := candidateConstant.(*Constant)
			return constant.Value(), constant.Type()
		}

		return nil, nil
	case *ast.SelectorExpr:
		importedType := c.visitor.findTypeByImportAndTypeName(exp.X.(*ast.Ident).Name, exp.Sel.Name)
		if importedType != nil {
			constant := importedType.Underlying().(*Constant)
			return constant.Value(), constant.Type()
		}
	case *ast.BinaryExpr:
		return c.evalBinaryExpr(exp, variableMap)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT:
			i, _ := strconv.Atoi(exp.Value)
			return i, basicTypes[UntypedInt]
		case token.FLOAT:
			f, _ := strconv.ParseFloat(exp.Value, 64)
			return f, basicTypes[UntypedFloat]
		case token.STRING:
			return exp.Value[1 : len(exp.Value)-1], basicTypes[String]
		}
	case *ast.UnaryExpr:
		result, typ := c.evalConstantExpression(exp.X, variableMap)

		switch result.(type) {
		case int:
			return -1 * result.(int), typ
		case float64:
			return -1.0 * result.(float64), typ
		}
	case *ast.ParenExpr:
		return c.evalConstantExpression(exp.X, variableMap)
	}

	return nil, nil
}

func (c *Constant) evalBinaryExpr(exp *ast.BinaryExpr, variableMap map[string]interface{}) (interface{}, T) {
	var expressionType T
	left, typLeft := c.evalConstantExpression(exp.X, variableMap)
	right, typRight := c.evalConstantExpression(exp.Y, variableMap)

	_, isTypeLeftBasic := typLeft.(*Basic)
	_, isTypeRightBasic := typRight.(*Basic)

	if isTypeLeftBasic && isTypeRightBasic {
		expressionType = typLeft
	} else if !isTypeLeftBasic {
		expressionType = typLeft
	} else if !isTypeRightBasic {
		expressionType = typRight
	}

	switch left.(type) {
	case int:
		switch exp.Op {
		case token.ADD:
			return left.(int) + right.(int), expressionType
		case token.SUB:
			return left.(int) - right.(int), expressionType
		case token.MUL:
			return left.(int) * right.(int), expressionType
		case token.QUO:
			return left.(int) / right.(int), expressionType
		case token.REM:
			return left.(int) % right.(int), expressionType
		case token.AND:
			return left.(int) & right.(int), expressionType
		case token.OR:
			return left.(int) | right.(int), expressionType
		case token.XOR:
			return left.(int) ^ right.(int), expressionType
		case token.SHL:
			return left.(int) << right.(int), expressionType
		case token.SHR:
			return left.(int) >> right.(int), expressionType
		case token.AND_NOT:
			return left.(int) &^ right.(int), expressionType
		case token.EQL:
			return left.(int) == right.(int), basicTypes[Bool]
		case token.NEQ:
			return left.(int) != right.(int), basicTypes[Bool]
		case token.LSS:
			return left.(int) < right.(int), basicTypes[Bool]
		case token.GTR:
			return left.(int) > right.(int), basicTypes[Bool]
		case token.LEQ:
			return left.(int) <= right.(int), basicTypes[Bool]
		case token.GEQ:
			return left.(int) >= right.(int), basicTypes[Bool]
		}
	case float64:
		switch exp.Op {
		case token.ADD:
			return left.(float64) + right.(float64), expressionType
		case token.SUB:
			return left.(float64) - right.(float64), expressionType
		case token.MUL:
			return left.(float64) * right.(float64), expressionType
		case token.QUO:
			return left.(float64) / right.(float64), expressionType
		}
	case string:
		switch exp.Op {
		case token.ADD:
			return left.(string) + right.(string), expressionType
		}
	}

	return nil, nil
}

type Constants struct {
	elements []*Constant
}

func (c *Constants) Len() int {
	return len(c.elements)
}

func (c *Constants) At(index int) *Constant {
	if index >= 0 && index < len(c.elements) {
		return c.elements[index]
	}

	return nil
}

func (c *Constants) FindByName(name string) (*Constant, bool) {
	for _, constant := range c.elements {
		if constant.name == name {
			return constant, true
		}
	}

	return nil, false
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

type Function struct {
	name       string
	isExported bool
	position   Position
	receiver   *Variable
	params     *Tuple
	results    *Tuple
	variadic   bool

	file *SourceFile
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

type Field struct {
	name       string
	isExported bool
	tags       string
	typ        T
	position   Position
	markers    MarkerValues
	file       *SourceFile
	isEmbedded bool
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

func (f *Field) IsEmbedded() bool {
	return f.isEmbedded
}

func (f *Field) Tags() string {
	return f.tags
}

type Struct struct {
	name        string
	isExported  bool
	isAnonymous bool
	position    Position
	markers     MarkerValues
	fields      []*Field
	allFields   []*Field
	methods     []*Function
	file        *SourceFile

	isProcessed bool

	specType  *ast.TypeSpec
	namedType *types.Named

	fieldsLoaded    bool
	allFieldsLoaded bool
	visitor         *PackageVisitor
}

func (s *Struct) loadFields() {
	if s.fieldsLoaded {
		return
	}

	s.fields = append(s.fields, s.visitor.getFieldsFromFieldList(s.specType.Type.(*ast.StructType).Fields.List)...)
	s.fieldsLoaded = true
}

func (s *Struct) loadAllFields() {
	if s.allFieldsLoaded {
		return
	}

	s.loadFields()

	for _, field := range s.fields {

		if !field.IsEmbedded() {
			s.allFields = append(s.allFields, field)
			continue
		}

		var baseType = field.Type()
		pointerType, ok := field.Type().(*Pointer)

		if ok {
			baseType = pointerType.Elem()
		}

		importedType, ok := baseType.(*ImportedType)

		if ok {
			baseType = importedType.Underlying()
		}

		structType, ok := baseType.(*Struct)

		if ok {
			s.allFields = append(s.allFields, structType.AllFields()...)
		}

	}

	s.allFieldsLoaded = true
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

func (s *Struct) IsAnonymous() bool {
	return s.isAnonymous
}

func (s *Struct) Markers() MarkerValues {
	return s.markers
}

func (s *Struct) NamedType() *types.Named {
	return s.namedType
}

func (s *Struct) NumEmbeddedFields() int {
	s.loadFields()

	numEmbeddedFields := 0

	for _, field := range s.fields {
		if field.IsEmbedded() {
			numEmbeddedFields++
		}
	}

	return numEmbeddedFields
}

func (s *Struct) EmbeddedFields() []*Field {
	s.loadFields()

	embeddedFields := make([]*Field, 0)

	for _, field := range s.fields {
		if field.IsEmbedded() {
			embeddedFields = append(embeddedFields, field)
		}
	}

	return embeddedFields
}

func (s *Struct) NumFields() int {
	s.loadFields()
	return len(s.fields)
}

func (s *Struct) Fields() []*Field {
	s.loadFields()
	return s.fields
}

func (s *Struct) NumAllFields() int {
	s.loadAllFields()
	return len(s.allFields)
}

func (s *Struct) AllFields() []*Field {
	s.loadAllFields()
	return s.allFields
}

func (s *Struct) Implements(i *Interface) bool {
	if i == nil || i.interfaceType == nil || s.namedType == nil {
		return false
	}

	if types.Implements(s.namedType, i.interfaceType) {
		return true
	}

	pointerType := types.NewPointer(s.namedType)

	if types.Implements(pointerType, i.interfaceType) {
		return true
	}

	return false
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

type Interface struct {
	name        string
	isExported  bool
	isAnonymous bool
	position    Position
	markers     MarkerValues
	embeddeds   []T
	allMethods  []*Function
	methods     []*Function
	file        *SourceFile

	isProcessed bool

	specType      *ast.TypeSpec
	interfaceType *types.Interface

	embeddedTypesLoaded bool
	methodsLoaded       bool
	allMethodsLoaded    bool
	visitor             *PackageVisitor
}

func (i *Interface) loadEmbeddedTypes() {
	if i.embeddedTypesLoaded {
		return
	}

	i.embeddeds = i.visitor.getInterfaceEmbeddedTypes(i.specType.Type.(*ast.InterfaceType).Methods.List)
	i.embeddedTypesLoaded = true
}

func (i *Interface) loadMethods() {
	if i.methodsLoaded {
		return
	}

	i.methods = i.visitor.getInterfaceMethods(i.specType.Type.(*ast.InterfaceType).Methods.List)
	i.allMethods = append(i.allMethods, i.methods...)
	i.methodsLoaded = true
}

func (i *Interface) loadAllMethods() {
	if i.allMethodsLoaded {
		return
	}

	i.loadMethods()
	i.loadEmbeddedTypes()

	for _, embeddedType := range i.embeddeds {
		interfaceType, ok := embeddedType.(*Interface)

		if ok {
			interfaceType.loadAllMethods()
			i.allMethods = append(i.allMethods, interfaceType.allMethods...)
		}
	}

	i.allMethodsLoaded = true
}

func (i *Interface) IsEmptyInterface() bool {
	return len(i.embeddeds) == 0 && len(i.methods) == 0
}

func (i *Interface) IsAnonymous() bool {
	return i.isAnonymous
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
	var builder strings.Builder
	return builder.String()
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
	i.loadMethods()
	return len(i.methods)
}

func (i *Interface) ExplicitMethods() []*Function {
	i.loadMethods()
	return i.methods
}

func (i *Interface) NumEmbeddedTypes() int {
	i.loadEmbeddedTypes()
	return len(i.embeddeds)
}

func (i *Interface) EmbeddedTypes() []T {
	i.loadEmbeddedTypes()
	return i.embeddeds
}

func (i *Interface) NumMethods() int {
	i.loadAllMethods()
	return len(i.allMethods)
}

func (i *Interface) Methods() []*Function {
	i.loadAllMethods()
	return i.allMethods
}

func (i *Interface) InterfaceType() *types.Interface {
	return i.interfaceType
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

type CustomType struct {
	name       string
	aliasType  T
	isExported bool
	position   Position
	markers    MarkerValues
	methods    []*Function
	file       *SourceFile

	isProcessed bool
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
	constants   *Constants
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
