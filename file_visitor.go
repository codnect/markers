package marker

import (
	"errors"
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"strconv"
)

type SourceFileCallback func(file *SourceFile, err error)

func eachFile(collector *Collector, pkgs []*Package, callback SourceFileCallback) {
	if collector == nil {
		callback(nil, errors.New("collector cannot be nil"))
	}

	if pkgs == nil {
		callback(nil, errors.New("pkgs(packages) cannot be nil"))
	}

	//var fileMap = make(map[*ast.File]*SourceFile)
	var errs []error

	packageMarkers := make(map[string]map[ast.Node]MarkerValues)

	for _, pkg := range pkgs {
		markers, err := collector.Collect(pkg)

		if err != nil {
			errs = append(errs, err.(ErrorList)...)
			continue
		}

		packageMarkers[pkg.ID] = markers
	}

	VisitPackages(pkgs, packageMarkers)

	/*fileNodeMap := EachPackage(cache, pkg, markers)

	for fileNode, file := range fileNodeMap {
		fileMap[fileNode] = file
	}*/
}

type packageCollector struct {
	hasSeen      map[string]bool
	hasProcessed map[string]bool
	files        map[string]*SourceFiles
	packages     map[string]*Package

	unprocessedTypes map[string]map[string]T

	importTypes map[string]*ImportedType
}

func newPackageCollector() *packageCollector {
	return &packageCollector{
		hasSeen:          make(map[string]bool),
		hasProcessed:     make(map[string]bool),
		files:            make(map[string]*SourceFiles),
		packages:         make(map[string]*Package),
		unprocessedTypes: make(map[string]map[string]T),
		importTypes:      make(map[string]*ImportedType),
	}
}

func (collector *packageCollector) getPackage(pkgId string) *Package {
	return collector.packages[pkgId]
}

func (collector *packageCollector) markAsSeen(pkgId string) {
	collector.hasSeen[pkgId] = true
}

func (collector *packageCollector) markAsProcessed(pkgId string) {
	collector.hasProcessed[pkgId] = true
}

func (collector *packageCollector) isVisited(pkgId string) bool {
	visited, ok := collector.hasSeen[pkgId]

	if !ok {
		return false
	}

	return visited
}

func (collector *packageCollector) isProcessed(pkgId string) bool {
	processed, ok := collector.hasProcessed[pkgId]

	if !ok {
		return false
	}

	return processed
}

func (collector *packageCollector) addFile(pkgId string, file *SourceFile) {
	if _, ok := collector.files[pkgId]; !ok {
		collector.files[pkgId] = &SourceFiles{
			elements: make([]*SourceFile, 0),
		}
	}

	if _, ok := collector.files[pkgId].FindByName(file.name); ok {
		return
	}

	collector.files[pkgId].elements = append(collector.files[pkgId].elements, file)
}

func (collector *packageCollector) findTypeByPkgIdAndName(pkgId, typeName string) (T, bool) {
	if files, ok := collector.files[pkgId]; ok {

		for i := 0; i < files.Len(); i++ {
			file := files.At(i)

			if structType, ok := file.structs.FindByName(typeName); ok {
				return structType, true
			}

			if interfaceType, ok := file.interfaces.FindByName(typeName); ok {
				return interfaceType, true
			}

			if customType, ok := file.customTypes.FindByName(typeName); ok {
				return customType, true
			}

			if constant, ok := file.constants.FindByName(typeName); ok {
				return constant, true
			}
		}

	}

	if typ, ok := collector.unprocessedTypes[pkgId][typeName]; ok {
		return typ, true
	}

	return nil, false
}

type PackageVisitor struct {
	collector *packageCollector

	pkg               *Package
	packageMarkers    map[ast.Node]MarkerValues
	allPackageMarkers map[string]map[ast.Node]MarkerValues

	currentFile *SourceFile

	genDecl  *ast.GenDecl
	funcDecl *ast.FuncDecl
	file     *ast.File
	typeSpec *ast.TypeSpec
}

func (visitor *PackageVisitor) VisitPackage() {
	visitor.packageMarkers = visitor.allPackageMarkers[visitor.pkg.ID]
	visitor.collector.markAsSeen(visitor.pkg.ID)

	for _, file := range visitor.pkg.Syntax {
		ast.Walk(visitor, file)
	}
}

func (visitor *PackageVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return visitor
	}

	switch typedNode := node.(type) {
	case *ast.File:
		visitor.file = typedNode
		visitor.createSourceFile()
		return visitor
	case *ast.GenDecl:
		visitor.genDecl = typedNode
		if typedNode.Tok == token.CONST {
			visitor.collectConstants(typedNode.Specs)
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

func (visitor *PackageVisitor) getPosition(tokenPosition token.Pos) Position {
	position := visitor.pkg.Fset.Position(tokenPosition)
	return Position{
		Line:   position.Line,
		Column: position.Column,
	}
}

func (visitor *PackageVisitor) createSourceFile() *SourceFile {
	position := visitor.pkg.Fset.Position(visitor.file.Pos())
	fileFullPath := position.Filename

	file := &SourceFile{
		name:          filepath.Base(fileFullPath),
		fullPath:      fileFullPath,
		allMarkers:    visitor.packageMarkers[visitor.file],
		pkg:           visitor.pkg,
		imports:       visitor.getFileImports(),
		fileMarkers:   make(MarkerValues, 0),
		importMarkers: make([]ImportMarker, 0),
		functions:     &Functions{},
		structs:       &Structs{},
		interfaces:    &Interfaces{},
		customTypes:   &CustomTypes{},
		constants:     &Constants{},
	}

	for markerName, markers := range file.allMarkers {
		if ImportMarkerName == markerName {
			for _, importMarker := range markers {
				file.importMarkers = append(file.importMarkers, importMarker.(ImportMarker))
			}
		} else {
			file.fileMarkers[markerName] = append(file.fileMarkers[markerName], markers...)
		}
	}

	visitor.currentFile = file
	visitor.collector.addFile(visitor.pkg.ID, file)
	return file
}

func (visitor *PackageVisitor) getFileImports() *Imports {
	imports := &Imports{}

	for _, importPackage := range visitor.file.Imports {
		importPosition := visitor.getPosition(importPackage.Pos())
		importName := ""

		if importPackage.Name != nil {
			importName = importPackage.Name.Name
		}

		imports.imports = append(imports.imports, &Import{
			name: importName,
			path: importPackage.Path.Value[1 : len(importPackage.Path.Value)-1],
			position: Position{
				importPosition.Line,
				importPosition.Column,
			},
		})
	}

	return imports
}

func (visitor *PackageVisitor) collectConstants(specs []ast.Spec) {
	var last *ast.ValueSpec
	for iota, s := range specs {
		valueSpec := s.(*ast.ValueSpec)

		switch {
		case valueSpec.Type != nil || len(valueSpec.Values) > 0:
			last = valueSpec
		case last == nil:
			last = new(ast.ValueSpec)
		}

		visitor.getConstants(valueSpec, last, iota)
	}
}

func (visitor *PackageVisitor) getConstants(valueSpec *ast.ValueSpec, lastValueSpec *ast.ValueSpec, iota int) []*Constant {
	constants := make([]*Constant, 0)

	for _, name := range valueSpec.Names {
		constant := &Constant{
			name:       name.Name,
			isExported: ast.IsExported(name.Name),
			iota:       iota,
			pkg:        visitor.pkg,
			visitor:    visitor,
		}

		if valueSpec.Values != nil {
			constant.expression = valueSpec.Values[0]
			constant.initType = valueSpec.Type
			//constant.value = visitor.evalConstantExpression(valueSpec.Values[0], params)
		} else {
			constant.expression = lastValueSpec.Values[0]
			constant.initType = lastValueSpec.Type
			//constant.value = visitor.evalConstantExpression(lastValueSpec.Values[0], params)
		}

		constants = append(constants, constant)
	}

	visitor.currentFile.constants.elements = append(visitor.currentFile.constants.elements, constants...)

	return constants
}

func (visitor *PackageVisitor) collectFunction() {
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

func (visitor *PackageVisitor) getTypeFromScope(name string) T {
	typ := visitor.pkg.Types.Scope().Lookup(name)

	typedName, ok := typ.Type().(*types.Named)

	if _, ok := visitor.collector.unprocessedTypes[visitor.pkg.ID]; !ok {
		visitor.collector.unprocessedTypes[visitor.pkg.ID] = make(map[string]T)
	}

	if ok {
		switch typedName.Underlying().(type) {
		case *types.Struct:
			structType := &Struct{
				name:        name,
				isProcessed: false,
			}
			visitor.collector.unprocessedTypes[visitor.pkg.ID][name] = structType
			return structType
		case *types.Interface:
			interfaceType := &Interface{
				name:        name,
				isProcessed: false,
			}
			visitor.collector.unprocessedTypes[visitor.pkg.ID][name] = interfaceType
			return interfaceType
		default:
			customType := &CustomType{
				name:        name,
				isProcessed: false,
			}
			visitor.collector.unprocessedTypes[visitor.pkg.ID][name] = customType
			return customType
		}
	}

	return nil
}

func (visitor *PackageVisitor) getTypeFromTypeSpec() T {
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

func (visitor *PackageVisitor) getInterface(specType *ast.TypeSpec) *Interface {
	interfaceType := &Interface{
		name:          specType.Name.Name,
		isExported:    ast.IsExported(specType.Name.Name),
		methods:       make([]*Function, 0),
		embeddeds:     make([]T, 0),
		position:      visitor.getPosition(specType.Pos()),
		markers:       visitor.packageMarkers[specType],
		file:          visitor.currentFile,
		isProcessed:   true,
		specType:      specType,
		visitor:       visitor,
		interfaceType: visitor.pkg.Types.Scope().Lookup(specType.Name.Name).Type().Underlying().(*types.Interface),
	}

	return interfaceType
}

func (visitor *PackageVisitor) processStruct(structType *Struct, specType *ast.TypeSpec) {
	structType.isExported = ast.IsExported(specType.Name.Name)
	structType.position = visitor.getPosition(specType.Pos())
	structType.markers = visitor.packageMarkers[specType]
	structType.file = visitor.currentFile
	structType.isProcessed = true
	structType.namedType = visitor.pkg.Types.Scope().Lookup(specType.Name.Name).Type().(*types.Named)

	fieldList := specType.Type.(*ast.StructType).Fields.List
	structType.fields = append(structType.fields, visitor.getFieldsFromFieldList(fieldList)...)
}

func (visitor *PackageVisitor) processInterface(interfaceType *Interface, specType *ast.TypeSpec) {
	interfaceType.isExported = ast.IsExported(specType.Name.Name)
	interfaceType.methods = visitor.getInterfaceMethods(specType.Type.(*ast.InterfaceType).Methods.List)
	interfaceType.embeddeds = visitor.getInterfaceEmbeddedTypes(specType.Type.(*ast.InterfaceType).Methods.List)
	interfaceType.position = visitor.getPosition(specType.Pos())
	interfaceType.markers = visitor.packageMarkers[specType]
	interfaceType.file = visitor.currentFile
	interfaceType.isProcessed = true
	interfaceType.interfaceType = visitor.pkg.Types.Scope().Lookup(specType.Name.Name).Type().Underlying().(*types.Interface)
}

func (visitor *PackageVisitor) processCustomType(customType *CustomType, specType *ast.TypeSpec) {
	customType.aliasType = visitor.getTypeFromExpression(specType.Type)
	customType.isExported = ast.IsExported(specType.Name.Name)
	customType.position = visitor.getPosition(specType.Pos())
	customType.file = visitor.currentFile
	customType.isProcessed = true
}

func (visitor *PackageVisitor) getStruct(specType *ast.TypeSpec) *Struct {

	structType := &Struct{
		name:        specType.Name.Name,
		isExported:  ast.IsExported(specType.Name.Name),
		position:    visitor.getPosition(specType.Pos()),
		markers:     visitor.packageMarkers[specType],
		file:        visitor.currentFile,
		fields:      make([]*Field, 0),
		allFields:   make([]*Field, 0),
		methods:     make([]*Function, 0),
		isProcessed: true,
		visitor:     visitor,
		specType:    specType,
		namedType:   visitor.pkg.Types.Scope().Lookup(specType.Name.Name).Type().(*types.Named),
	}

	return structType
}

func (visitor *PackageVisitor) getCustomType(specType *ast.TypeSpec) *CustomType {

	customType := &CustomType{
		name:        specType.Name.Name,
		aliasType:   visitor.getTypeFromExpression(specType.Type),
		isExported:  ast.IsExported(specType.Name.Name),
		position:    visitor.getPosition(specType.Pos()),
		markers:     visitor.packageMarkers[specType],
		methods:     make([]*Function, 0),
		file:        visitor.currentFile,
		isProcessed: true,
	}

	return customType
}

func (visitor *PackageVisitor) getTypeFromExpression(expr ast.Expr) T {

	switch typed := expr.(type) {
	case *ast.Ident:
		var typ T
		var ok bool
		typ, ok = basicTypesMap[typed.Name]

		if ok {
			return typ
		}

		if typed.Name == "error" {
			errorType, exists := visitor.collector.findTypeByPkgIdAndName("builtin", "error")

			if !exists {
				visitor.loadPackageAndVisit("builtin")
			}

			errorType, _ = visitor.collector.findTypeByPkgIdAndName("builtin", "error")

			return errorType

		}

		typ, ok = visitor.collector.findTypeByPkgIdAndName(visitor.pkg.ID, typed.Name)

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
			base: visitor.getTypeFromExpression(typed.X),
		}
	case *ast.ArrayType:

		if typed.Len == nil {
			return &Slice{
				elem: visitor.getTypeFromExpression(typed.Elt),
			}
		} else {
			basicLit, isBasicLit := typed.Len.(*ast.BasicLit)

			if isBasicLit {
				length, _ := strconv.ParseInt(basicLit.Value, 10, 64)
				return &Array{
					elem: visitor.getTypeFromExpression(typed.Elt),
					len:  length,
				}
			}

			return &Array{
				elem: visitor.getTypeFromExpression(typed.Elt),
				len:  -1,
			}
		}
	case *ast.ChanType:
		chanType := &Chan{
			elem: visitor.getTypeFromExpression(typed.Value),
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
			key:  visitor.getTypeFromExpression(typed.Key),
			elem: visitor.getTypeFromExpression(typed.Value),
		}
	case *ast.InterfaceType:
		interfaceType := &Interface{
			position:    visitor.getPosition(typed.Pos()),
			visitor:     visitor,
			isAnonymous: true,
			isProcessed: true,
		}
		interfaceType.methods = append(interfaceType.methods, visitor.getInterfaceMethods(typed.Methods.List)...)
		return interfaceType
	case *ast.StructType:
		structType := &Struct{
			position:    visitor.getPosition(typed.Pos()),
			isAnonymous: true,
			isProcessed: true,
		}

		return structType
	}

	return nil
}

func (visitor *PackageVisitor) loadPackageAndVisit(pkgId string) {
	if visitor.collector.isVisited(pkgId) {
		return
	}

	loadResult, err := LoadPackages(pkgId)

	if err != nil {
		panic(err)
	}

	pkg, _ := loadResult.Lookup(pkgId)

	visitPackage(pkg, visitor.collector, visitor.allPackageMarkers)
}

func (visitor *PackageVisitor) findTypeByImportAndTypeName(importName, typeName string) *ImportedType {
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

	visitor.loadPackageAndVisit(packageImport.path)

	typ, exists = visitor.collector.findTypeByPkgIdAndName(packageImport.path, typeName)

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

func (visitor *PackageVisitor) getVariables(fieldList []*ast.Field) *Tuple {
	tuple := &Tuple{
		variables: make([]*Variable, 0),
	}

	for _, field := range fieldList {

		typ := visitor.getTypeFromExpression(field.Type)

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

func (visitor *PackageVisitor) getInterfaceMethods(fieldList []*ast.Field) []*Function {
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

func (visitor *PackageVisitor) getInterfaceEmbeddedTypes(fieldList []*ast.Field) []T {
	embeddedTypes := make([]T, 0)

	for _, field := range fieldList {
		_, ok := field.Type.(*ast.FuncType)

		if !ok {
			embeddedTypes = append(embeddedTypes, visitor.getTypeFromExpression(field.Type))
		}
	}

	return embeddedTypes
}

func (visitor *PackageVisitor) getFieldsFromFieldList(fieldList []*ast.Field) []*Field {
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

func VisitPackages(pkgList []*Package, allPackageMarkers map[string]map[ast.Node]MarkerValues) {
	pkgCollector := newPackageCollector()

	for _, pkg := range pkgList {
		if !pkgCollector.isVisited(pkg.ID) || !pkgCollector.isProcessed(pkg.ID) {
			visitPackage(pkg, pkgCollector, allPackageMarkers)
		}
	}

	if pkgCollector == nil {

	}

	file, _ := pkgCollector.files["github.com/procyon-projects/marker/test/package2"].FindByName("test.go")

	//x, _ := file.Interfaces().FindByName("IFace")
	x, _ := file.Structs().FindByName("ComplexRequest")
	fields := x.AllFields()

	if fields == nil {

	}

	i, _ := file.Interfaces().FindByName("IFace")
	methods := i.ExplicitMethods()

	if methods == nil {

	}

	log.Printf("")
}

func visitPackage(pkg *Package, collector *packageCollector, allPackageMarkers map[string]map[ast.Node]MarkerValues) {
	pkgVisitor := &PackageVisitor{
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

func IsInterfaceType(t T) bool {
	_, ok := t.(*Interface)
	return ok
}

func IsStructType(t T) bool {
	_, ok := t.(*Struct)
	return ok
}

func IsErrorType(t T) bool {
	interfaceType, ok := t.(*Interface)
	if !ok {
		return false
	}

	if interfaceType.file == nil || interfaceType.file.pkg == nil {
		return false
	}

	return interfaceType.file.pkg.ID == "builtin"
}
