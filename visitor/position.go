package visitor

import (
	"github.com/procyon-projects/marker/packages"
	"go/token"
)

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
