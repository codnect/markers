package marker

import (
	"go/ast"
	"go/token"
)

type commentVisitor struct {
	allComments      []*ast.CommentGroup
	nextCommentIndex int

	packageMarkers     []markerComment
	importMarkers      []markerComment
	declarationMarkers []markerComment
	nodeMarkers        map[ast.Node][]markerComment
}

func newCommentVisitor(allComments []*ast.CommentGroup) *commentVisitor {
	return &commentVisitor{
		allComments: allComments,
		nodeMarkers: make(map[ast.Node][]markerComment),
	}
}

func (visitor *commentVisitor) Visit(node ast.Node) (w ast.Visitor) {

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
		docCommentGroup := visitor.getCommentsForNode(node)

		markerCommentIndex := lastCommentIndex

		if docCommentGroup != nil && visitor.allComments[markerCommentIndex].Pos() == docCommentGroup.Pos() {
			markerCommentIndex--
		}

		if markerCommentIndex >= visitor.nextCommentIndex {
			markersFromComment = visitor.getMarkerComments(markerCommentIndex, markerCommentIndex+1)
			markersFromDocument = visitor.getMarkerComments(markerCommentIndex+1, lastCommentIndex+1)
		} else {
			markersFromDocument = visitor.getMarkerComments(markerCommentIndex+1, lastCommentIndex+1)
		}
	}

	switch typedNode := node.(type) {
	case *ast.File:
		visitor.packageMarkers = append(visitor.packageMarkers, markersFromComment...)
		visitor.packageMarkers = append(visitor.packageMarkers, markersFromDocument...)
	case *ast.TypeSpec:
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], visitor.declarationMarkers...)
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromComment...)
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromDocument...)
		visitor.declarationMarkers = nil
	case *ast.GenDecl:
		if typedNode.Tok == token.IMPORT {
			visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromComment...)
			visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromDocument...)
		} else {
			visitor.declarationMarkers = append(visitor.declarationMarkers, markersFromComment...)
			visitor.declarationMarkers = append(visitor.declarationMarkers, markersFromDocument...)
		}
	case *ast.Field, *ast.FuncDecl:
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromComment...)
		visitor.nodeMarkers[node] = append(visitor.nodeMarkers[node], markersFromDocument...)
	}

	visitor.nextCommentIndex = lastCommentIndex + 1

	return visitor
}

func (visitor *commentVisitor) getMarkerComments(startIndex, endIndex int) []markerComment {
	if startIndex < 0 || endIndex < 0 {
		return nil
	}

	markerComments := make([]markerComment, 0)

	for index := startIndex; index < endIndex; index++ {
		commentGroup := visitor.allComments[index]

		var markerComment *markerComment
		var hasContinuation bool

		for _, comment := range commentGroup.List {
			containsMarker := isMarkerComment(comment.Text)

			if containsMarker {
				if markerComment != nil {
					markerComments = append(markerComments, *markerComment)
				}

				markerComment = newMarkerComment(comment)
				hasContinuation = false
			} else if hasContinuation {
				if markerComment != nil {
					markerComment.append(comment)
				}
			}

			hasContinuation = hasContinuationCharacter(comment.Text)
		}

		if markerComment != nil {
			markerComments = append(markerComments, *markerComment)
		}
	}

	return markerComments
}

func (visitor *commentVisitor) getCommentsForNode(node ast.Node) (docCommentGroup *ast.CommentGroup) {

	switch typedNode := node.(type) {
	case *ast.File:
		docCommentGroup = typedNode.Doc
	case *ast.ImportSpec:
		docCommentGroup = typedNode.Doc
	case *ast.TypeSpec:
		docCommentGroup = typedNode.Doc
	case *ast.GenDecl:
		docCommentGroup = typedNode.Doc
	case *ast.Field:
		docCommentGroup = typedNode.Doc
	case *ast.FuncDecl:
		docCommentGroup = typedNode.Doc
	case *ast.ValueSpec:
		docCommentGroup = typedNode.Doc
	}

	return docCommentGroup
}

type importCallback func(file *ast.File, decl *ast.GenDecl)
type typeCallback func(file *ast.File, decl *ast.GenDecl, node *ast.TypeSpec)
type functionCallback func(file *ast.File, decl *ast.FuncDecl, funcType *ast.FuncType)
type fileCallback func(file *ast.File)

type fileElementVisitor struct {
	pkg              *Package
	fileCallback     fileCallback
	importCallback   importCallback
	typeCallback     typeCallback
	functionCallback functionCallback
	genDecl          *ast.GenDecl
	funcDecl         *ast.FuncDecl
	file             *ast.File
}

func (visitor *fileElementVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return visitor
	}

	switch typedNode := node.(type) {
	case *ast.File:
		visitor.file = typedNode
		visitor.fileCallback(typedNode)
		return visitor
	case *ast.GenDecl:
		visitor.genDecl = typedNode
		if typedNode.Tok == token.IMPORT {
			visitor.importCallback(visitor.file, visitor.genDecl)
		}
		return visitor
	case *ast.FuncDecl:
		visitor.funcDecl = typedNode
		return visitor
	case *ast.FuncType:
		visitor.functionCallback(visitor.file, visitor.funcDecl, typedNode)
		return nil
	case *ast.TypeSpec:
		visitor.typeCallback(visitor.file, visitor.genDecl, typedNode)
		return nil
	default:
		return nil
	}
}

func visitPackageFiles(pkg *Package, fileCallback fileCallback,
	importCallback importCallback,
	typeCallback typeCallback,
	functionCallback functionCallback) {
	visitor := &fileElementVisitor{
		pkg:              pkg,
		fileCallback:     fileCallback,
		importCallback:   importCallback,
		typeCallback:     typeCallback,
		functionCallback: functionCallback,
	}

	for _, file := range pkg.Syntax {
		ast.Walk(visitor, file)
	}
}
