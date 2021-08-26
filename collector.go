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

func (collector *Collector) Collect(pkg *Package) (map[ast.Node]MarkerValues, error) {

	if pkg == nil {
		return nil, errors.New("pkg(package) cannot be nil")
	}

	nodeMarkers := collector.collectPackageMarkerComments(pkg)
	markers, err := collector.parseMarkerComments(pkg, nodeMarkers)

	if err != nil {
		return nil, err
	}

	return markers, nil
}

func (collector *Collector) collectPackageMarkerComments(pkg *Package) map[ast.Node][]markerComment {
	packageNodeMarkers := make(map[ast.Node][]markerComment)

	for _, file := range pkg.Syntax {
		fileNodeMarkers := collector.collectFileMarkerComments(file)

		for node, markers := range fileNodeMarkers {
			packageNodeMarkers[node] = append(packageNodeMarkers[node], markers...)
		}
	}

	return packageNodeMarkers
}

func (collector *Collector) collectFileMarkerComments(file *ast.File) map[ast.Node][]markerComment {
	visitor := newVisitor(file.Comments)
	ast.Walk(visitor, file)
	visitor.nodeMarkers[file] = visitor.packageMarkers

	return visitor.nodeMarkers
}

func (collector *Collector) parseMarkerComments(pkg *Package, nodeMarkerComments map[ast.Node][]markerComment) (map[ast.Node]MarkerValues, error) {
	var errs []error
	nodeMarkerValues := make(map[ast.Node]MarkerValues)

	for node, markerComments := range nodeMarkerComments {

		markerValues := make(MarkerValues)

		for _, markerComment := range markerComments {
			markerText := markerComment.Text()
			definition := collector.Lookup(markerText)

			if definition == nil {
				continue
			}

			switch typedNode := node.(type) {
			case *ast.File:

				if definition.Level&PackageLevel != PackageLevel {
					continue
				}

			case *ast.TypeSpec:

				_, isStructType := typedNode.Type.(*ast.StructType)

				if isStructType && (definition.Level&TypeLevel != TypeLevel && definition.Level&StructTypeLevel != StructTypeLevel) {
					continue
				}

				_, isInterfaceType := typedNode.Type.(*ast.InterfaceType)

				if isInterfaceType && (definition.Level&TypeLevel != TypeLevel && definition.Level&InterfaceTypeLevel != InterfaceTypeLevel) {
					continue
				}

			case *ast.Field:

				_, isFuncType := typedNode.Type.(*ast.FuncType)

				if !isFuncType && definition.Level&FieldLevel != FieldLevel {
					continue
				} else if isFuncType && !(definition.Level&MethodLevel != MethodLevel || definition.Level&InterfaceMethodLevel != InterfaceMethodLevel) {
					continue
				}

			case *ast.FuncDecl:

				if typedNode.Recv != nil && !(definition.Level&MethodLevel != MethodLevel || definition.Level&StructMethodLevel != StructMethodLevel) {
					continue
				} else if typedNode.Recv == nil && definition.Level&FunctionLevel != FunctionLevel {
					continue
				}

			}

			value, err := definition.Parse(markerText)

			if err != nil {
				position := pkg.Fset.Position(node.Pos())
				errs = append(errs, toParseError(err, markerComment, position))
				continue
			}

			markerValues[definition.Name] = append(markerValues[definition.Name], value)
		}

		nodeMarkerValues[node] = markerValues
	}

	return nodeMarkerValues, NewErrorList(errs)
}
