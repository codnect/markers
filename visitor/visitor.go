package visitor

import (
	"errors"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"go/token"
)

type FileCallback func(file *File, err error) error

type packageVisitor struct {
	collector *packageCollector

	pkg               *packages.Package
	packageMarkers    map[ast.Node]marker.MarkerValues
	allPackageMarkers map[string]map[ast.Node]marker.MarkerValues

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
		newFunction(typedNode, nil, visitor.file, visitor.pkg, visitor)
		return nil
	case *ast.TypeSpec:
		collectTypeFromTypeSpec(typedNode, visitor)
		return nil
	default:
		return nil
	}
}

func visitPackage(pkg *packages.Package, collector *packageCollector, allPackageMarkers map[string]map[ast.Node]marker.MarkerValues) {
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

func EachFile(collector *marker.Collector, pkgs []*packages.Package, callback FileCallback) error {
	if collector == nil {
		return errors.New("collector cannot be nil")
	}

	if pkgs == nil {
		return errors.New("packages cannot be nil")
	}

	var errs []error
	packageMarkers := make(map[string]map[ast.Node]marker.MarkerValues)

	for _, pkg := range pkgs {
		markers, err := collector.Collect(pkg)

		if err != nil {
			errs = append(errs, err.(marker.ErrorList)...)
			continue
		}

		packageMarkers[pkg.ID] = markers
	}

	if len(errs) != 0 {
		return marker.NewErrorList(errs)
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
