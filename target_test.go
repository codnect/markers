package markers

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestFindTargetLevelFromNode(t *testing.T) {
	assert.Equal(t, StructTypeLevel, FindTargetLevel(&ast.TypeSpec{
		Type: &ast.StructType{},
	}))
	assert.Equal(t, InterfaceTypeLevel, FindTargetLevel(&ast.TypeSpec{
		Type: &ast.InterfaceType{},
	}))
	assert.Equal(t, FieldLevel, FindTargetLevel(&ast.Field{}))
	assert.Equal(t, InterfaceMethodLevel, FindTargetLevel(&ast.Field{
		Type: &ast.FuncType{},
	}))
	assert.Equal(t, StructMethodLevel, FindTargetLevel(&ast.FuncDecl{
		Recv: &ast.FieldList{},
	}))
	assert.Equal(t, FunctionLevel, FindTargetLevel(&ast.FuncDecl{}))
	assert.Equal(t, PackageLevel, FindTargetLevel(&ast.Package{}))
	assert.Equal(t, InvalidLevel, FindTargetLevel(nil))
}
