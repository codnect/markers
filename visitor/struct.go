package visitor

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"sync"
)

type Field struct {
	name       string
	isExported bool
	tags       string
	typ        Type
	position   Position
	markers    markers.Values
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

func (f *Field) Markers() markers.Values {
	return f.markers
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
	markers     markers.Values
	fields      []*Field
	allFields   []*Field
	methods     []*Function
	allMethods  []*Function
	typeParams  *TypeParameters
	file        *File

	isProcessed bool

	specType  *ast.TypeSpec
	namedType *types.Named
	fieldList []*ast.Field

	pkg     *packages.Package
	visitor *packageVisitor

	typeParamsOnce sync.Once
	methodsOnce    sync.Once
	allMethodsOnce sync.Once
	fieldsOnce     sync.Once
	allFieldsOnce  sync.Once
}

func newStruct(specType *ast.TypeSpec, structType *ast.StructType, file *File, pkg *packages.Package, visitor *packageVisitor, markers markers.Values) *Struct {
	s := &Struct{
		markers:   markers,
		file:      file,
		fields:    make([]*Field, 0),
		allFields: make([]*Field, 0),
		methods:   make([]*Function, 0),
		typeParams: &TypeParameters{
			[]*TypeParameter{},
		},
		isProcessed: true,
		specType:    specType,
		pkg:         pkg,
		visitor:     visitor,
	}

	return s.initialize(specType, structType, file, pkg)
}

func (s *Struct) initialize(specType *ast.TypeSpec, structType *ast.StructType, file *File, pkg *packages.Package) *Struct {
	s.isProcessed = true
	s.specType = specType
	s.file = file
	s.pkg = pkg

	if specType != nil {
		s.name = specType.Name.Name
		s.isExported = ast.IsExported(specType.Name.Name)
		s.position = getPosition(pkg, specType.Pos())
		s.namedType = file.pkg.Types.Scope().Lookup(specType.Name.Name).Type().(*types.Named)
		s.fieldList = s.specType.Type.(*ast.StructType).Fields.List
		if _, exists := file.structs.FindByName(s.name); !exists {
			s.file.structs.elements = append(s.file.structs.elements, s)
		}
	} else if structType != nil {
		if structType.Pos() != token.NoPos {
			//i.position = getPosition(pkg, interfaceType.Pos())
		}
		s.fieldList = structType.Fields.List
		s.isAnonymous = true
	}

	return s
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
	if s.name == "" && len(s.fieldList) == 0 {
		return "struct{}"
	}

	var builder strings.Builder
	if s.file != nil && s.file.pkg.Name != "builtin" {
		builder.WriteString(fmt.Sprintf("%s.%s", s.file.Package().Name, s.name))
	} else if s.name != "" {
		builder.WriteString(s.name)
	}

	if s.TypeParameters().Len() != 0 {
		builder.WriteString("[")

		for index := 0; index < s.TypeParameters().Len(); index++ {
			typeParam := s.TypeParameters().At(index)
			builder.WriteString(typeParam.String())

			if index != s.TypeParameters().Len()-1 {
				builder.WriteString(",")
			}
		}

		builder.WriteString("]")
	}

	return builder.String()
}

func (s *Struct) Name() string {
	if s.name == "" && len(s.fieldList) == 0 {
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

func (s *Struct) Markers() markers.Values {
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

func (s *Struct) TypeParameters() *TypeParameters {
	s.loadTypeParams()
	return s.typeParams
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

func (s *Struct) getFieldsFromFieldList() []*Field {
	fields := make([]*Field, 0)

	markers := s.visitor.allPackageMarkers[s.pkg.ID]

	for _, rawField := range s.fieldList {
		tags := ""

		if rawField.Tag != nil {
			tags = rawField.Tag.Value
		}

		if rawField.Names == nil {
			embeddedType := getTypeFromExpression(rawField.Type, s.file, s.visitor, nil, nil)
			typ := embeddedType

			pointerType, isPointerType := typ.(*Pointer)
			if isPointerType {
				typ = pointerType.Elem()
			}

			genericType, isGenericType := typ.(*GenericType)

			if isGenericType {
				typ = genericType.RawType()
			}

			nameParts := strings.SplitN(typ.Name(), ".", 2)
			name := nameParts[0]

			if len(nameParts) == 2 {
				name = nameParts[1]
			}

			field := &Field{
				name:       name,
				isExported: ast.IsExported(name),
				// TODO set position
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
			typ := getTypeFromExpression(rawField.Type, s.file, s.visitor, nil, nil)

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
	s.fieldsOnce.Do(func() {
		s.loadTypeParams()
		s.fields = append(s.fields, s.getFieldsFromFieldList()...)
	})
}

func (s *Struct) loadAllFields() {
	s.allFieldsOnce.Do(func() {
		s.loadFields()

		for _, field := range s.fields {

			if !field.IsEmbedded() {
				s.allFields = append(s.allFields, field)
				continue
			}

			var baseType = field.Type()
			pointerType, isPointer := field.Type().(*Pointer)

			if isPointer {
				baseType = pointerType.Elem()
			}

			genericType, isGenericType := baseType.(*GenericType)

			if isGenericType {
				baseType = genericType.RawType()
			}

			structType, isStruct := baseType.(*Struct)

			if isStruct {
				s.allFields = append(s.allFields, structType.FieldsInHierarchy().ToSlice()...)
			}
		}
	})
}

func (s *Struct) loadMethods() {
	s.methodsOnce.Do(func() {
		s.loadTypeParams()
		s.allMethods = append(s.allMethods, s.methods...)
	})
}

func (s *Struct) loadAllMethods() {
	s.allMethodsOnce.Do(func() {
		s.loadMethods()
		s.loadFields()

		for _, field := range s.fields {

			if !field.IsEmbedded() {
				continue
			}

			var baseType = field.Type()
			pointerType, isPointer := field.Type().(*Pointer)

			if isPointer {
				baseType = pointerType.Elem()
			}

			genericType, isGenericType := baseType.(*GenericType)

			if isGenericType {
				baseType = genericType.RawType()
			}

			structType, isStructType := baseType.(*Struct)

			if isStructType {
				s.allMethods = append(s.allMethods, structType.MethodsInHierarchy().ToSlice()...)
			}

			interfaceType, isInterfaceType := baseType.(*Interface)

			if isInterfaceType {
				s.allMethods = append(s.allMethods, interfaceType.Methods().ToSlice()...)
			}
		}
	})

}

func (s *Struct) loadTypeParams() {
	s.typeParamsOnce.Do(func() {
		if s.specType == nil || s.specType.TypeParams == nil {
			return
		}

		for _, field := range s.specType.TypeParams.List {
			for _, fieldName := range field.Names {
				typeParameter := &TypeParameter{
					name: fieldName.Name,
					constraints: &TypeConstraints{
						[]*TypeConstraint{},
					},
				}
				s.typeParams.elements = append(s.typeParams.elements, typeParameter)
			}
		}

		for _, field := range s.specType.TypeParams.List {
			constraints := make([]*TypeConstraint, 0)
			typ := getTypeFromExpression(field.Type, s.file, s.visitor, nil, s.typeParams)

			if typeSets, isTypeSets := typ.(TypeSets); isTypeSets {
				for _, item := range typeSets {
					if constraint, isConstraint := item.(*TypeConstraint); isConstraint {
						constraints = append(constraints, constraint)
					} else {
						constraints = append(constraints, &TypeConstraint{typ: item})
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
				typeParam, exists := s.typeParams.FindByName(fieldName.Name)

				if exists {
					typeParam.constraints.elements = append(typeParam.constraints.elements, constraints...)
				}
			}
		}

	})
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
