package visitor

import (
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
)

type Type interface {
	Name() string
	Underlying() Type
	String() string
}

type Types struct {
	elements []Type
}

func (t *Types) Len() int {
	return len(t.elements)
}

func (t *Types) At(index int) Type {
	if index >= 0 && index < len(t.elements) {
		return t.elements[index]
	}

	return nil
}

func (t *Types) FindByName(name string) (Type, bool) {
	for _, typ := range t.elements {
		if typ.Name() == name {
			return typ, true
		}
	}

	return nil, false
}

type TypeSets []Type

func (t TypeSets) Name() string {
	return ""
}

func (t TypeSets) Len() int {
	return len(t)
}

func (t TypeSets) At(index int) Type {
	if index >= 0 && index < len(t) {
		return t[index]
	}

	return nil
}

func (t TypeSets) Underlying() Type {
	return t
}

func (t TypeSets) String() string {
	return ""
}

func getTypeFromScope(name string, visitor *packageVisitor) Type {
	pkg := visitor.pkg
	typ := pkg.Types.Scope().Lookup(name)

	typedName, isNamedType := typ.Type().(*types.Named)

	if _, ok := visitor.collector.unprocessedTypes[pkg.ID]; !ok {
		visitor.collector.unprocessedTypes[pkg.ID] = make(map[string]Type)
	}

	if unprocessedType, ok := visitor.collector.unprocessedTypes[pkg.ID][name]; ok {
		return unprocessedType
	}

	if isNamedType {
		switch typedName.Underlying().(type) {
		case *types.Struct:
			structType := newStruct(nil, nil, nil, pkg, visitor, nil)
			structType.isProcessed = false
			visitor.collector.unprocessedTypes[pkg.ID][name] = structType
			return structType
		case *types.Interface:
			interfaceType := newInterface(nil, nil, nil, pkg, visitor, nil)
			interfaceType.isProcessed = false
			visitor.collector.unprocessedTypes[pkg.ID][name] = interfaceType
			return interfaceType
		default:
			customType := newCustomType(nil, nil, pkg, visitor, nil)
			customType.isProcessed = false
			visitor.collector.unprocessedTypes[pkg.ID][name] = customType
			return customType
		}
	}

	return nil
}

func collectTypeFromTypeSpec(typeSpec *ast.TypeSpec, visitor *packageVisitor) Type {
	file := visitor.file
	pkg := visitor.pkg
	typeName := typeSpec.Name.Name

	typ, ok := visitor.collector.findTypeByPkgIdAndName(pkg.ID, typeName)

	if ok {
		switch t := typ.(type) {
		case *Interface:
			if !t.isProcessed {
				if _, exists := file.interfaces.FindByName(t.name); !exists {
					file.interfaces.elements = append(file.interfaces.elements, t)
				}
				t.initialize(typeSpec, nil, file, pkg)
			}
			t.markers = visitor.packageMarkers[typeSpec]
			return t
		case *Struct:
			if !t.isProcessed {
				if _, exists := file.structs.FindByName(t.name); !exists {
					file.structs.elements = append(file.structs.elements, t)
				}
				t.initialize(typeSpec, nil, file, pkg)
			}
			t.markers = visitor.packageMarkers[typeSpec]
			return t
		case *CustomType:
			if !t.isProcessed {
				if _, exists := file.customTypes.FindByName(t.name); !exists {
					file.customTypes.elements = append(file.customTypes.elements, t)
				}
				t.initialize(typeSpec, file, pkg)

			}
			t.markers = visitor.packageMarkers[typeSpec]
			return t
		}
	}

	switch typeSpec.Type.(type) {
	case *ast.InterfaceType:
		return newInterface(typeSpec, nil, file, pkg, visitor, visitor.packageMarkers[typeSpec])
	case *ast.StructType:
		return newStruct(typeSpec, nil, file, pkg, visitor, visitor.packageMarkers[typeSpec])
	default:
		return newCustomType(typeSpec, file, pkg, visitor, visitor.packageMarkers[typeSpec])
	}
}

func getTypeFromExpression(expr ast.Expr, file *File, visitor *packageVisitor, ownerType Type, typeParameters *TypeParameters) Type {
	pkg := visitor.pkg
	collector := visitor.collector

	switch typed := expr.(type) {
	case *ast.Ident:
		var typ Type
		var ok bool
		typ, ok = basicTypesMap[typed.Name]

		if ok {
			return typ
		}

		typ, ok = collector.findTypeByPkgIdAndName("builtin", typed.Name)

		if ok {
			return typ
		}

		typ, ok = collector.findTypeByPkgIdAndName(pkg.ID, typed.Name)

		if ok {
			return typ
		}

		if typed.Obj == nil {
			if typeParameters != nil {
				if typeParameter, exists := typeParameters.FindByName(typed.Name); exists {
					return typeParameter
				}
			}
			return getTypeFromScope(typed.Name, visitor)
		}

		if field, isField := typed.Obj.Decl.(*ast.Field); isField {
			if typeParameters == nil {
				//TODO: return invalid type
				return nil
			}

			if typeParameter, exists := typeParameters.FindByName(field.Names[0].Name); exists {
				return typeParameter
			}

			//TODO: return invalid type
			return nil
		}

		return collectTypeFromTypeSpec(typed.Obj.Decl.(*ast.TypeSpec), visitor)
	case *ast.SelectorExpr:
		importName := typed.X.(*ast.Ident).Name
		typeName := typed.Sel.Name
		return collector.findTypeByImportAndTypeName(importName, typeName, file)
	case *ast.StarExpr:
		return &Pointer{
			base: getTypeFromExpression(typed.X, file, visitor, ownerType, typeParameters),
		}
	case *ast.ArrayType:

		if typed.Len == nil {
			return &Slice{
				elem: getTypeFromExpression(typed.Elt, file, visitor, ownerType, typeParameters),
			}
		} else {
			basicLit, isBasicLit := typed.Len.(*ast.BasicLit)

			if isBasicLit {
				length, _ := strconv.ParseInt(basicLit.Value, 10, 64)
				return &Array{
					elem: getTypeFromExpression(typed.Elt, file, visitor, ownerType, typeParameters),
					len:  length,
				}
			}

			return &Array{
				elem: getTypeFromExpression(typed.Elt, file, visitor, ownerType, typeParameters),
				len:  -1,
			}
		}
	case *ast.ChanType:
		chanType := &Chan{
			elem: getTypeFromExpression(typed.Value, file, visitor, ownerType, typeParameters),
		}

		if typed.Dir&ast.SEND == ast.SEND {
			chanType.direction |= SendDir
		}

		if typed.Dir&ast.RECV == ast.RECV {
			chanType.direction |= ReceiveDir
		}

		return chanType
	case *ast.Ellipsis:
		return &Variadic{
			elem: getTypeFromExpression(typed.Elt, file, visitor, ownerType, typeParameters),
		}
	case *ast.FuncType:
		return newFunction(nil, typed, nil, ownerType, file, nil, visitor, nil)
	case *ast.MapType:
		return &Map{
			key:  getTypeFromExpression(typed.Key, file, visitor, ownerType, typeParameters),
			elem: getTypeFromExpression(typed.Value, file, visitor, ownerType, typeParameters),
		}
	case *ast.InterfaceType:
		return newInterface(nil, typed, nil, nil, visitor, nil)
	case *ast.StructType:
		return newStruct(nil, typed, nil, nil, visitor, nil)
	case *ast.IndexExpr:
		genericType := &GenericType{
			rawType:   getTypeFromExpression(typed.X, file, visitor, ownerType, typeParameters),
			arguments: make([]Type, 0),
		}

		genericType.arguments = append(genericType.arguments, getTypeFromExpression(typed.Index, file, visitor, ownerType, typeParameters))
		return genericType
	case *ast.IndexListExpr:
		genericType := &GenericType{
			rawType:   getTypeFromExpression(typed.X, file, visitor, ownerType, typeParameters),
			arguments: make([]Type, 0),
		}

		for _, argument := range typed.Indices {
			genericType.arguments = append(genericType.arguments, getTypeFromExpression(argument, file, visitor, ownerType, typeParameters))
		}

		return genericType
	case *ast.BinaryExpr:
		constraints := make(TypeSets, 0)

		firstType := getTypeFromExpression(typed.X, file, visitor, ownerType, typeParameters)

		if typeSets, isTypeSet := firstType.(TypeSets); isTypeSet {
			for _, typ := range typeSets {
				if constraint, isConstraint := typ.(*TypeConstraint); isConstraint {
					constraints = append(constraints, constraint)
				} else {
					constraints = append(constraints, &TypeConstraint{typ: typ})
				}
			}
		} else {
			if constraint, isConstraint := firstType.(*TypeConstraint); isConstraint {
				constraints = append(constraints, constraint)
			} else {
				constraints = append(constraints, &TypeConstraint{typ: firstType})
			}
		}

		secondType := getTypeFromExpression(typed.Y, file, visitor, ownerType, typeParameters)

		if typeSets, isTypeSet := secondType.(TypeSets); isTypeSet {
			for _, typ := range typeSets {
				if constraint, isConstraint := typ.(*TypeConstraint); isConstraint {
					constraints = append(constraints, constraint)
				} else {
					constraints = append(constraints, &TypeConstraint{typ: typ})
				}
			}
		} else {
			if constraint, isConstraint := secondType.(*TypeConstraint); isConstraint {
				constraints = append(constraints, constraint)
			} else {
				constraints = append(constraints, &TypeConstraint{typ: secondType})
			}
		}

		return constraints
	case *ast.UnaryExpr:
		if typed.Op == token.TILDE {
			return &TypeConstraint{
				tildeOperator: true,
				typ:           getTypeFromExpression(typed.X, nil, visitor, ownerType, typeParameters),
			}
		}
	}

	return nil
}
