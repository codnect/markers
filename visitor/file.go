package visitor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"path/filepath"
)

type File struct {
	name string
	path string
	pkg  *packages.Package

	allMarkers  markers.MarkerValues
	fileMarkers markers.MarkerValues

	imports       *Imports
	importMarkers []markers.ImportMarker

	functions   *Functions
	structs     *Structs
	interfaces  *Interfaces
	customTypes *CustomTypes
	constants   *Constants

	rawFile *ast.File

	visitor *packageVisitor
}

func newFile(rawFile *ast.File, pkg *packages.Package, markerValues markers.MarkerValues, visitor *packageVisitor) *File {
	position := pkg.Fset.Position(rawFile.Pos())
	path := position.Filename

	file := &File{
		name:          filepath.Base(path),
		path:          path,
		pkg:           pkg,
		allMarkers:    markerValues,
		fileMarkers:   make(markers.MarkerValues, 0),
		imports:       &Imports{},
		importMarkers: make([]markers.ImportMarker, 0),
		functions:     &Functions{},
		structs:       &Structs{},
		interfaces:    &Interfaces{},
		customTypes:   &CustomTypes{},
		constants:     &Constants{},
		rawFile:       rawFile,
		visitor:       visitor,
	}

	return file.initialize()
}

func (f *File) initialize() *File {
	for markerName, markerValues := range f.allMarkers {
		if markers.ImportMarkerName == markerName {
			for _, importMarker := range markerValues {
				f.importMarkers = append(f.importMarkers, importMarker.(markers.ImportMarker))
			}
		} else {
			f.fileMarkers[markerName] = append(f.fileMarkers[markerName], markerValues...)
		}
	}

	for _, importPackage := range f.rawFile.Imports {
		importPosition := getPosition(f.pkg, importPackage.Pos())
		importName := ""

		if importPackage.Name != nil {
			importName = importPackage.Name.Name
		}

		sideEffect := false

		if importName == "_" {
			sideEffect = true
		}

		f.imports.elements = append(f.imports.elements, &Import{
			name:       importName,
			path:       importPackage.Path.Value[1 : len(importPackage.Path.Value)-1],
			sideEffect: sideEffect,
			position: Position{
				importPosition.Line,
				importPosition.Column,
			},
		})
	}

	return f
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Path() string {
	return f.path
}

func (f *File) Markers() markers.MarkerValues {
	return f.fileMarkers
}

func (f *File) Package() *packages.Package {
	return f.pkg
}

func (f *File) Imports() *Imports {
	return f.imports
}

func (f *File) ImportMarkers() []markers.ImportMarker {
	return f.importMarkers
}

func (f *File) NumImportMarkers() int {
	return len(f.importMarkers)
}

func (f *File) Constants() *Constants {
	return f.constants
}

func (f *File) Functions() *Functions {
	return f.functions
}

func (f *File) Structs() *Structs {
	return f.structs
}

func (f *File) Interfaces() *Interfaces {
	return f.interfaces
}

func (f *File) CustomTypes() *CustomTypes {
	return f.customTypes
}

type Files struct {
	elements []*File
}

func (f *Files) FindByName(name string) (*File, bool) {
	for _, file := range f.elements {
		if file.name == name {
			return file, true
		}
	}

	return nil, false
}

func (f *Files) Len() int {
	return len(f.elements)
}

func (f *Files) At(index int) *File {
	if index >= 0 && index < len(f.elements) {
		return f.elements[index]
	}

	return nil
}
