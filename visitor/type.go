package visitor

import (
	"github.com/procyon-projects/marker/packages"
	"strings"
)

type Type interface {
	Underlying() Type
	String() string
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

type Position struct {
	Line   int
	Column int
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
