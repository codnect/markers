package visitor

import (
	"errors"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"strconv"
)

type FileCallback func(file *File, err error) error

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

type packageVisitor struct {
	collector *packageCollector

	pkg               *packages.Package
	packageMarkers    map[ast.Node]marker.MarkerValues
	allPackageMarkers map[string]map[ast.Node]marker.MarkerValues

	file *File

	genDecl  *ast.GenDecl
	funcDecl *ast.FuncDecl
	rawFile  *ast.File
	typeSpec *ast.TypeSpec
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
		visitor.file = newFile(typedNode, visitor.pkg, visitor.packageMarkers[typedNode])
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
		visitor.collectFunction()
		return nil
	case *ast.TypeSpec:
		visitor.typeSpec = typedNode
		visitor.getTypeFromTypeSpec()
		return nil
	default:
		return nil
	}
}

func collectConstantsFromSpecs(specs []ast.Spec, file *File) {
	var last *ast.ValueSpec
	for iota, s := range specs {
		valueSpec := s.(*ast.ValueSpec)

		switch {
		case valueSpec.Type != nil || len(valueSpec.Values) > 0:
			last = valueSpec
		case last == nil:
			last = new(ast.ValueSpec)
		}

		collectConstants(valueSpec, last, iota, file)
	}
}

func collectConstants(valueSpec *ast.ValueSpec, lastValueSpec *ast.ValueSpec, iota int, file *File) {
	for _, name := range valueSpec.Names {
		constant := &Constant{
			name:       name.Name,
			isExported: ast.IsExported(name.Name),
			iota:       iota,
			pkg:        file.pkg,
			//visitor:    visitor,
		}

		if valueSpec.Values != nil {
			constant.expression = valueSpec.Values[0]
			constant.initType = valueSpec.Type
		} else {
			constant.expression = lastValueSpec.Values[0]
			constant.initType = lastValueSpec.Type
		}

		file.constants.elements = append(file.constants.elements, constant)
	}
}

func (visitor *packageVisitor) collectFunction() {
	function := &Function{
		name:       visitor.funcDecl.Name.Name,
		isExported: ast.IsExported(visitor.funcDecl.Name.Name),
		file:       visitor.currentFile,
		position:   visitor.getPosition(visitor.funcDecl.Pos()),
		params:     &Tuple{},
		results:    &Tuple{},
	}

	funcType := visitor.funcDecl.Type

	if funcType.Params != nil {
		function.params.variables = append(function.params.variables, visitor.getVariables(funcType.Params.List).variables...)
	}

	if funcType.Results != nil {
		function.results.variables = append(function.results.variables, visitor.getVariables(funcType.Results.List).variables...)
	}

	if visitor.funcDecl.Recv == nil {
		visitor.currentFile.functions.elements = append(visitor.currentFile.functions.elements, function)
	} else {
		receiverVariable := &Variable{}

		if visitor.funcDecl.Recv.List[0].Names != nil {
			receiverVariable.name = visitor.funcDecl.Recv.List[0].Names[0].Name
		}

		var receiverTypeSpec *ast.TypeSpec
		receiver := visitor.funcDecl.Recv.List[0].Type

		receiverTypeName := ""
		isPointerReceiver := false
		isStructMethod := false

		switch typedReceiver := receiver.(type) {
		case *ast.Ident:
			if typedReceiver.Obj == nil {
				receiverTypeName = typedReceiver.Name
				unprocessedype := visitor.getTypeFromScope(receiverTypeName)
				_, isStructMethod = unprocessedype.(*Struct)
			} else {
				receiverTypeSpec = typedReceiver.Obj.Decl.(*ast.TypeSpec)
				receiverTypeName = receiverTypeSpec.Name.Name
				_, isStructMethod = receiverTypeSpec.Type.(*ast.StructType)
			}
		case *ast.StarExpr:
			if typedReceiver.X.(*ast.Ident).Obj == nil {
				receiverTypeName = typedReceiver.X.(*ast.Ident).Name
				unprocessedype := visitor.getTypeFromScope(receiverTypeName)
				_, isStructMethod = unprocessedype.(*Struct)
			} else {
				receiverTypeSpec = typedReceiver.X.(*ast.Ident).Obj.Decl.(*ast.TypeSpec)
				receiverTypeName = receiverTypeSpec.Name.Name
				_, isStructMethod = receiverTypeSpec.Type.(*ast.StructType)
			}
			isPointerReceiver = true
		}

		candidateType, ok := visitor.collector.findTypeByPkgIdAndName(visitor.pkg.ID, receiverTypeName)

		if isStructMethod {
			var structType *Struct

			if !ok {
				structType = visitor.getStruct(receiverTypeSpec)
				visitor.currentFile.structs.elements = append(visitor.currentFile.structs.elements, structType)
				candidateType = structType
			} else {
				structType = candidateType.(*Struct)
			}

			structType.methods = append(structType.methods, function)
		} else {
			var customType *CustomType

			if !ok {
				customType = visitor.getCustomType(receiverTypeSpec)
				visitor.currentFile.customTypes.elements = append(visitor.currentFile.customTypes.elements, customType)
				candidateType = customType
			} else {
				customType = candidateType.(*CustomType)
			}

			customType.methods = append(customType.methods, function)
		}

		if isPointerReceiver {
			receiverVariable.typ = &Pointer{
				base: candidateType,
			}
		} else {
			receiverVariable.typ = candidateType
		}

		function.receiver = receiverVariable
	}
}

func getTypeFromScope(name string, pkg *packages.Package, collector *packageCollector) Type {
	typ := pkg.Types.Scope().Lookup(name)

	typedName, ok := typ.Type().(*types.Named)

	if _, ok := collector.unprocessedTypes[pkg.ID]; !ok {
		collector.unprocessedTypes[pkg.ID] = make(map[string]Type)
	}

	if ok {
		switch typedName.Underlying().(type) {
		case *types.Struct:
			structType := &Struct{
				name:        name,
				isProcessed: false,
			}
			collector.unprocessedTypes[pkg.ID][name] = structType
			return structType
		case *types.Interface:
			interfaceType := &Interface{
				name:        name,
				isProcessed: false,
			}
			collector.unprocessedTypes[pkg.ID][name] = interfaceType
			return interfaceType
		default:
			customType := &CustomType{
				name:        name,
				isProcessed: false,
			}
			collector.unprocessedTypes[pkg.ID][name] = customType
			return customType
		}
	}

	return nil
}

func (visitor *packageVisitor) getTypeFromTypeSpec() Type {
	typeName := visitor.typeSpec.Name.Name

	switch visitor.typeSpec.Type.(type) {
	case *ast.InterfaceType:
		interfaceCandidate, ok := visitor.collector.findTypeByPkgIdAndName(visitor.pkg.ID, typeName)

		var interfaceType *Interface
		if ok {
			interfaceType = interfaceCandidate.(*Interface)

			if interfaceType.isProcessed {
				//interfaceType.rawGenDecl = visitor.genDecl
			} else {
				visitor.processInterface(interfaceType, visitor.typeSpec)
				visitor.currentFile.interfaces.elements = append(visitor.currentFile.interfaces.elements, interfaceType)
			}
		} else {
			interfaceType = visitor.getInterface(visitor.typeSpec)
			visitor.currentFile.interfaces.elements = append(visitor.currentFile.interfaces.elements, interfaceType)
		}

		return interfaceType
	case *ast.StructType:
		structCandidate, ok := visitor.collector.findTypeByPkgIdAndName(visitor.pkg.ID, typeName)

		var structType *Struct
		if ok {
			structType = structCandidate.(*Struct)

			if structType.isProcessed {
				//structType.rawGenDecl = visitor.genDecl
			} else {
				visitor.processStruct(structType, visitor.typeSpec)
				visitor.currentFile.structs.elements = append(visitor.currentFile.structs.elements, structType)
			}
		} else {
			structType = visitor.getStruct(visitor.typeSpec)
			visitor.currentFile.structs.elements = append(visitor.currentFile.structs.elements, structType)
		}

		return structType
	}

	customTypeCandidate, ok := visitor.collector.findTypeByPkgIdAndName(visitor.pkg.ID, typeName)

	var customType *CustomType
	if ok {
		defer func() {
			if r := recover(); r != nil {
				log.Printf(customTypeCandidate.String())
				log.Printf(visitor.pkg.ID)
				log.Printf(typeName)
			}
		}()
		customType = customTypeCandidate.(*CustomType)

		if customType.isProcessed {
			//customType.rawGenDecl = visitor.genDecl
		} else {
			visitor.processCustomType(customType, visitor.typeSpec)
			visitor.currentFile.customTypes.elements = append(visitor.currentFile.customTypes.elements, customType)
		}
	} else {
		customType := visitor.getCustomType(visitor.typeSpec)
		visitor.currentFile.customTypes.elements = append(visitor.currentFile.customTypes.elements, customType)
	}

	return customType
}

func getTypeFromExpression(expr ast.Expr, pkg *packages.Package, collector *packageCollector) Type {

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
			return visitor.getTypeFromScope(typed.Name)
		}

		visitor.typeSpec = typed.Obj.Decl.(*ast.TypeSpec)
		return visitor.getTypeFromTypeSpec()
	case *ast.SelectorExpr:
		importName := typed.X.(*ast.Ident).Name
		typeName := typed.Sel.Name
		return visitor.findTypeByImportAndTypeName(importName, typeName)
	case *ast.StarExpr:
		return &Pointer{
			base: getTypeFromExpression(typed.X, pkg, collector),
		}
	case *ast.ArrayType:

		if typed.Len == nil {
			return &Slice{
				elem: getTypeFromExpression(typed.Elt, pkg, collector),
			}
		} else {
			basicLit, isBasicLit := typed.Len.(*ast.BasicLit)

			if isBasicLit {
				length, _ := strconv.ParseInt(basicLit.Value, 10, 64)
				return &Array{
					elem: getTypeFromExpression(typed.Elt, pkg, collector),
					len:  length,
				}
			}

			return &Array{
				elem: getTypeFromExpression(typed.Elt, pkg, collector),
				len:  -1,
			}
		}
	case *ast.ChanType:
		chanType := &Chan{
			elem: getTypeFromExpression(typed.Value, pkg, collector),
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
			key:  getTypeFromExpression(typed.Key, pkg, collector),
			elem: getTypeFromExpression(typed.Value, pkg, collector),
		}
	case *ast.InterfaceType:
		interfaceType := &Interface{
			position: getPosition(pkg, typed.Pos()),
			//visitor:     visitor,
			isAnonymous: true,
			isProcessed: true,
		}
		//interfaceType.methods = append(interfaceType.methods, visitor.getInterfaceMethods(typed.Methods.List)...)
		return interfaceType
	case *ast.StructType:
		structType := &Struct{
			position:    getPosition(pkg, typed.Pos()),
			isAnonymous: true,
			isProcessed: true,
		}

		return structType
	}

	return nil
}

func (visitor *packageVisitor) findTypeByImportAndTypeName(importName, typeName string) *ImportedType {
	if importedType, ok := visitor.collector.importTypes[importName+"#"+typeName]; ok {
		return importedType
	}

	packageImport, _ := visitor.currentFile.imports.FindByName(importName)

	if packageImport == nil {
		packageImport, _ = visitor.currentFile.imports.FindByPath(importName)
	}

	if importedType, ok := visitor.collector.importTypes[packageImport.path+"#"+typeName]; ok {
		return importedType
	}

	typ, exists := visitor.collector.findTypeByPkgIdAndName(packageImport.path, typeName)

	if exists {
		importedType := &ImportedType{
			visitor.collector.packages[packageImport.path],
			typ,
		}
		visitor.collector.importTypes[packageImport.path+"#"+typeName] = importedType
	}

	importedType := &ImportedType{
		pkg: visitor.collector.packages[packageImport.path],
		typ: typ,
	}
	visitor.collector.importTypes[packageImport.path+"#"+typeName] = importedType
	return importedType
}

func getVariables(fieldList []*ast.Field) *Tuple {
	tuple := &Tuple{
		variables: make([]*Variable, 0),
	}

	for _, field := range fieldList {

		typ := getTypeFromExpression(field.Type)

		if field.Names == nil {
			tuple.variables = append(tuple.variables, &Variable{
				typ: typ,
			})
		}

		for _, fieldName := range field.Names {
			tuple.variables = append(tuple.variables, &Variable{
				name: fieldName.Name,
				typ:  typ,
			})
		}

	}

	return tuple
}

func (visitor *packageVisitor) getInterfaceMethods(fieldList []*ast.Field) []*Function {
	methods := make([]*Function, 0)

	for _, rawMethod := range fieldList {
		funcType, ok := rawMethod.Type.(*ast.FuncType)

		if ok {
			method := &Function{
				params:  &Tuple{},
				results: &Tuple{},
			}

			if rawMethod.Names != nil {
				method.name = rawMethod.Names[0].Name
			}

			if funcType.Params != nil {
				method.params.variables = append(method.params.variables, visitor.getVariables(funcType.Params.List).variables...)
			}

			if funcType.Results != nil {
				method.results.variables = append(method.results.variables, visitor.getVariables(funcType.Results.List).variables...)
			}

			methods = append(methods, method)
		}
	}

	return methods
}

func (visitor *packageVisitor) getInterfaceEmbeddedTypes(fieldList []*ast.Field) []Type {
	embeddedTypes := make([]Type, 0)

	for _, field := range fieldList {
		_, ok := field.Type.(*ast.FuncType)

		if !ok {
			embeddedTypes = append(embeddedTypes, visitor.getTypeFromExpression(field.Type))
		}
	}

	return embeddedTypes
}

func (visitor *packageVisitor) getFieldsFromFieldList(fieldList []*ast.Field) []*Field {
	fields := make([]*Field, 0)

	for _, rawField := range fieldList {
		tags := ""

		if rawField.Tag != nil {
			tags = rawField.Tag.Value
		}

		if rawField.Names == nil {
			embeddedType := visitor.getTypeFromExpression(rawField.Type)

			field := &Field{
				name:       "",
				isExported: false,
				position:   Position{},
				markers:    visitor.packageMarkers[rawField],
				file:       visitor.currentFile,
				tags:       tags,
				typ:        embeddedType,
				isEmbedded: true,
			}

			fields = append(fields, field)
			continue
		}

		for _, fieldName := range rawField.Names {
			typ := visitor.getTypeFromExpression(rawField.Type)

			field := &Field{
				name:       fieldName.Name,
				isExported: ast.IsExported(fieldName.Name),
				position:   visitor.getPosition(fieldName.Pos()),
				markers:    visitor.packageMarkers[rawField],
				file:       visitor.currentFile,
				tags:       tags,
				typ:        typ,
				isEmbedded: false,
			}

			fields = append(fields, field)
		}

	}

	return fields
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

	sourceFiles := collector.files[pkg.ID]

	if sourceFiles == nil {
		return
	}

	for i := 0; i < sourceFiles.Len(); i++ {
		file := sourceFiles.At(i)
		for j := 0; j < file.constants.Len(); j++ {
			constant := file.constants.At(j)
			constant.evaluateExpression()
		}
	}

}
