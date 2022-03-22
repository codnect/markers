package visitor

import "strings"

type Import struct {
	name       string
	path       string
	sideEffect bool
	position   Position
}

func (i *Import) Name() string {
	return i.name
}

func (i *Import) Path() string {
	return i.path
}

func (i *Import) SideEffect() bool {
	return i.sideEffect
}

func (i *Import) Position() Position {
	return i.position
}

type Imports struct {
	imports []*Import
}

func (i *Imports) Len() int {
	return len(i.imports)
}

func (i *Imports) At(index int) *Import {
	if index >= 0 && index < len(i.imports) {
		return i.imports[index]
	}

	return nil
}

func (i *Imports) FindByName(name string) (*Import, bool) {
	for _, importItem := range i.imports {
		if importItem.name == name || strings.HasSuffix(importItem.path, "/"+name) {
			return importItem, true
		}
	}

	return nil, false
}

func (i *Imports) FindByPath(path string) (*Import, bool) {
	for _, importItem := range i.imports {
		if importItem.path == path {
			return importItem, true
		}
	}

	return nil, false
}
