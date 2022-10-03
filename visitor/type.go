package visitor

import (
	"fmt"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
	"strings"
)

type Type interface {
	Name() string
	Underlying() Type
	String() string
}

type Position struct {
	Line   int
	Column int
}

func getPosition(pkg *packages.Package, tokenPosition token.Pos) Position {
	position := pkg.Fset.Position(tokenPosition)
	return Position{
		Line:   position.Line,
		Column: position.Column,
	}
}

type ImportedType struct {
	pkg *packages.Package
	typ Type
}

func (i *ImportedType) Package() *packages.Package {
	return i.pkg
}

func (i *ImportedType) Underlying() Type {
	return i.typ
}

func (i *ImportedType) String() string {
	return ""
}

func (i *ImportedType) Name() string {
	return fmt.Sprintf("%s.%s", i.pkg.Name, i.typ.Name())
}

type Variadic struct {
	elem Type
}

func (v *Variadic) Name() string {
	return v.elem.Name()
}

func (v *Variadic) Elem() Type {
	return v.elem
}

func (v *Variadic) Underlying() Type {
	return v
}

func (v *Variadic) String() string {
	return ""
}

type Pointer struct {
	base Type
}

func (p *Pointer) Name() string {
	return ""
}

func (p *Pointer) Elem() Type {
	return p.base
}

func (p *Pointer) Underlying() Type {
	return p
}

func (p *Pointer) String() string {
	var builder strings.Builder
	builder.WriteString("*")
	builder.WriteString(p.base.Name())
	return builder.String()
}

type TypeParam struct {
	name string
	typ  Type
}

func (t *TypeParam) Name() string {
	return t.name
}

func (t *TypeParam) Type() Type {
	return t.typ
}

type TypeParams struct {
	params []*TypeParam
}

func (t *TypeParams) Len() int {
	return len(t.params)
}

func (t *TypeParams) At(index int) *TypeParam {
	if index >= 0 && index < len(t.params) {
		return t.params[index]
	}

	return nil
}

type Generic struct {
	typeParam *TypeParam
}

func (g *Generic) Name() string {
	return g.typeParam.name
}

func (g *Generic) ParamName() string {
	return g.typeParam.name
}

func (g *Generic) TypeParam() *TypeParam {
	return g.typeParam
}

func (g *Generic) Underlying() Type {
	return g.typeParam.typ
}

func (g *Generic) String() string {
	return ""
}

func getTypeFromScope(name string, visitor *packageVisitor) Type {
	pkg := visitor.pkg
	typ := pkg.Types.Scope().Lookup(name)

	typedName, ok := typ.Type().(*types.Named)

	if _, ok := visitor.collector.unprocessedTypes[pkg.ID]; !ok {
		visitor.collector.unprocessedTypes[pkg.ID] = make(map[string]Type)
	}

	if ok {
		switch typedName.Underlying().(type) {
		case *types.Struct:
			structType := &Struct{
				name:        name,
				isProcessed: false,
			}
			visitor.collector.unprocessedTypes[pkg.ID][name] = structType
			return structType
		case *types.Interface:
			interfaceType := &Interface{
				name:        name,
				isProcessed: false,
			}
			visitor.collector.unprocessedTypes[pkg.ID][name] = interfaceType
			return interfaceType
		default:
			customType := &CustomType{
				name:        name,
				isProcessed: false,
			}
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
				file.interfaces.elements = append(file.interfaces.elements, t)
			}
			t.markers = visitor.packageMarkers[typeSpec]
			return t
		case *Struct:
			if !t.isProcessed {
				file.structs.elements = append(file.structs.elements, t)
			}
			t.markers = visitor.packageMarkers[typeSpec]
			return t
		case *CustomType:
			if !t.isProcessed {
				file.customTypes.elements = append(file.customTypes.elements, t)
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

func getTypeFromExpression(expr ast.Expr, visitor *packageVisitor) Type {
	file := visitor.file
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

		if typed.Name == "error" {
			errorType, _ := collector.findTypeByPkgIdAndName("builtin", "error")
			return errorType
		} else if typed.Name == "any" {
			anyType, _ := collector.findTypeByPkgIdAndName("builtin", "any")
			return anyType
		}

		typ, ok = collector.findTypeByPkgIdAndName(pkg.ID, typed.Name)

		if ok {
			return typ
		}

		if typed.Obj == nil {
			return getTypeFromScope(typed.Name, visitor)
		}

		return collectTypeFromTypeSpec(typed.Obj.Decl.(*ast.TypeSpec), visitor)
	case *ast.SelectorExpr:
		importName := typed.X.(*ast.Ident).Name
		typeName := typed.Sel.Name
		return collector.findTypeByImportAndTypeName(importName, typeName, file)
	case *ast.StarExpr:
		return &Pointer{
			base: getTypeFromExpression(typed.X, visitor),
		}
	case *ast.ArrayType:

		if typed.Len == nil {
			return &Slice{
				elem: getTypeFromExpression(typed.Elt, visitor),
			}
		} else {
			basicLit, isBasicLit := typed.Len.(*ast.BasicLit)

			if isBasicLit {
				length, _ := strconv.ParseInt(basicLit.Value, 10, 64)
				return &Array{
					elem: getTypeFromExpression(typed.Elt, visitor),
					len:  length,
				}
			}

			return &Array{
				elem: getTypeFromExpression(typed.Elt, visitor),
				len:  -1,
			}
		}
	case *ast.ChanType:
		chanType := &Chan{
			elem: getTypeFromExpression(typed.Value, visitor),
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
			elem: getTypeFromExpression(typed.Elt, visitor),
		}
	case *ast.FuncType:
		return &Function{}
	case *ast.MapType:
		return &Map{
			key:  getTypeFromExpression(typed.Key, visitor),
			elem: getTypeFromExpression(typed.Value, visitor),
		}
	case *ast.InterfaceType:
		return newInterface(nil, typed, nil, nil, visitor, nil)
	case *ast.StructType:
		return newStruct(nil, typed, nil, nil, visitor, nil)
	}

	return nil
}
