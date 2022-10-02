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
	elements []*Import
}

func (i *Imports) ToSlice() []*Import {
	return i.elements
}

func (i *Imports) Len() int {
	return len(i.elements)
}

func (i *Imports) At(index int) *Import {
	if index >= 0 && index < len(i.elements) {
		return i.elements[index]
	}

	return nil
}

func (i *Imports) SideEffects() []*Import {
	result := make([]*Import, 0)
	for _, importItem := range i.elements {
		if importItem.SideEffect() {
			result = append(result, importItem)
		}
	}

	return result
}

func (i *Imports) FindByName(name string) (*Import, bool) {
	for _, importItem := range i.elements {
		if importItem.name == name || strings.HasSuffix(importItem.path, "/"+name) {
			return importItem, true
		}
	}

	return nil, false
}

func (i *Imports) FindByPath(path string) (*Import, bool) {
	for _, importItem := range i.elements {
		if importItem.path == path {
			return importItem, true
		}
	}

	return nil, false
}
