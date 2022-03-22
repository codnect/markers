package visitor

import "github.com/procyon-projects/marker"

type CustomType struct {
	name       string
	aliasType  Type
	isExported bool
	position   Position
	markers    marker.MarkerValues
	methods    []*Function
	file       *File

	isProcessed bool
}

func (c *CustomType) Name() string {
	return c.name
}

func (c *CustomType) AliasType() Type {
	return c.aliasType
}

func (c *CustomType) Underlying() Type {
	return c
}

func (c *CustomType) String() string {
	return ""
}

type CustomTypes struct {
	elements []*CustomType
}

func (c *CustomTypes) Len() int {
	return len(c.elements)
}

func (c *CustomTypes) At(index int) *CustomType {
	if index >= 0 && index < len(c.elements) {
		return c.elements[index]
	}

	return nil
}

func (c *CustomTypes) FindByName(name string) (*CustomType, bool) {
	for _, customType := range c.elements {
		if customType.name == name {
			return customType, true
		}
	}

	return nil, false
}
