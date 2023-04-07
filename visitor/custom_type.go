package visitor

import (
	"fmt"
	"github.com/procyon-projects/markers"
	"github.com/procyon-projects/markers/packages"
	"go/ast"
	"strings"
	"sync"
)

type CustomType struct {
	name           string
	underlyingType Type
	isExported     bool
	position       Position
	markers        markers.Values
	typeParams     *TypeParameters
	methods        []*Function
	file           *File

	pkg      *packages.Package
	specType *ast.TypeSpec

	isProcessed bool
	visitor     *packageVisitor

	typeParamsOnce sync.Once
}

func newCustomType(specType *ast.TypeSpec, file *File, pkg *packages.Package, visitor *packageVisitor, markers markers.Values) *CustomType {
	customType := &CustomType{
		markers: markers,
		methods: make([]*Function, 0),
		typeParams: &TypeParameters{
			[]*TypeParameter{},
		},
		isProcessed: true,
		file:        file,
		visitor:     visitor,
		pkg:         pkg,
		specType:    specType,
	}

	return customType.initialize(specType, file, pkg)
}

func (c *CustomType) initialize(specType *ast.TypeSpec, file *File, pkg *packages.Package) *CustomType {
	c.isProcessed = true
	c.specType = specType
	c.file = file
	c.pkg = pkg

	if specType != nil {
		c.name = specType.Name.Name
		c.isExported = ast.IsExported(specType.Name.Name)
		c.position = getPosition(file.Package(), specType.Pos())

		c.loadTypeParams()
		c.underlyingType = getTypeFromExpression(specType.Type, file, c.visitor, c, c.typeParams)
		if _, exists := file.customTypes.FindByName(c.name); !exists {
			c.file.customTypes.elements = append(c.file.customTypes.elements, c)
		}
	}

	return c
}

func (c *CustomType) Name() string {
	return c.name
}

func (c *CustomType) IsExported() bool {
	return c.isExported
}

func (c *CustomType) Underlying() Type {
	return c.underlyingType
}

func (c *CustomType) TypeParameters() *TypeParameters {
	c.loadTypeParams()
	return c.typeParams
}

func (c *CustomType) NumMethods() int {
	return len(c.methods)
}

func (c *CustomType) Methods() *Functions {
	return &Functions{
		elements: c.methods,
	}
}

func (c *CustomType) Markers() markers.Values {
	return c.markers
}

func (c *CustomType) SpecType() *ast.TypeSpec {
	return c.specType
}

func (c *CustomType) String() string {
	var builder strings.Builder

	if c.file != nil && c.file.pkg.Name != "builtin" {
		builder.WriteString(fmt.Sprintf("%s.%s", c.file.Package().Name, c.name))
	} else if c.name != "" {
		builder.WriteString(c.name)
	}

	if c.TypeParameters().Len() != 0 {
		builder.WriteString("[")

		for index := 0; index < c.TypeParameters().Len(); index++ {
			typeParam := c.TypeParameters().At(index)
			builder.WriteString(typeParam.String())

			if index != c.TypeParameters().Len()-1 {
				builder.WriteString(",")
			}

		}

		builder.WriteString("]")
	}

	return builder.String()
}

func (c *CustomType) loadTypeParams() {
	c.typeParamsOnce.Do(func() {
		if c.specType == nil || c.specType.TypeParams == nil {
			return
		}

		for _, field := range c.specType.TypeParams.List {
			for _, fieldName := range field.Names {
				typeParameter := &TypeParameter{
					name: fieldName.Name,
					constraints: &TypeConstraints{
						[]*TypeConstraint{},
					},
				}
				c.typeParams.elements = append(c.typeParams.elements, typeParameter)
			}
		}

		for _, field := range c.specType.TypeParams.List {
			constraints := make([]*TypeConstraint, 0)
			typ := getTypeFromExpression(field.Type, c.file, c.visitor, nil, c.typeParams)

			if typeSets, isTypeSets := typ.(TypeSets); isTypeSets {
				for _, item := range typeSets {
					if constraint, isConstraint := item.(*TypeConstraint); isConstraint {
						constraints = append(constraints, constraint)
					} else {
						constraints = append(constraints, &TypeConstraint{typ: item})
					}
				}
			} else {
				if constraint, isConstraint := typ.(*TypeConstraint); isConstraint {
					constraints = append(constraints, constraint)
				} else {
					constraints = append(constraints, &TypeConstraint{typ: typ})
				}
			}

			for _, fieldName := range field.Names {
				typeParam, exists := c.typeParams.FindByName(fieldName.Name)

				if exists {
					typeParam.constraints.elements = append(typeParam.constraints.elements, constraints...)
				}
			}
		}

	})
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
