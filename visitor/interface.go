package visitor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
)

type Constraint struct {
}

func (c *Constraint) Name() string {
	return ""
}

func (c *Constraint) Underlying() Type {
	return c
}

func (c *Constraint) String() string {
	return ""
}

type Interface struct {
	name        string
	isExported  bool
	isAnonymous bool
	position    Position
	markers     marker.MarkerValues
	embeddeds   []Type
	constrains  []*Constraint
	allMethods  []*Function
	methods     []*Function
	file        *File

	isProcessed bool

	specType      *ast.TypeSpec
	interfaceType *types.Interface
	fieldList     []*ast.Field

	pkg     *packages.Package
	visitor *packageVisitor

	constraintsLoaded   bool
	embeddedTypesLoaded bool
	methodsLoaded       bool
	allMethodsLoaded    bool
}

func newInterface(specType *ast.TypeSpec, interfaceType *ast.InterfaceType, file *File, pkg *packages.Package, visitor *packageVisitor, markers marker.MarkerValues) *Interface {
	i := &Interface{
		methods:     make([]*Function, 0),
		allMethods:  make([]*Function, 0),
		embeddeds:   make([]Type, 0),
		constrains:  make([]*Constraint, 0),
		markers:     markers,
		file:        file,
		isProcessed: true,
		specType:    specType,
		pkg:         pkg,
		visitor:     visitor,
	}

	return i.initialize(specType, interfaceType, pkg)
}

func (i *Interface) initialize(specType *ast.TypeSpec, interfaceType *ast.InterfaceType, pkg *packages.Package) *Interface {
	if specType != nil {
		i.name = specType.Name.Name
		i.isExported = ast.IsExported(specType.Name.Name)
		i.position = getPosition(pkg, specType.Pos())
		typ := pkg.Types.Scope().Lookup(specType.Name.Name).Type()
		underlyingType := typ.Underlying()

		switch underlyingType.(type) {
		case *types.Interface:
			i.interfaceType = underlyingType.(*types.Interface)
			i.fieldList = i.specType.Type.(*ast.InterfaceType).Methods.List
		default:
		}

		i.file.interfaces.elements = append(i.file.interfaces.elements, i)
	} else if interfaceType != nil {
		if interfaceType.Pos() != token.NoPos {
			//i.position = getPosition(pkg, interfaceType.Pos())
		}
		i.fieldList = interfaceType.Methods.List
		i.isAnonymous = true
	}
	return i
}

func (i *Interface) getInterfaceMethods() []*Function {
	methods := make([]*Function, 0)

	markers := i.visitor.allPackageMarkers[i.pkg.ID]

	for _, rawMethod := range i.fieldList {
		_, ok := rawMethod.Type.(*ast.FuncType)

		if ok {
			methods = append(methods, newFunction(nil, rawMethod, i.file, i.pkg, i.visitor, markers[rawMethod]))
		}
	}

	return methods
}

func (i *Interface) getInterfaceEmbeddedTypes() []Type {
	embeddedTypes := make([]Type, 0)

	for _, field := range i.fieldList {
		_, ok := field.Type.(*ast.FuncType)

		if !ok {
			embeddedTypes = append(embeddedTypes, getTypeFromExpression(field.Type, i.visitor))
		}
	}

	return embeddedTypes
}

func (i *Interface) loadEmbeddedTypes() {
	if i.embeddedTypesLoaded {
		return
	}

	i.embeddeds = i.getInterfaceEmbeddedTypes()
	i.embeddedTypesLoaded = true
}

func (i *Interface) loadMethods() {
	if i.methodsLoaded {
		return
	}

	i.methods = i.getInterfaceMethods()
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
	return len(i.fieldList) == 0
}

func (i *Interface) IsAnonymous() bool {
	return i.isAnonymous
}

func (i *Interface) File() *File {
	return i.file
}

func (i *Interface) Position() Position {
	return i.position
}

func (i *Interface) Underlying() Type {
	return i
}

func (i *Interface) String() string {
	var builder strings.Builder
	return builder.String()
}

func (i *Interface) IsConstraint() bool {
	return false
}

func (i *Interface) Constraints() []*Constraint {
	return i.constrains
}

func (i *Interface) Name() string {
	if len(i.fieldList) == 0 {
		return "interface{}"
	}

	return i.name
}

func (i *Interface) IsExported() bool {
	return i.isExported
}

func (i *Interface) Markers() marker.MarkerValues {
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

func (i *Interface) EmbeddedTypes() []Type {
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
