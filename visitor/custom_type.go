package visitor

import (
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
)

type CustomType struct {
	name       string
	aliasType  Type
	isExported bool
	position   Position
	markers    markers.MarkerValues
	methods    []*Function
	file       *File

	isProcessed bool
	visitor     *packageVisitor
}

func newCustomType(specType *ast.TypeSpec, file *File, pkg *packages.Package, visitor *packageVisitor, markers markers.MarkerValues) *CustomType {
	customType := &CustomType{
		name:        specType.Name.Name,
		aliasType:   getTypeFromExpression(specType.Type, file, visitor),
		isExported:  ast.IsExported(specType.Name.Name),
		position:    getPosition(file.Package(), specType.Pos()),
		markers:     markers,
		methods:     make([]*Function, 0),
		file:        file,
		isProcessed: true,
		visitor:     visitor,
	}

	return customType.initialize()
}

func (c *CustomType) initialize() *CustomType {
	c.file.customTypes.elements = append(c.file.customTypes.elements, c)
	return c
}

func (c *CustomType) Name() string {
	return c.name
}

func (c *CustomType) IsExported() bool {
	return c.isExported
}

func (c *CustomType) AliasType() Type {
	return c.aliasType
}

func (c *CustomType) Underlying() Type {
	return c
}

func (c *CustomType) String() string {
	return fmt.Sprintf("type %s %s", c.name, c.aliasType.Name())
}

type CustomTypes struct {
	elements []*CustomType
}

func (c *CustomTypes) ToSlice() []*CustomType {
	return c.elements
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
