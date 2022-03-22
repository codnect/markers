package visitor

import (
	"github.com/procyon-projects/marker"
	"go/ast"
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

type Struct struct {
	name        string
	isExported  bool
	isAnonymous bool
	position    Position
	markers     marker.MarkerValues
	fields      []*Field
	allFields   []*Field
	methods     []*Function
	file        *File

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
	return s.name
}

func (s *Struct) IsExported() bool {
	return s.isExported
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
