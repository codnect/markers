package visitor

import (
	"errors"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
)

type FileCallback func(file *File, err error) error

type packageVisitor struct {
	collector *packageCollector

	pkg               *packages.Package
	packageMarkers    map[ast.Node]marker.MarkerValues
	allPackageMarkers map[string]map[ast.Node]marker.MarkerValues

	file *File

	genDecl  *ast.GenDecl
	funcDecl *ast.FuncDecl
	rawFile  *ast.File
}

func (visitor *packageVisitor) VisitPackage() {
	visitor.packageMarkers = visitor.allPackageMarkers[visitor.pkg.ID]
	visitor.collector.markAsSeen(visitor.pkg.ID)

	for _, file := range visitor.pkg.Syntax {
		ast.Walk(visitor, file)
	}
}

func (visitor *packageVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return visitor
	}

	switch typedNode := node.(type) {
	case *ast.File:
		visitor.rawFile = typedNode
		visitor.file = newFile(typedNode, visitor.pkg, visitor.packageMarkers[typedNode], visitor)
		visitor.collector.addFile(visitor.pkg.ID, visitor.file)
		return visitor
	case *ast.GenDecl:
		visitor.genDecl = typedNode

		if typedNode.Tok == token.CONST {
			collectConstantsFromSpecs(typedNode.Specs, visitor.file)
		}

		return visitor
	case *ast.FuncDecl:
		visitor.funcDecl = typedNode
		newFunction(typedNode, nil, visitor.file)
		return nil
	case *ast.TypeSpec:
		collectTypeFromTypeSpec(typedNode, visitor)
		return nil
	default:
		return nil
	}
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
			return t
		case *Struct:
			if !t.isProcessed {
				file.structs.elements = append(file.structs.elements, t)
			}
			return t
		case *CustomType:
			if !t.isProcessed {
				file.customTypes.elements = append(file.customTypes.elements, t)
			}
			return t
		}
	}

	switch typeSpec.Type.(type) {
	case *ast.InterfaceType:
		return newInterface(typeSpec, nil, file, pkg, nil)
	case *ast.StructType:
		return newStruct(typeSpec, nil, file, pkg, nil)
	default:
		return newCustomType(typeSpec, file, pkg, nil, nil)
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
		}

		typ, ok = collector.findTypeByPkgIdAndName(pkg.ID, typed.Name)

		if ok {
			return typ
		}

		if typed.Obj == nil {
			return getTypeFromScope(typed.Name, visitor)
		}

		return collectTypeFromTypeSpec(typed.Obj.Decl.(*ast.TypeSpec), nil)
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
			key:  getTypeFromExpression(typed.Key, visitor),
			elem: getTypeFromExpression(typed.Value, visitor),
		}
	case *ast.InterfaceType:
		return newInterface(nil, typed, nil, nil, nil)
	case *ast.StructType:
		return newStruct(nil, typed, nil, nil, nil)
	}

	return nil
}

func visitPackage(pkg *packages.Package, collector *packageCollector, allPackageMarkers map[string]map[ast.Node]marker.MarkerValues) {
	pkgVisitor := &packageVisitor{
		collector:         collector,
		pkg:               pkg,
		allPackageMarkers: allPackageMarkers,
	}

	if _, ok := collector.packages[pkg.ID]; !ok {
		collector.packages[pkg.ID] = pkg
	}

	pkgVisitor.VisitPackage()
}

func EachFile(collector *marker.Collector, pkgs []*packages.Package, callback FileCallback) error {
	if collector == nil {
		return errors.New("collector cannot be nil")
	}

	if pkgs == nil {
		return errors.New("packages cannot be nil")
	}

	var errs []error
	packageMarkers := make(map[string]map[ast.Node]marker.MarkerValues)

	for _, pkg := range pkgs {
		markers, err := collector.Collect(pkg)

		if err != nil {
			errs = append(errs, err.(marker.ErrorList)...)
			continue
		}

		packageMarkers[pkg.ID] = markers
	}

	if len(errs) != 0 {
		return marker.NewErrorList(errs)
	}

	pkgCollector := newPackageCollector()

	for _, pkg := range pkgs {
		if !pkgCollector.isVisited(pkg.ID) || !pkgCollector.isProcessed(pkg.ID) {
			visitPackage(pkg, pkgCollector, packageMarkers)
		}
	}

	return nil
}
