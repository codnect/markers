package visitor

import (
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"go/token"
	"strconv"
)

type Constant struct {
	name       string
	isExported bool
	value      any
	typ        Type
	expression ast.Expr
	initType   ast.Expr

	iota                int
	expressionEvaluated bool

	file    *File
	pkg     *packages.Package
	visitor *packageVisitor
}

func (c *Constant) Name() string {
	return c.name
}

func (c *Constant) Value() any {
	c.evaluateExpression()
	return c.value
}

func (c *Constant) evaluateExpression() {
	if c.expressionEvaluated {
		return
	}

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	params := make(map[string]any, 0)
	params["iota"] = c.iota
	c.value, c.typ = c.evalConstantExpression(c.expression, params)

	if c.initType != nil {
		switch typed := c.initType.(type) {
		case *ast.Ident:
			c.typ, _ = c.visitor.collector.findTypeByPkgIdAndName(c.pkg.ID, typed.Name)
		case *ast.SelectorExpr:
			c.typ = c.visitor.collector.findTypeByImportAndTypeName(typed.X.(*ast.Ident).Name, typed.Sel.Name, c.file)
		}
	}

	c.expressionEvaluated = true
}

func (c *Constant) Type() Type {
	return c.typ
}

func (c *Constant) IsExported() string {
	return c.name
}

func (c *Constant) Underlying() Type {
	return c
}

func (c *Constant) String() string {
	return ""
}

func (c *Constant) evalConstantExpression(exp ast.Expr, variableMap map[string]any) (any, Type) {
	switch exp := exp.(type) {
	case *ast.Ident:
		if value, ok := variableMap[exp.Name]; ok {
			return value, basicTypes[UntypedInt]
		}

		candidateConstant, ok := c.visitor.collector.findTypeByPkgIdAndName(c.pkg.ID, exp.Name)
		if ok {
			constant := candidateConstant.(*Constant)
			return constant.Value(), constant.Type()
		}

		return nil, nil
	case *ast.SelectorExpr:
		importedType := c.visitor.collector.findTypeByImportAndTypeName(exp.X.(*ast.Ident).Name, exp.Sel.Name, c.file)
		if importedType != nil {
			constant := importedType.Underlying().(*Constant)
			return constant.Value(), constant.Type()
		}
	case *ast.BinaryExpr:
		return c.evalBinaryExpr(exp, variableMap)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT:
			i, _ := strconv.Atoi(exp.Value)
			return i, basicTypes[UntypedInt]
		case token.FLOAT:
			f, _ := strconv.ParseFloat(exp.Value, 64)
			return f, basicTypes[UntypedFloat]
		case token.STRING:
			return exp.Value[1 : len(exp.Value)-1], basicTypes[String]
		}
	case *ast.UnaryExpr:
		result, typ := c.evalConstantExpression(exp.X, variableMap)

		switch result.(type) {
		case int:
			return -1 * result.(int), typ
		case float64:
			return -1.0 * result.(float64), typ
		}
	case *ast.ParenExpr:
		return c.evalConstantExpression(exp.X, variableMap)
	}

	return nil, nil
}

func (c *Constant) evalBinaryExpr(exp *ast.BinaryExpr, variableMap map[string]any) (any, Type) {
	var expressionType Type
	left, typLeft := c.evalConstantExpression(exp.X, variableMap)
	right, typRight := c.evalConstantExpression(exp.Y, variableMap)

	_, isTypeLeftBasic := typLeft.(*Basic)
	_, isTypeRightBasic := typRight.(*Basic)

	if isTypeLeftBasic && isTypeRightBasic {
		expressionType = typLeft
	} else if !isTypeLeftBasic {
		expressionType = typLeft
	} else if !isTypeRightBasic {
		expressionType = typRight
	}

	switch left.(type) {
	case int:
		switch exp.Op {
		case token.ADD:
			return left.(int) + right.(int), expressionType
		case token.SUB:
			return left.(int) - right.(int), expressionType
		case token.MUL:
			return left.(int) * right.(int), expressionType
		case token.QUO:
			return left.(int) / right.(int), expressionType
		case token.REM:
			return left.(int) % right.(int), expressionType
		case token.AND:
			return left.(int) & right.(int), expressionType
		case token.OR:
			return left.(int) | right.(int), expressionType
		case token.XOR:
			return left.(int) ^ right.(int), expressionType
		case token.SHL:
			return left.(int) << right.(int), expressionType
		case token.SHR:
			return left.(int) >> right.(int), expressionType
		case token.AND_NOT:
			return left.(int) &^ right.(int), expressionType
		case token.EQL:
			return left.(int) == right.(int), basicTypes[Bool]
		case token.NEQ:
			return left.(int) != right.(int), basicTypes[Bool]
		case token.LSS:
			return left.(int) < right.(int), basicTypes[Bool]
		case token.GTR:
			return left.(int) > right.(int), basicTypes[Bool]
		case token.LEQ:
			return left.(int) <= right.(int), basicTypes[Bool]
		case token.GEQ:
			return left.(int) >= right.(int), basicTypes[Bool]
		}
	case float64:
		switch exp.Op {
		case token.ADD:
			return left.(float64) + right.(float64), expressionType
		case token.SUB:
			return left.(float64) - right.(float64), expressionType
		case token.MUL:
			return left.(float64) * right.(float64), expressionType
		case token.QUO:
			return left.(float64) / right.(float64), expressionType
		}
	case string:
		switch exp.Op {
		case token.ADD:
			return left.(string) + right.(string), expressionType
		}
	}

	return nil, nil
}

type Constants struct {
	elements []*Constant
}

func (c *Constants) Len() int {
	return len(c.elements)
}

func (c *Constants) At(index int) *Constant {
	if index >= 0 && index < len(c.elements) {
		return c.elements[index]
	}

	return nil
}

func (c *Constants) FindByName(name string) (*Constant, bool) {
	for _, constant := range c.elements {
		if constant.name == name {
			return constant, true
		}
	}

	return nil, false
}

func collectConstantsFromSpecs(specs []ast.Spec, file *File) {
	var last *ast.ValueSpec
	for iota, s := range specs {
		valueSpec := s.(*ast.ValueSpec)

		switch {
		case valueSpec.Type != nil || len(valueSpec.Values) > 0:
			last = valueSpec
		case last == nil:
			last = new(ast.ValueSpec)
		}

		collectConstants(valueSpec, last, iota, file)
	}
}

func collectConstants(valueSpec *ast.ValueSpec, lastValueSpec *ast.ValueSpec, iota int, file *File) {
	for _, name := range valueSpec.Names {
		constant := &Constant{
			name:       name.Name,
			isExported: ast.IsExported(name.Name),
			iota:       iota,
			pkg:        file.pkg,
			//visitor:    visitor,
		}

		if valueSpec.Values != nil {
			constant.expression = valueSpec.Values[0]
			constant.initType = valueSpec.Type
		} else {
			constant.expression = lastValueSpec.Values[0]
			constant.initType = lastValueSpec.Type
		}

		file.constants.elements = append(file.constants.elements, constant)
	}
}
