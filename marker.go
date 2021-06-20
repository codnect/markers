package marker

import (
	"errors"
	"go/ast"
	"strings"
)

type TargetLevel int

const (
	PackageLevel TargetLevel = 1 << iota
	TypeLevel
	StructTypeLevel
	InterfaceTypeLevel
	FieldLevel
	FunctionLevel
	MethodLevel
	StructMethodLevel
	InterfaceMethodLevel
)

type Marker struct {
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

func Scan(collector *Collector, pkg *Package) error {

	if collector == nil {
		return errors.New("collector cannot be nil")
	}

	if pkg == nil {
		return errors.New("pkg(package) cannot be nil")
	}

	err := collector.Collect(pkg)

	if err != nil {
		return err
	}

	return nil
}
