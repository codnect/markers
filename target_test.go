package markers

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestFindTargetLevelFromNode(t *testing.T) {
	assert.Equal(t, StructTypeLevel, FindTargetLevelFromNode(&ast.TypeSpec{
		Type: &ast.StructType{},
	}))
	assert.Equal(t, InterfaceTypeLevel, FindTargetLevelFromNode(&ast.TypeSpec{
		Type: &ast.InterfaceType{},
	}))
	assert.Equal(t, FieldLevel, FindTargetLevelFromNode(&ast.Field{}))
	assert.Equal(t, InterfaceMethodLevel, FindTargetLevelFromNode(&ast.Field{
		Type: &ast.FuncType{},
	}))
	assert.Equal(t, StructMethodLevel, FindTargetLevelFromNode(&ast.FuncDecl{
		Recv: &ast.FieldList{},
	}))
	assert.Equal(t, FunctionLevel, FindTargetLevelFromNode(&ast.FuncDecl{}))
	assert.Equal(t, PackageLevel, FindTargetLevelFromNode(&ast.Package{}))
	assert.Equal(t, InvalidLevel, FindTargetLevelFromNode(nil))
}
