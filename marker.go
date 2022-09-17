package marker

import (
	"go/ast"
	"go/token"
	"strings"
)

type Validate interface {
	Validate() error
}

type MarkerValues map[string][]any

func (markerValues MarkerValues) Count() int {
	return len(markerValues)
}

func (markerValues MarkerValues) AllMarkers(name string) []any {
	result := markerValues[name]

	if len(result) == 0 {
		return nil
	}

	return result
}

func (markerValues MarkerValues) First(name string) any {
	result := markerValues[name]

	if len(result) == 0 {
		return nil
	}

	return result[0]
}

func (markerValues MarkerValues) CountByName(name string) int {
	result := markerValues[name]
	return len(result)
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
		value := strings.TrimRight(nameFieldParts[0], "\t ")
		value = strings.TrimLeft(value, "\t ")
		return value, value, ""
	}

	fields := nameFieldParts[1]
	anonymousName = strings.TrimRight(nameFieldParts[0], "\t ")
	anonymousName = strings.TrimLeft(anonymousName, "\t ")
	name = anonymousName

	nameParts := strings.Split(name, ":")

	if len(nameParts) > 1 {
		name = strings.Join(nameParts[:len(nameParts)-1], ":")
	}

	if len(anonymousName) >= len(name)+1 {
		fields = anonymousName[len(name)+1:] + "=" + fields
	} else {
		fields = "=" + fields
	}

	return name, anonymousName, fields
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

func isMultiLineComment(comment string) bool {
	stripped := strings.TrimSpace(comment[2:])
	return strings.HasSuffix(stripped, "\\")
}
