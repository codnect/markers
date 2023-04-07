package visitor

import (
	"errors"
	"github.com/procyon-projects/markers"
	"github.com/procyon-projects/markers/packages"
	"go/ast"
	"go/token"
)

type FileCallback func(file *File, err error) error

type packageVisitor struct {
	collector *packageCollector

	pkg               *packages.Package
	packageMarkers    map[ast.Node]markers.Values
	allPackageMarkers map[string]map[ast.Node]markers.Values

	file *File

	genDecl  *ast.GenDecl
	funcDecl *ast.FuncDecl
	rawFile  *ast.File
}

func (visitor *packageVisitor) VisitPackage() {
	visitor.packageMarkers = visitor.allPackageMarkers[visitor.pkg.ID]
	visitor.collector.markAsSeen(visitor.pkg.ID)

	for _, file := range visitor.pkg.Syntax {
		ast.Walk(visitor, file)
	}
}

func (visitor *packageVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return visitor
	}

	switch typedNode := node.(type) {
	case *ast.File:
		visitor.rawFile = typedNode
		visitor.file = newFile(typedNode, visitor.pkg, visitor.packageMarkers[typedNode], visitor)
		visitor.collector.addFile(visitor.pkg.ID, visitor.file)
		return visitor
	case *ast.GenDecl:
		visitor.genDecl = typedNode

		if typedNode.Tok == token.CONST {
			collectConstantsFromSpecs(typedNode.Specs, visitor.file)
		}

		return visitor
	case *ast.FuncDecl:
		visitor.funcDecl = typedNode
		newFunction(typedNode, nil, nil, nil, visitor.file, visitor.pkg, visitor, visitor.packageMarkers[typedNode])
		return nil
	case *ast.TypeSpec:
		collectTypeFromTypeSpec(typedNode, visitor)
		return nil
	default:
		return nil
	}
}

func visitPackage(pkg *packages.Package, collector *packageCollector, allPackageMarkers map[string]map[ast.Node]markers.Values) {
	pkgVisitor := &packageVisitor{
		collector:         collector,
		pkg:               pkg,
		allPackageMarkers: allPackageMarkers,
	}

	if _, ok := collector.packages[pkg.ID]; !ok {
		collector.packages[pkg.ID] = pkg
	}

	pkgVisitor.VisitPackage()
}

func EachFile(collector *markers.Collector, pkgs []*packages.Package, callback FileCallback) error {
	if collector == nil {
		return errors.New("collector cannot be nil")
	}

	if pkgs == nil {
		return errors.New("packages cannot be nil")
	}

	var errs []error
	packageMarkers := make(map[string]map[ast.Node]markers.Values)

	for _, pkg := range pkgs {
		markerValues, err := collector.Collect(pkg)

		if err != nil {
			errs = append(errs, err.(markers.ErrorList)...)
			continue
		}

		packageMarkers[pkg.ID] = markerValues
	}

	if len(errs) != 0 {
		return markers.NewErrorList(errs)
	}

	pkgCollector := newPackageCollector()

	for _, pkg := range pkgs {
		if !pkgCollector.isVisited(pkg.ID) || !pkgCollector.isProcessed(pkg.ID) {
			visitPackage(pkg, pkgCollector, packageMarkers)
		}
	}

	for _, pkg := range pkgCollector.files {
		for _, file := range pkg.elements {
			callback(file, nil)
		}
	}

	return nil
}
