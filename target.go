package markers

import "go/ast"

// TargetLevel describes which kind of nodes a given marker are associated with.
type TargetLevel int

const (
	InvalidLevel TargetLevel = 1 << iota
	// PackageLevel indicates that a marker is associated with a package.
	PackageLevel
	// StructTypeLevel indicates that a marker is associated with a struct type.
	StructTypeLevel
	// InterfaceTypeLevel indicates that a marker is associated with an interface type.
	InterfaceTypeLevel
	// FieldLevel indicates that a marker is associated with a struct field.
	FieldLevel
	// FunctionLevel indicates that a marker is associated with a function.
	FunctionLevel
	// StructMethodLevel indicates that a marker is associated with a struct method.
	StructMethodLevel
	// InterfaceMethodLevel indicates that a marker is associated with an interface method.
	InterfaceMethodLevel
)

// Combined levels
const (
	// TypeLevel indicates that a marker is associated with any type.
	TypeLevel = StructTypeLevel | InterfaceTypeLevel
	// MethodLevel indicates that a marker is associated with a struct method or an interface method.
	MethodLevel = StructMethodLevel | InterfaceMethodLevel
	AllLevels   = PackageLevel | TypeLevel | MethodLevel | FieldLevel | FunctionLevel
)

func FindTargetLevelFromNode(node ast.Node) TargetLevel {
	switch typedNode := node.(type) {
	case *ast.TypeSpec:
		_, isStructType := typedNode.Type.(*ast.StructType)
		if isStructType {
			return StructTypeLevel
		}

		_, isInterfaceType := typedNode.Type.(*ast.InterfaceType)
		if isInterfaceType {
			return InterfaceTypeLevel
		}
	case *ast.Field:
		_, isFuncType := typedNode.Type.(*ast.FuncType)
		if !isFuncType {
			return FieldLevel
		} else if isFuncType {
			return InterfaceMethodLevel
		}
	case *ast.FuncDecl:
		if typedNode.Recv != nil {
			return StructMethodLevel
		} else {
			return FunctionLevel
		}
	case *ast.Package:
		return PackageLevel
	}

	return InvalidLevel
}
