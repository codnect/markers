package visitor

import (
	"fmt"
	"github.com/procyon-projects/markers"
	"github.com/procyon-projects/markers/packages"
	"go/ast"
	"go/token"
	"go/types"
	"strings"
	"sync"
)

type Interface struct {
	name               string
	isExported         bool
	isAnonymous        bool
	position           Position
	markers            markers.Values
	embeddedInterfaces []*Interface
	embeddedTypes      []Type
	typeParams         *TypeParameters
	allMethods         []*Function
	methods            []*Function
	file               *File

	isProcessed bool

	specType      *ast.TypeSpec
	interfaceType *types.Interface
	fieldList     []*ast.Field

	pkg     *packages.Package
	visitor *packageVisitor

	typeParamsOnce         sync.Once
	embeddedInterfacesOnce sync.Once
	embeddedTypesOnce      sync.Once
	methodsOnce            sync.Once
	allMethodsOnce         sync.Once
}

func newInterface(specType *ast.TypeSpec, interfaceType *ast.InterfaceType, file *File, pkg *packages.Package, visitor *packageVisitor, markers markers.Values) *Interface {
	i := &Interface{
		methods:            make([]*Function, 0),
		allMethods:         make([]*Function, 0),
		embeddedTypes:      make([]Type, 0),
		embeddedInterfaces: make([]*Interface, 0),
		typeParams:         &TypeParameters{},
		markers:            markers,
		file:               file,
		isProcessed:        true,
		specType:           specType,
		pkg:                pkg,
		visitor:            visitor,
	}

	return i.initialize(specType, interfaceType, file, pkg)
}

func (i *Interface) initialize(specType *ast.TypeSpec, interfaceType *ast.InterfaceType, file *File, pkg *packages.Package) *Interface {
	i.isProcessed = true
	i.specType = specType
	i.file = file
	i.pkg = pkg

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

		if _, exists := file.interfaces.FindByName(i.name); !exists {
			i.file.interfaces.elements = append(i.file.interfaces.elements, i)
		}
	} else if interfaceType != nil {
		if interfaceType.Pos() != token.NoPos {
			//i.position = getPosition(pkg, interfaceType.Pos())
		}
		i.fieldList = interfaceType.Methods.List
		i.isAnonymous = true
	}
	return i
}

func (i *Interface) Name() string {
	if i.name == "" && len(i.fieldList) == 0 {
		return "interface{}"
	}

	return i.name
}

func (i *Interface) IsEmpty() bool {
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

func (i *Interface) IsExported() bool {
	return i.isExported
}

func (i *Interface) Markers() markers.Values {
	return i.markers
}

func (i *Interface) NumExplicitMethods() int {
	i.loadMethods()
	return len(i.methods)
}

func (i *Interface) ExplicitMethods() *Functions {
	i.loadMethods()
	return &Functions{
		elements: i.methods,
	}
}

func (i *Interface) NumEmbeddedInterfaces() int {
	i.loadEmbeddedInterfaces()
	return len(i.embeddedInterfaces)
}

func (i *Interface) EmbeddedInterfaces() *Interfaces {
	i.loadEmbeddedInterfaces()
	return &Interfaces{
		elements: i.embeddedInterfaces,
	}
}

func (i *Interface) NumEmbeddedTypes() int {
	i.loadEmbeddedTypes()
	return len(i.embeddedTypes)
}

func (i *Interface) EmbeddedTypes() *Types {
	i.loadEmbeddedTypes()
	return &Types{
		i.embeddedTypes,
	}
}

func (i *Interface) NumMethods() int {
	i.loadAllMethods()
	return len(i.allMethods)
}

func (i *Interface) Methods() *Functions {
	i.loadAllMethods()
	return &Functions{
		elements: i.allMethods,
	}
}

func (i *Interface) IsConstraint() bool {
	i.loadEmbeddedInterfaces()

	// diff is greater than 1 means that there are many non-interface types defined in interface
	// and this constraint can never be satisfied
	if len(i.embeddedTypes)-len(i.embeddedInterfaces) > 1 {
		return false
	}

	hasConstraintTypes := len(i.embeddedTypes)-len(i.embeddedInterfaces) == 1

	for _, embeddedInterface := range i.embeddedInterfaces {
		if embeddedInterface.IsConstraint() {
			if hasConstraintTypes {
				return false
			} else {
				hasConstraintTypes = true
			}
		}
	}

	return hasConstraintTypes
}

func (i *Interface) TypeParameters() *TypeParameters {
	i.loadTypeParams()
	return i.typeParams
}

func (i *Interface) String() string {
	if i.name == "" && len(i.fieldList) == 0 {
		return "interface{}"
	}

	var builder strings.Builder
	if i.file != nil && i.file.pkg.Name != "builtin" {
		builder.WriteString(fmt.Sprintf("%s.%s", i.file.Package().Name, i.name))
	} else if i.name != "" {
		builder.WriteString(i.name)
	}

	if i.TypeParameters().Len() != 0 {
		builder.WriteString("[")

		for index := 0; index < i.TypeParameters().Len(); index++ {
			typeParam := i.TypeParameters().At(index)
			builder.WriteString(typeParam.String())

			if index != i.TypeParameters().Len()-1 {
				builder.WriteString(",")
			}
		}

		builder.WriteString("]")
	}

	return builder.String()
}

func (i *Interface) InterfaceType() *types.Interface {
	return i.interfaceType
}

func (i *Interface) getInterfaceMethods() []*Function {
	methods := make([]*Function, 0)

	markers := i.visitor.allPackageMarkers[i.pkg.ID]

	for _, rawMethod := range i.fieldList {
		_, ok := rawMethod.Type.(*ast.FuncType)

		if ok {
			methods = append(methods, newFunction(nil, nil, rawMethod, i, i.file, i.pkg, i.visitor, markers[rawMethod]))
		}
	}

	return methods
}

func (i *Interface) getEmbeddedTypes() []Type {
	embeddedTypes := make([]Type, 0)

	for _, field := range i.fieldList {
		_, ok := field.Type.(*ast.FuncType)

		if !ok {
			embeddedTypes = append(embeddedTypes, getTypeFromExpression(field.Type, i.file, i.visitor, nil, i.typeParams))
		}
	}

	return embeddedTypes
}

func (i *Interface) getEmbeddedInterfaces() []*Interface {
	embeddedInterfaces := make([]*Interface, 0)

	for _, embeddedType := range i.embeddedTypes {
		if iface, isInterface := embeddedType.(*Interface); isInterface {
			embeddedInterfaces = append(embeddedInterfaces, iface)
		}
	}

	return embeddedInterfaces
}

func (i *Interface) loadEmbeddedTypes() {
	i.embeddedTypesOnce.Do(func() {
		i.loadTypeParams()
		i.embeddedTypes = i.getEmbeddedTypes()
	})
}

func (i *Interface) loadEmbeddedInterfaces() {
	i.embeddedInterfacesOnce.Do(func() {
		i.loadEmbeddedTypes()
		i.embeddedInterfaces = i.getEmbeddedInterfaces()
	})
}

func (i *Interface) loadMethods() {
	i.methodsOnce.Do(func() {
		i.loadTypeParams()
		i.methods = i.getInterfaceMethods()
		i.allMethods = append(i.allMethods, i.methods...)
	})
}

func (i *Interface) loadTypeParams() {
	i.typeParamsOnce.Do(func() {
		if i.specType == nil || i.specType.TypeParams == nil {
			return
		}

		for _, field := range i.specType.TypeParams.List {
			for _, fieldName := range field.Names {
				typeParameter := &TypeParameter{
					name: fieldName.Name,
					constraints: &TypeConstraints{
						[]*TypeConstraint{},
					},
				}
				i.typeParams.elements = append(i.typeParams.elements, typeParameter)
			}
		}

		for _, field := range i.specType.TypeParams.List {
			constraints := make([]*TypeConstraint, 0)
			typ := getTypeFromExpression(field.Type, i.file, i.visitor, nil, i.typeParams)

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
				typeParam, exists := i.typeParams.FindByName(fieldName.Name)

				if exists {
					typeParam.constraints.elements = append(typeParam.constraints.elements, constraints...)
				}
			}
		}

	})
}

func (i *Interface) loadAllMethods() {
	i.allMethodsOnce.Do(func() {
		i.loadMethods()
		i.loadEmbeddedInterfaces()

		for _, embeddedInterface := range i.embeddedInterfaces {
			embeddedInterface.loadAllMethods()
			i.allMethods = append(i.allMethods, embeddedInterface.allMethods...)
		}
	})
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
