package marker

import (
	"go/ast"
	"go/token"
	"strings"
)

// TargetLevel describes which kind of nodes a given marker are associated with.
type TargetLevel int

const (
	// PackageLevel indicates that a marker is associated with a package.
	PackageLevel TargetLevel = 1 << iota
	// ImportLevel indicates that a marker is associated with an import block.
	ImportLevel
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

const ValueArgument = "Value"

type ImportMarker struct {
	Value string `marker:"Value"`
	Name  string `marker:"Marker"`
	Pkg   string `marker:"Package"`
}

type MarkerValues map[string][]interface{}

func (markerValues MarkerValues) Get(name string) interface{} {
	result := markerValues[name]

	if len(result) == 0 {
		return nil
	}

	return result[0]
}

type markerComment struct {
	commentLines []*ast.Comment
}

func newMarkerComment(comment *ast.Comment) *markerComment {
	markerComment := &markerComment{
		make([]*ast.Comment, 0),
	}

	markerComment.commentLines = append(markerComment.commentLines, comment)

	return markerComment
}

func (c markerComment) Pos() token.Pos {
	return c.commentLines[0].Pos()
}

func (c markerComment) End() token.Pos {
	return c.commentLines[len(c.commentLines)].End()
}

func (c *markerComment) append(comment *ast.Comment) {
	c.commentLines = append(c.commentLines, comment)
}

func (c *markerComment) Text() string {
	var text string
	for _, line := range c.commentLines {
		comment := strings.TrimSpace(line.Text[2:])

		if strings.HasSuffix(comment, "\\") {
			comment = strings.TrimSpace(comment[:len(comment)-1])
		}

		if text == "" {
			text = comment
		} else {
			text = text + " " + comment
		}
	}
	return text
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

func hasContinuationCharacter(comment string) bool {
	stripped := strings.TrimSpace(comment[2:])
	return strings.HasSuffix(stripped, "\\")
}
