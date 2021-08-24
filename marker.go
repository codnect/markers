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

type Callback func(element *File)

type File struct {
	Name        string
	FullPath    string
	PackageName string
	Markers     MarkerValues

	Functions      []Function
	StructTypes    []StructType
	InterfaceTypes []InterfaceType
	RawFile        *ast.File
}

type Field struct {
	Name     string
	Markers  MarkerValues
	Type     Type
	RawFile  *ast.File
	RawField *ast.Field
}

type Type struct {
	ImportAlias string
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
	Markers     MarkerValues
	Methods     []Method
	File        *File
	RawFile     *ast.File
	RawGenDecl  *ast.GenDecl
	RawTypeSpec *ast.TypeSpec
}

func EachFile(collector *Collector, pkg *Package, callback Callback) error {

	if collector == nil {
		return errors.New("collector cannot be nil")
	}

	if pkg == nil {
		return errors.New("pkg(package) cannot be nil")
	}

	markers, err := collector.Collect(pkg)

	if err != nil {
		return err
	}

	var fileInfoMap = make(map[*ast.File]*File)

	visitFiles(pkg, func(file *ast.File) {
		fileInfo, ok := fileInfoMap[file]

		if ok {
			return
		}

		position := pkg.Fset.Position(file.Pos())
		fileFullPath := position.Filename

		fileInfo = &File{
			Name:           filepath.Base(fileFullPath),
			FullPath:       fileFullPath,
			PackageName:    file.Name.Name,
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

		switch specType := spec.Type.(type) {
		case *ast.InterfaceType:
			interfaceType := InterfaceType{
				Name:        spec.Name.Name,
				Markers:     markers[spec],
				File:        fileInfo,
				RawFile:     file,
				RawGenDecl:  decl,
				RawTypeSpec: spec,
			}

			for _, methodInfo := range specType.Methods.List {
				method := &Method{
					Name:        methodInfo.Names[0].Name,
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
				Markers:     markers[spec],
				File:        fileInfo,
				RawFile:     file,
				RawGenDecl:  decl,
				RawTypeSpec: spec,
			}

			fieldTypeInfoList := getTypesInfo(specType.Fields.List)

			for _, fieldTypeInfo := range fieldTypeInfoList {

				for _, fieldName := range fieldTypeInfo.Names {
					field := &Field{
						Name:     fieldName,
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

		// If Recv is nil, it is a function, not a method
		if decl.Recv == nil {

			function := &Function{
				Name:        decl.Name.Name,
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
				Name:    decl.Name.Name,
				Markers: markers[decl],
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
		callback(fileInfo)
	}

	return nil
}
