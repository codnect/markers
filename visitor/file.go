package visitor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
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
