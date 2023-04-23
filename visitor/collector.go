package visitor

import "github.com/procyon-projects/markers/packages"

type packageCollector struct {
	hasSeen  map[string]bool
	files    map[string]*Files
	packages map[string]*packages.Package

	unprocessedTypes map[string]map[string]Type

	importTypes map[string]Type
}

func newPackageCollector() *packageCollector {
	return &packageCollector{
		hasSeen:          make(map[string]bool),
		files:            make(map[string]*Files),
		packages:         make(map[string]*packages.Package),
		unprocessedTypes: make(map[string]map[string]Type),
		importTypes:      make(map[string]Type),
	}
}

func (collector *packageCollector) markAsSeen(pkgId string) {
	collector.hasSeen[pkgId] = true
}

func (collector *packageCollector) isVisited(pkgId string) bool {
	visited, ok := collector.hasSeen[pkgId]

	if !ok {
		return false
	}

	return visited
}

func (collector *packageCollector) addFile(pkgId string, file *File) {
	if _, ok := collector.files[pkgId]; !ok {
		collector.files[pkgId] = &Files{
			elements: make([]*File, 0),
		}
	}

	if _, ok := collector.files[pkgId].FindByName(file.name); ok {
		return
	}

	collector.files[pkgId].elements = append(collector.files[pkgId].elements, file)
}

func (collector *packageCollector) findTypeByImportAndTypeName(importName, typeName string, file *File) Type {
	if importedType, ok := collector.importTypes[importName+"#"+typeName]; ok {
		return importedType
	}

	packageImport, _ := file.imports.FindByName(importName)

	if packageImport == nil {
		packageImport, _ = file.imports.FindByPath(importName)
	}

	if importedType, ok := collector.importTypes[packageImport.path+"#"+typeName]; ok {
		return importedType
	}

	typ, exists := collector.findTypeByPkgIdAndName(packageImport.path, typeName)

	if exists {
		collector.importTypes[packageImport.path+"#"+typeName] = typ
	}

	return typ
}

func (collector *packageCollector) findTypeByPkgIdAndName(pkgId, typeName string) (Type, bool) {
	if files, ok := collector.files[pkgId]; ok {

		for i := 0; i < files.Len(); i++ {
			file := files.At(i)

			if structType, ok := file.structs.FindByName(typeName); ok {
				return structType, true
			}

			if interfaceType, ok := file.interfaces.FindByName(typeName); ok {
				return interfaceType, true
			}

			if customType, ok := file.customTypes.FindByName(typeName); ok {
				return customType, true
			}

			if constant, ok := file.constants.FindByName(typeName); ok {
				return constant, true
			}
		}

	} else if !collector.isVisited(pkgId) {
		loadResult, err := packages.LoadPackages(pkgId)

		if err != nil {
			panic(err)
		}

		pkg, _ := loadResult.Lookup(pkgId)

		visitPackage(pkg, collector, nil) //collector.allPackageMarkers)

		typ, ok := collector.findTypeByPkgIdAndName(pkgId, typeName)

		if ok {
			return typ, true
		}
	}

	if typ, ok := collector.unprocessedTypes[pkgId][typeName]; ok {
		return typ, true
	}

	return nil, false
}
