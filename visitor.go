package marker

import (
	"go/ast"
)

type Visitor struct {
	allComments      []*ast.CommentGroup
	nextCommentIndex int
	lastCommentGroup *ast.CommentGroup

	packageMarkers     []markerComment
	declarationMarkers []markerComment
	nodeMarkers        map[ast.Node][]markerComment
}

func newVisitor(allComments []*ast.CommentGroup) *Visitor {
	return &Visitor{
		allComments: allComments,
		nodeMarkers: make(map[ast.Node][]markerComment),
	}
}

func (visitor *Visitor) Visit(node ast.Node) (w ast.Visitor) {

	if node == nil {
		return nil
	}

	switch node.(type) {
	case *ast.CommentGroup:
		return nil
	case *ast.Ident:
		return nil
	case *ast.ImportSpec:
		return nil
	case *ast.FieldList:
		return visitor
	case *ast.InterfaceType:
		return visitor
	case *ast.FuncType:
		return nil
	}

	lastCommentIndex := visitor.nextCommentIndex

	var markersFromComment []markerComment
	var markersFromDocument []markerComment

	if visitor.nextCommentIndex < len(visitor.allComments) {
		nextCommentGroup := visitor.allComments[visitor.nextCommentIndex]

		for nextCommentGroup.Pos() < node.Pos() {
			lastCommentIndex++

			if lastCommentIndex >= len(visitor.allComments) {
				break
			}

			nextCommentGroup = visitor.allComments[lastCommentIndex]
		}

		lastCommentIndex--

		var docCommentGroup *ast.CommentGroup
		docCommentGroup, visitor.lastCommentGroup = getCommentsForNode(node)

		markerCommentIndex := lastCommentIndex

		if docCommentGroup != nil && visitor.allComments[markerCommentIndex].Pos() == docCommentGroup.Pos() {
			markerCommentIndex--
		}

		if markerCommentIndex >= visitor.nextCommentIndex {
			markersFromComment = visitor.getMarkers(markerCommentIndex, markerCommentIndex+1)
			markersFromDocument = visitor.getMarkers(markerCommentIndex+1, lastCommentIndex+1)
		} else {
			markersFromDocument = visitor.getMarkers(markerCommentIndex+1, lastCommentIndex+1)
		}
	}

	switch node.(type) {
	case *ast.File:
		visitor.packageMarkers = append(visitor.packageMarkers, markersFromComment...)
		visitor.packageMarkers = append(visitor.packageMarkers, markersFromDocument...)
	case *ast.TypeSpec:
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], visitor.declarationMarkers...)
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromComment...)
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromDocument...)
		visitor.declarationMarkers = nil
	case *ast.GenDecl:
		visitor.declarationMarkers = append(visitor.declarationMarkers, markersFromComment...)
		visitor.declarationMarkers = append(visitor.declarationMarkers, markersFromDocument...)
	case *ast.Field:
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromComment...)
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromDocument...)
	case *ast.FuncDecl:
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromComment...)
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromDocument...)
	}

	visitor.nextCommentIndex = lastCommentIndex + 1

	return visitor
}

func (visitor *Visitor) getMarkers(startIndex, endIndex int) []markerComment {
	if startIndex < 0 || endIndex < 0 {
		return nil
	}

	markerComments := make([]markerComment, 0)

	for index := startIndex; index < endIndex; index++ {
		commentGroup := visitor.allComments[index]

		for _, comment := range commentGroup.List {
			if !isMarkerComment(comment.Text) {
				continue
			}

			markerComments = append(markerComments, newMarkerComment(comment))
		}
	}

	return markerComments
}
