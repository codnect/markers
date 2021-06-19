package marker

import (
	"errors"
	"go/ast"
)

type Collector struct {
	*Registry
}

func NewCollector(registry *Registry) *Collector {
	return &Collector{
		registry,
	}
}

func (collector *Collector) Collect(pkg *Package) error {

	if pkg == nil {
		return errors.New("pkg(package) cannot be nil")
	}

	collector.collectPackageMarkers(pkg)

	return nil
}

func (collector *Collector) collectPackageMarkers(pkg *Package) map[ast.Node][]markerComment {
	packageMarkers := make(map[ast.Node][]markerComment)

	for _, file := range pkg.Syntax {
		fileMarkers := collector.collectFileMarkers(file)

		for node, markers := range fileMarkers {
			packageMarkers[node] = append(packageMarkers[node], markers...)
		}
	}

	return packageMarkers
}

func (collector *Collector) collectFileMarkers(file *ast.File) map[ast.Node][]markerComment {
	visitor := newVisitor(file.Comments)
	ast.Walk(visitor, file)

	return nil
}

func getCommentsForNode(node ast.Node) (docCommentGroup *ast.CommentGroup, lastCommentGroup *ast.CommentGroup) {

	switch typedNode := node.(type) {
	case *ast.File:
		docCommentGroup = typedNode.Doc
	case *ast.ImportSpec:
		docCommentGroup = typedNode.Doc
		lastCommentGroup = typedNode.Comment
	case *ast.TypeSpec:
		docCommentGroup = typedNode.Doc
		lastCommentGroup = typedNode.Comment
	case *ast.GenDecl:
		docCommentGroup = typedNode.Doc
	case *ast.Field:
		docCommentGroup = typedNode.Doc
		lastCommentGroup = typedNode.Comment
	case *ast.FuncDecl:
		docCommentGroup = typedNode.Doc
	case *ast.ValueSpec:
		docCommentGroup = typedNode.Doc
		lastCommentGroup = typedNode.Comment
	default:
		lastCommentGroup = nil
	}

	return docCommentGroup, lastCommentGroup
}
