package marker

import (
	"go/ast"
	"strings"
)

// TargetLevel describes which kind of node a given marker is associated with.
type TargetLevel int

const (
	// PackageLevel indicates that a marker is associated with a package.
	PackageLevel TargetLevel = 1 << iota
	// TypeLevel indicates that a marker is associated with any type.
	TypeLevel
	// StructTypeLevel indicates that a marker is associated with a struct type.
	StructTypeLevel
	// InterfaceTypeLevel indicates that a marker is associated with an interface type.
	InterfaceTypeLevel
	// FieldLevel indicates that a marker is associated with a struct field.
	FieldLevel
	// FunctionLevel indicates that a marker is associated with a function.
	FunctionLevel
	// MethodLevel indicates that a marker is associated with a struct method or an interface method.
	MethodLevel
	// StructMethodLevel indicates that a marker is associated with a struct method.
	StructMethodLevel
	// InterfaceMethodLevel indicates that a marker is associated with an interface method.
	InterfaceMethodLevel
)

type MarkerValues map[string][]interface{}

func (markerValues MarkerValues) Get(name string) interface{} {
	result := markerValues[name]

	if len(result) == 0 {
		return nil
	}

	return result[0]
}

type markerComment struct {
	*ast.Comment
}

func newMarkerComment(comment *ast.Comment) markerComment {
	return markerComment{
		comment,
	}
}

func (comment *markerComment) Text() string {
	return strings.TrimSpace(comment.Comment.Text[2:])
}

func splitMarker(marker string) (name string, anonymousName string, options string) {
	marker = marker[1:]

	nameFieldParts := strings.SplitN(marker, "=", 2)

	if len(nameFieldParts) == 1 {
		return nameFieldParts[0], nameFieldParts[0], ""
	}

	anonymousName = nameFieldParts[0]
	name = anonymousName

	nameParts := strings.Split(name, ":")

	if len(nameParts) > 1 {
		name = strings.Join(nameParts[:len(nameParts)-1], ":")
	}

	return name, anonymousName, nameFieldParts[1]
}

func isMarkerComment(comment string) bool {
	if comment[0:2] != "//" {
		return false
	}

	stripped := strings.TrimSpace(comment[2:])

	if len(stripped) < 1 || stripped[0] != '+' {
		return false
	}

	return true
}
