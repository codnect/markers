package visitor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"path/filepath"
)

type File struct {
	name     string
	fullPath string
	pkg      *packages.Package

	allMarkers  marker.MarkerValues
	fileMarkers marker.MarkerValues

	imports       *Imports
	importMarkers []marker.ImportMarker

	functions   *Functions
	structs     *Structs
	interfaces  *Interfaces
	customTypes *CustomTypes
	constants   *Constants

	rawFile *ast.File

	visitor *packageVisitor
}

func newFile(rawFile *ast.File, pkg *packages.Package, markers marker.MarkerValues, visitor *packageVisitor) *File {
	position := pkg.Fset.Position(rawFile.Pos())
	fileFullPath := position.Filename

	file := &File{
		name:          filepath.Base(fileFullPath),
		fullPath:      fileFullPath,
		pkg:           pkg,
		allMarkers:    markers,
		fileMarkers:   make(marker.MarkerValues, 0),
		imports:       &Imports{},
		importMarkers: make([]marker.ImportMarker, 0),
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
	for markerName, markers := range f.allMarkers {
		if marker.ImportMarkerName == markerName {
			for _, importMarker := range markers {
				f.importMarkers = append(f.importMarkers, importMarker.(marker.ImportMarker))
			}
		} else {
			f.fileMarkers[markerName] = append(f.fileMarkers[markerName], markers...)
		}
	}

	for _, importPackage := range f.rawFile.Imports {
		importPosition := getPosition(f.pkg, importPackage.Pos())
		importName := ""

		if importPackage.Name != nil {
			importName = importPackage.Name.Name
		}

		f.imports.elements = append(f.imports.elements, &Import{
			name: importName,
			path: importPackage.Path.Value[1 : len(importPackage.Path.Value)-1],
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

func (f *File) FullPath() string {
	return f.name
}

func (f *File) Markers() marker.MarkerValues {
	return f.fileMarkers
}

func (f *File) Package() *packages.Package {
	return f.pkg
}

func (f *File) Imports() *Imports {
	return f.imports
}

func (f *File) ImportMarkers() []marker.ImportMarker {
	return f.importMarkers
}

func (f *File) NumImportMarkers() int {
	return len(f.importMarkers)
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
