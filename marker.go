package marker

import (
	"errors"
	"go/ast"
	"path/filepath"
	"strings"
)

type TargetLevel int

const (
	PackageLevel TargetLevel = 1 << iota
	TypeLevel
	StructTypeLevel
	InterfaceTypeLevel
	FieldLevel
	FunctionLevel
	MethodLevel
	StructMethodLevel
	InterfaceMethodLevel
)

type MarkerValues map[string][]interface{}

func (markerValues MarkerValues) Get(name string) interface{} {
	result := markerValues[name]

	if len(result) == 0 {
		return nil
	}

	return result[0]
}

type markerComment struct {
	*ast.Comment
}

func newMarkerComment(comment *ast.Comment) markerComment {
	return markerComment{
		comment,
	}
}

func (comment *markerComment) Text() string {
	return strings.TrimSpace(comment.Comment.Text[2:])
}

func splitMarker(marker string) (name string, anonymousName string, options string) {
	marker = marker[1:]

	nameFieldParts := strings.SplitN(marker, "=", 2)

	if len(nameFieldParts) == 1 {
		return nameFieldParts[0], nameFieldParts[0], ""
	}

	anonymousName = nameFieldParts[0]
	name = anonymousName

	nameParts := strings.Split(name, ":")

	if len(nameParts) > 1 {
		name = strings.Join(nameParts[:len(nameParts)-1], ":")
	}

	return name, anonymousName, nameFieldParts[1]
}

func isMarkerComment(comment string) bool {
	if comment[0:2] != "//" {
		return false
	}

	stripped := strings.TrimSpace(comment[2:])

	if len(stripped) < 1 || stripped[0] != '+' {
		return false
	}

	return true
}

type Callback func(element *File, error error)

type Position struct {
	Line   int
	Column int
}

type Import struct {
	Name          string
	Path          string
	Position      Position
	RawImportSpec *ast.ImportSpec
}

type File struct {
	Name        string
	FullPath    string
	PackageName string
	Imports     []Import
	Markers     MarkerValues

	Functions      []Function
	StructTypes    []StructType
	InterfaceTypes []InterfaceType
	RawFile        *ast.File
}

type Field struct {
	Name     string
	Position Position
	Markers  MarkerValues
	Type     Type
	RawFile  *ast.File
	RawField *ast.Field
}

type Type struct {
	PackageName string
	Name        string
	IsPointer   bool
	RawObject   *ast.Object
}

type TypeInfo struct {
	Names    []string
	Type     Type
	RawField *ast.Field
}

type Method struct {
	Name         string
	Position     Position
	Markers      MarkerValues
	Receiver     *TypeInfo
	Parameters   []TypeInfo
	ReturnValues []TypeInfo
	File         *File
	RawFile      *ast.File
	RawField     *ast.Field
	RawFuncDecl  *ast.FuncDecl
	RawFuncType  *ast.FuncType
}

type StructType struct {
	Name        string
	Position    Position
	Markers     MarkerValues
	Fields      []Field
	Methods     []Method
	File        *File
	RawFile     *ast.File
	RawGenDecl  *ast.GenDecl
	RawTypeSpec *ast.TypeSpec
}

type Function struct {
	Name         string
	Position     Position
	Markers      MarkerValues
	Parameters   []TypeInfo
	ReturnValues []TypeInfo
	File         *File
	RawFile      *ast.File
	RawFuncDecl  *ast.FuncDecl
	RawFuncType  *ast.FuncType
}

type InterfaceType struct {
	Name        string
	Position    Position
	Markers     MarkerValues
	Methods     []Method
	File        *File
	RawFile     *ast.File
	RawGenDecl  *ast.GenDecl
	RawTypeSpec *ast.TypeSpec
}

func EachFile(collector *Collector, pkg *Package, callback Callback) {

	if collector == nil {
		callback(nil, errors.New("collector cannot be nil"))
		return
	}

	if pkg == nil {
		callback(nil, errors.New("pkg(package) cannot be nil"))
		return
	}

	markers, err := collector.Collect(pkg)

	if err != nil {
		callback(nil, err)
		return
	}

	var fileInfoMap = make(map[*ast.File]*File)

	visitFiles(pkg, func(file *ast.File) {
		fileInfo, ok := fileInfoMap[file]

		if ok {
			return
		}

		position := pkg.Fset.Position(file.Pos())

		fileFullPath := position.Filename

		imports := make([]Import, 0)

		for _, importInfo := range file.Imports {
			importPosition := pkg.Fset.Position(importInfo.Pos())
			importName := ""

			if importInfo.Name != nil {
				importName = importInfo.Name.Name
			}

			imports = append(imports, Import{
				Name: importName,
				Path: importInfo.Path.Value,
				Position: Position{
					importPosition.Line,
					importPosition.Column,
				},
				RawImportSpec: importInfo,
			})
		}

		fileInfo = &File{
			Name:           filepath.Base(fileFullPath),
			FullPath:       fileFullPath,
			PackageName:    file.Name.Name,
			Imports:        imports,
			Markers:        markers[file],
			Functions:      make([]Function, 0),
			StructTypes:    make([]StructType, 0),
			InterfaceTypes: make([]InterfaceType, 0),
			RawFile:        file,
		}

		fileInfoMap[file] = fileInfo
	})

	visitTypeElements(pkg, func(file *ast.File, decl *ast.GenDecl, spec *ast.TypeSpec) {

		fileInfo, ok := fileInfoMap[file]

		if !ok {
			return
		}

		position := pkg.Fset.Position(spec.Pos())
		typePosition := Position{
			Line:   position.Line,
			Column: position.Column,
		}

		switch specType := spec.Type.(type) {
		case *ast.InterfaceType:
			interfaceType := InterfaceType{
				Name:        spec.Name.Name,
				Position:    typePosition,
				Markers:     markers[spec],
				File:        fileInfo,
				RawFile:     file,
				RawGenDecl:  decl,
				RawTypeSpec: spec,
			}

			for _, methodInfo := range specType.Methods.List {

				position := pkg.Fset.Position(methodInfo.Pos())
				methodPosition := Position{
					Line:   position.Line,
					Column: position.Column,
				}

				method := &Method{
					Name:        methodInfo.Names[0].Name,
					Position:    methodPosition,
					Markers:     markers[methodInfo],
					File:        fileInfo,
					RawFile:     file,
					RawFuncType: methodInfo.Type.(*ast.FuncType),
					RawField:    methodInfo,
				}

				if methodInfo.Type.(*ast.FuncType).Params != nil {
					method.Parameters = getTypesInfo(methodInfo.Type.(*ast.FuncType).Params.List)
				}

				if methodInfo.Type.(*ast.FuncType).Results != nil {
					method.ReturnValues = getTypesInfo(methodInfo.Type.(*ast.FuncType).Results.List)
				}

				interfaceType.Methods = append(interfaceType.Methods, *method)
			}

			fileInfo.InterfaceTypes = append(fileInfo.InterfaceTypes, interfaceType)
		case *ast.StructType:
			structType := StructType{
				Name:        spec.Name.Name,
				Position:    typePosition,
				Markers:     markers[spec],
				File:        fileInfo,
				RawFile:     file,
				RawGenDecl:  decl,
				RawTypeSpec: spec,
			}

			fieldTypeInfoList := getTypesInfo(specType.Fields.List)

			for _, fieldTypeInfo := range fieldTypeInfoList {

				position := pkg.Fset.Position(fieldTypeInfo.RawField.Pos())
				fieldPosition := Position{
					Line:   position.Line,
					Column: position.Column,
				}

				for _, fieldName := range fieldTypeInfo.Names {
					field := &Field{
						Name:     fieldName,
						Position: fieldPosition,
						Markers:  markers[fieldTypeInfo.RawField],
						Type:     fieldTypeInfo.Type,
						RawFile:  file,
						RawField: fieldTypeInfo.RawField,
					}

					structType.Fields = append(structType.Fields, *field)
				}

			}

			fileInfo.StructTypes = append(fileInfo.StructTypes, structType)
		}

	})

	visitFunctions(pkg, func(file *ast.File, decl *ast.FuncDecl, funcType *ast.FuncType) {

		fileInfo, ok := fileInfoMap[file]

		if !ok {
			return
		}

		position := pkg.Fset.Position(funcType.Pos())
		functionPosition := Position{
			Line:   position.Line,
			Column: position.Column,
		}

		// If Recv is nil, it is a function, not a method
		if decl.Recv == nil {

			function := &Function{
				Name:        decl.Name.Name,
				Position:    functionPosition,
				Markers:     markers[decl],
				File:        fileInfo,
				RawFile:     file,
				RawFuncDecl: decl,
				RawFuncType: funcType,
			}

			if funcType.Params != nil {
				function.Parameters = getTypesInfo(funcType.Params.List)
			}

			if funcType.Results != nil {
				function.ReturnValues = getTypesInfo(funcType.Results.List)
			}

			fileInfo.Functions = append(fileInfo.Functions, *function)
		} else {
			method := &Method{
				Name:     decl.Name.Name,
				Position: functionPosition,
				Markers:  markers[decl],
				Receiver: &TypeInfo{
					Names: make([]string, 0),
					Type:  Type{},
				},
				File:        fileInfo,
				RawFile:     file,
				RawFuncDecl: decl,
				RawFuncType: funcType,
			}

			if funcType.Params != nil {
				method.Parameters = getTypesInfo(funcType.Params.List)
			}

			if funcType.Results != nil {
				method.ReturnValues = getTypesInfo(funcType.Results.List)
			}

			receiver := decl.Recv.List[0]
			receiverType := getType(receiver)

			method.Receiver.Type = receiverType
			method.Receiver.Names = append(method.Receiver.Names, receiver.Names[0].Name)

			//  Find the struct type to add the method into its list.
			for _, fileInfo := range fileInfoMap {

				for typeIndex, structType := range fileInfo.StructTypes {

					// if RawObject is nil, try to resolve the struct type by receiver type name
					if receiverType.RawObject == nil {
						if file.Name.Name != fileInfo.PackageName && structType.Name != receiverType.Name {
							continue
						}
					} else if structType.RawTypeSpec != receiverType.RawObject.Decl.(*ast.TypeSpec) {
						continue
					}

					fileInfo.StructTypes[typeIndex].Methods = append(fileInfo.StructTypes[typeIndex].Methods, *method)

				}
			}

		}

	})

	for _, fileInfo := range fileInfoMap {
		callback(fileInfo, nil)
	}
}
