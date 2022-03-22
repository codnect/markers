package visitor

import (
	"github.com/procyon-projects/marker/packages"
	"go/token"
	"strings"
)

type Type interface {
	Underlying() Type
	String() string
}

type Position struct {
	Line   int
	Column int
}

func getPosition(pkg *packages.Package, tokenPosition token.Pos) Position {
	position := pkg.Fset.Position(tokenPosition)
	return Position{
		Line:   position.Line,
		Column: position.Column,
	}
}

type ImportedType struct {
	pkg *packages.Package
	typ Type
}

func (i *ImportedType) Package() *packages.Package {
	return i.pkg
}

func (i *ImportedType) Underlying() Type {
	return i.typ
}

func (i *ImportedType) String() string {
	return ""
}

func (i *ImportedType) Name() string {
	return ""
}

type Pointer struct {
	base Type
}

func (p *Pointer) Elem() Type {
	return p.base
}

func (p *Pointer) Underlying() Type {
	return p
}

func (p *Pointer) String() string {
	var builder strings.Builder
	builder.WriteString("*")
	builder.WriteString(p.base.String())
	return builder.String()
}
