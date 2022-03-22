package visitor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"go/types"
	"strings"
)

type Interface struct {
	name        string
	isExported  bool
	isAnonymous bool
	position    Position
	markers     marker.MarkerValues
	embeddeds   []Type
	allMethods  []*Function
	methods     []*Function
	file        *File

	isProcessed bool

	specType      *ast.TypeSpec
	interfaceType *types.Interface

	embeddedTypesLoaded bool
	methodsLoaded       bool
	allMethodsLoaded    bool
	visitor             *packageVisitor
}

func newInterface(specType *ast.TypeSpec, file *File, pkg *packages.Package, markers marker.MarkerValues) *Interface {
	interfaceType := &Interface{
		name:          specType.Name.Name,
		isExported:    ast.IsExported(specType.Name.Name),
		methods:       make([]*Function, 0),
		embeddeds:     make([]Type, 0),
		position:      getPosition(pkg, specType.Pos()),
		markers:       markers,
		file:          file,
		isProcessed:   true,
		specType:      specType,
		interfaceType: pkg.Types.Scope().Lookup(specType.Name.Name).Type().Underlying().(*types.Interface),
	}

	return interfaceType
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

func (i *Interface) Name() string {
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
