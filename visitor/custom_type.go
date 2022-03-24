package visitor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
)

type CustomType struct {
	name       string
	aliasType  Type
	isExported bool
	position   Position
	markers    marker.MarkerValues
	methods    []*Function
	file       *File

	isProcessed bool
	visitor     *packageVisitor
}

func newCustomType(specType *ast.TypeSpec, file *File, pkg *packages.Package, visitor *packageVisitor, markers marker.MarkerValues) *CustomType {
	customType := &CustomType{
		name:        specType.Name.Name,
		aliasType:   getTypeFromExpression(specType.Type, visitor),
		isExported:  ast.IsExported(specType.Name.Name),
		position:    getPosition(file.Package(), specType.Pos()),
		markers:     markers,
		methods:     make([]*Function, 0),
		file:        file,
		isProcessed: true,
		visitor:     visitor,
	}

	return customType
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
