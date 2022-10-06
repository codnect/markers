package visitor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"go/token"
	"go/types"
)

type Field struct {
	name       string
	isExported bool
	tags       string
	typ        Type
	position   Position
	markers    marker.MarkerValues
	file       *File
	isEmbedded bool
}

func (f *Field) Name() string {
	return f.name
}

func (f *Field) Type() Type {
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

type Fields struct {
	elements []*Field
}

func (f *Fields) ToSlice() []*Field {
	return f.elements
}

func (f *Fields) Len() int {
	return len(f.elements)
}

func (f *Fields) At(index int) *Field {
	if index >= 0 && index < len(f.elements) {
		return f.elements[index]
	}

	return nil
}

func (f *Fields) FindByName(name string) (*Field, bool) {
	for _, field := range f.elements {
		if field.name == name {
			return field, true
		}
	}

	return nil, false
}

type Struct struct {
	name        string
	isExported  bool
	isAnonymous bool
	position    Position
	markers     marker.MarkerValues
	fields      []*Field
	allFields   []*Field
	methods     []*Function
	allMethods  []*Function
	file        *File

	isProcessed bool

	specType  *ast.TypeSpec
	namedType *types.Named
	fieldList []*ast.Field

	pkg     *packages.Package
	visitor *packageVisitor

	methodsLoaded    bool
	allMethodsLoaded bool

	fieldsLoaded    bool
	allFieldsLoaded bool
}

func newStruct(specType *ast.TypeSpec, structType *ast.StructType, file *File, pkg *packages.Package, visitor *packageVisitor, markers marker.MarkerValues) *Struct {
	s := &Struct{
		markers:     markers,
		file:        file,
		fields:      make([]*Field, 0),
		allFields:   make([]*Field, 0),
		methods:     make([]*Function, 0),
		isProcessed: true,
		specType:    specType,
		pkg:         pkg,
		visitor:     visitor,
	}

	return s.initialize(specType, structType, file, pkg)
}

func (s *Struct) initialize(specType *ast.TypeSpec, structType *ast.StructType, file *File, pkg *packages.Package) *Struct {
	if specType != nil {
		s.name = specType.Name.Name
		s.isExported = ast.IsExported(specType.Name.Name)
		s.position = getPosition(pkg, specType.Pos())
		s.namedType = file.pkg.Types.Scope().Lookup(specType.Name.Name).Type().(*types.Named)
		s.fieldList = s.specType.Type.(*ast.StructType).Fields.List
		s.file.structs.elements = append(s.file.structs.elements, s)
	} else if structType != nil {
		if structType.Pos() != token.NoPos {
			//i.position = getPosition(pkg, interfaceType.Pos())
		}
		s.fieldList = structType.Fields.List
		s.isAnonymous = true
	}

	return s
}

func (s *Struct) getFieldsFromFieldList() []*Field {
	fields := make([]*Field, 0)

	markers := s.visitor.allPackageMarkers[s.pkg.ID]

	for _, rawField := range s.fieldList {
		tags := ""

		if rawField.Tag != nil {
			tags = rawField.Tag.Value
		}

		if rawField.Names == nil {
			embeddedType := getTypeFromExpression(rawField.Type, s.file, s.visitor)

			field := &Field{
				name:       embeddedType.Name(),
				isExported: ast.IsExported(embeddedType.Name()),
				position:   Position{},
				markers:    markers[rawField],
				file:       s.file,
				tags:       tags,
				typ:        embeddedType,
				isEmbedded: true,
			}

			fields = append(fields, field)
			continue
		}

		for _, fieldName := range rawField.Names {
			typ := getTypeFromExpression(rawField.Type, s.file, s.visitor)

			field := &Field{
				name:       fieldName.Name,
				isExported: ast.IsExported(fieldName.Name),
				position:   getPosition(s.file.pkg, fieldName.Pos()),
				markers:    markers[rawField],
				file:       s.file,
				tags:       tags,
				typ:        typ,
				isEmbedded: false,
			}

			fields = append(fields, field)
		}

	}

	return fields
}

func (s *Struct) loadFields() {
	if s.fieldsLoaded {
		return
	}

	s.fields = append(s.fields, s.getFieldsFromFieldList()...)
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
			s.allFields = append(s.allFields, structType.FieldsInHierarchy().ToSlice()...)
		}

	}

	s.allFieldsLoaded = true
}

func (s *Struct) loadMethods() {
	if s.methodsLoaded {
		return
	}

	s.allMethods = append(s.allMethods, s.methods...)
	s.methodsLoaded = true
}

func (s *Struct) loadAllMethods() {
	if s.allMethodsLoaded {
		return
	}

	s.loadMethods()
	s.loadFields()

	for _, field := range s.fields {

		if !field.IsEmbedded() {
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
			s.allMethods = append(s.allMethods, structType.MethodsInHierarchy().ToSlice()...)
		}

		interfaceType, ok := baseType.(*Interface)

		if ok {
			s.allMethods = append(s.allMethods, interfaceType.Methods().ToSlice()...)
		}
	}

	s.allMethodsLoaded = true
}

func (s *Struct) File() *File {
	return s.file
}

func (s *Struct) Position() Position {
	return s.position
}

func (s *Struct) Underlying() Type {
	return s
}

func (s *Struct) String() string {
	return ""
}

func (s *Struct) Name() string {
	if len(s.fieldList) == 0 {
		return "struct{}"
	}

	return s.name
}

func (s *Struct) IsExported() bool {
	return s.isExported
}

func (s *Struct) IsEmpty() bool {
	return len(s.fieldList) == 0
}

func (s *Struct) IsAnonymous() bool {
	return s.isAnonymous
}

func (s *Struct) Markers() marker.MarkerValues {
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

func (s *Struct) EmbeddedFields() *Fields {
	s.loadFields()

	embeddedFields := make([]*Field, 0)

	for _, field := range s.fields {
		if field.IsEmbedded() {
			embeddedFields = append(embeddedFields, field)
		}
	}

	return &Fields{
		elements: embeddedFields,
	}
}

func (s *Struct) NumFields() int {
	s.loadFields()
	return len(s.fields)
}

func (s *Struct) Fields() *Fields {
	s.loadFields()
	return &Fields{
		elements: s.fields,
	}
}

func (s *Struct) NumFieldsInHierarchy() int {
	s.loadAllFields()
	return len(s.allFields)
}

func (s *Struct) FieldsInHierarchy() *Fields {
	s.loadAllFields()
	return &Fields{
		elements: s.allFields,
	}
}

func (s *Struct) NumMethods() int {
	s.loadMethods()
	return len(s.methods)
}

func (s *Struct) Methods() *Functions {
	s.loadMethods()
	return &Functions{
		elements: s.methods,
	}
}

func (s *Struct) NumMethodsInHierarchy() int {
	s.loadAllMethods()
	return len(s.allMethods)
}

func (s *Struct) MethodsInHierarchy() *Functions {
	s.loadAllMethods()
	return &Functions{
		elements: s.allMethods,
	}
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
