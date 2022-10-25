package markers

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

func (collector *Collector) Collect(pkg *packages.Package) (map[ast.Node]Values, error) {

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

func (collector *Collector) parseMarkerComments(pkg *packages.Package, nodeMarkerComments map[ast.Node][]markerComment) (map[ast.Node]Values, error) {
	importNodeMarkers, err := collector.parseImportMarkerComments(pkg, nodeMarkerComments)

	if err != nil {
		return nil, err
	}

	nodeMarkerValues := make(map[ast.Node]Values)

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

		markerValues := make(Values)
		file := pkg.Fset.File(node.Pos())
		importAliases := fileImportAliases[file]

		for _, markerComment := range markerComments {
			markerText := markerComment.Text()
			markerName, _, _ := splitMarker(markerText)
			targetLevel := FindTargetLevel(node)
			alias := strings.SplitN(markerName, ":", 2)[0]

			var definition *Definition
			var exists bool

			if importMarker, ok := importAliases[alias]; ok {
				markerName = strings.Replace(markerName, fmt.Sprintf("+%s", alias), fmt.Sprintf("+%s", importMarker.Value), 1)
				definition, exists = collector.Lookup(markerName, importMarker.Pkg, targetLevel)
			} else if _, isReservedMarker := reservedMarkerMap[markerName]; !isReservedMarker {
				continue
			} else {
				definition, exists = collector.Lookup(markerName, "", targetLevel)
			}

			if !exists {
				continue
			}

			value, err := definition.Parse(markerText)

			if err != nil {
				position := pkg.Fset.Position(markerComment.Pos())
				errs = append(errs, toParseError(err, markerComment, position))
				continue
			}

			/*if marker, ok := value.(Marker); ok {
				err = marker.Validate()
			}*/

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

func (collector *Collector) parseImportMarkerComments(pkg *packages.Package, nodeMarkerComments map[ast.Node][]markerComment) (map[ast.Node]Values, error) {
	var errs []error
	importNodeMarkers := make(map[ast.Node]Values)

	for node, markerComments := range nodeMarkerComments {

		markerValues := make(Values)

		for _, markerComment := range markerComments {
			markerText := markerComment.Text()
			name, anonymousName, fields := splitMarker(markerText)

			if fields == "" {

			}

			if ImportMarkerName != name || ImportMarkerName != anonymousName {
				continue
			}

			definition, exists := collector.Lookup("import", "", PackageLevel)

			if !exists {
				continue
			}

			value, err := definition.Parse(markerText)

			if err != nil {
				position := pkg.Fset.Position(markerComment.Pos())
				errs = append(errs, toParseError(err, markerComment, position))
				continue
			}

			/*if marker, ok := value.(Marker); ok {
				err = marker.Validate()
			}*/

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

type AliasMap map[string]Import

func (collector *Collector) extractFileImportAliases(pkg *packages.Package, importNodeMarkers map[ast.Node]Values) (map[*token.File]AliasMap, error) {
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
			importMarker := marker.(Import)

			if _, ok := pkgIdMap[importMarker.Pkg]; ok {
				position := pkg.Fset.Position(node.Pos())
				err := fmt.Errorf("processor with Pkg '%s' has alrealdy been imported", importMarker.Pkg)
				errs = append(errs, toParseError(err, node, position))
				continue
			}

			pkgIdMap[importMarker.Pkg] = true

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
