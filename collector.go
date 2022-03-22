package marker

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/marker/packages"
	"go/ast"
	"go/token"
	"strings"
)

type Collector struct {
	*Registry
}

func NewCollector(registry *Registry) *Collector {
	return &Collector{
		registry,
	}
}

func (collector *Collector) Collect(pkg *packages.Package) (map[ast.Node]MarkerValues, error) {

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

func (collector *Collector) collectPackageMarkerComments(pkg *packages.Package) map[ast.Node][]markerComment {
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
	visitor := newCommentVisitor(file.Comments)
	ast.Walk(visitor, file)
	visitor.nodeMarkers[file] = visitor.packageMarkers

	return visitor.nodeMarkers
}

func (collector *Collector) parseMarkerComments(pkg *packages.Package, nodeMarkerComments map[ast.Node][]markerComment) (map[ast.Node]MarkerValues, error) {
	importNodeMarkers, err := collector.parseImportMarkerComments(pkg, nodeMarkerComments)

	if err != nil {
		return nil, err
	}

	nodeMarkerValues := make(map[ast.Node]MarkerValues)

	if importNodeMarkers != nil {
		for importNode, importMarker := range importNodeMarkers {
			nodeMarkerValues[importNode] = importMarker
		}
	}

	var fileImportAliases map[*token.File]AliasMap
	fileImportAliases, err = collector.extractFileImportAliases(pkg, importNodeMarkers)

	if err != nil {
		return nil, err
	}

	var errs []error
	for node, markerComments := range nodeMarkerComments {

		markerValues := make(MarkerValues)
		file := pkg.Fset.File(node.Pos())
		importAliases := fileImportAliases[file]

		for _, markerComment := range markerComments {
			markerText := markerComment.Text()

			// first we need to check if there is any import
			aliasName, _, _ := splitMarker(markerText)
			aliasName = strings.Split(aliasName, ":")[0]
			// markers can be syntax free such as +build
			aliasName = strings.Split(aliasName, " ")[0]

			var definition *Definition
			if importMarker, ok := importAliases[aliasName]; ok {
				markerText = strings.Replace(markerText, fmt.Sprintf("+%s", aliasName), fmt.Sprintf("+%s", importMarker.Value), 1)
				definition = collector.Lookup(markerText, importMarker.GetPkgId())
			} else {
				definition = collector.Lookup(markerText, "")
			}

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
				position := pkg.Fset.Position(markerComment.Pos())
				errs = append(errs, toParseError(err, markerComment, position))
				continue
			}

			if marker, ok := value.(Marker); ok {
				err = marker.Validate()
			}

			if err != nil {
				position := pkg.Fset.Position(markerComment.Pos())
				errs = append(errs, toParseError(err, markerComment, position))
				continue
			}

			markerValues[definition.Name] = append(markerValues[definition.Name], value)
		}

		if len(markerValues) != 0 {
			nodeMarkerValues[node] = markerValues
		}

	}

	return nodeMarkerValues, NewErrorList(errs)
}

func (collector *Collector) parseImportMarkerComments(pkg *packages.Package, nodeMarkerComments map[ast.Node][]markerComment) (map[ast.Node]MarkerValues, error) {
	var errs []error
	importNodeMarkers := make(map[ast.Node]MarkerValues)

	for node, markerComments := range nodeMarkerComments {

		markerValues := make(MarkerValues)

		for _, markerComment := range markerComments {
			markerText := markerComment.Text()
			definition := collector.Lookup(markerText, "")

			if definition == nil {
				continue
			}

			if ImportMarkerName != definition.Name {
				continue
			}

			value, err := definition.Parse(markerText)

			if err != nil {
				position := pkg.Fset.Position(markerComment.Pos())
				errs = append(errs, toParseError(err, markerComment, position))
				continue
			}

			if marker, ok := value.(Marker); ok {
				err = marker.Validate()
			}

			if err != nil {
				position := pkg.Fset.Position(markerComment.Pos())
				errs = append(errs, toParseError(err, markerComment, position))
				continue
			}

			markerValues[definition.Name] = append(markerValues[definition.Name], value)
		}

		if len(markerValues) != 0 {
			importNodeMarkers[node] = markerValues
		}

	}

	return importNodeMarkers, NewErrorList(errs)
}

type AliasMap map[string]ImportMarker

func (collector *Collector) extractFileImportAliases(pkg *packages.Package, importNodeMarkers map[ast.Node]MarkerValues) (map[*token.File]AliasMap, error) {
	var errs []error
	var fileImportAliases = make(map[*token.File]AliasMap, 0)

	if importNodeMarkers == nil {
		return fileImportAliases, nil
	}

	for node, markerValues := range importNodeMarkers {
		file := pkg.Fset.File(node.Pos())

		if file == nil {
			continue
		}

		markers, ok := markerValues[ImportMarkerName]

		if !ok {
			continue
		}

		aliasMap := make(AliasMap, 0)
		pkgIdMap := make(map[string]bool, 0)

		for _, marker := range markers {
			importMarker := marker.(ImportMarker)

			if _, ok := pkgIdMap[importMarker.GetPkgId()]; ok {
				position := pkg.Fset.Position(node.Pos())
				err := fmt.Errorf("processor with Pkg '%s' has alrealdy been imported", importMarker.GetPkgId())
				errs = append(errs, toParseError(err, node, position))
				continue
			}

			pkgIdMap[importMarker.GetPkgId()] = true

			if importMarker.Alias == "" {
				aliasMap[importMarker.Value] = importMarker
			} else {
				aliasMap[importMarker.Alias] = importMarker
			}
		}

		fileImportAliases[file] = aliasMap
	}

	return fileImportAliases, NewErrorList(errs)
}
