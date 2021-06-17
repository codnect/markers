package marker

import (
	"errors"
	"strings"
)

type TargetLevel int

const (
	TypeLevel TargetLevel = 1 << iota
	ImportLevel
	FieldLevel
	FunctionLevel
	MethodLevel
	PackageLevel
)

func splitMarker(marker string) (name string, anonymousName string, options string) {
	marker = marker[1:]

	nameFieldParts := strings.SplitN(marker, "=", 2)

	if len(nameFieldParts) == 1 {
		return nameFieldParts[0], nameFieldParts[0], ""
	}

	anonymousName = nameFieldParts[0]
	name = anonymousName

	nameParts := strings.Split(name, ":")

	if len(nameParts) > 1 {
		name = strings.Join(nameParts[:len(nameParts)-1], ":")
	}

	return name, anonymousName, nameFieldParts[1]
}

func Scan(collector *Collector, pkg *Package) error {

	if collector == nil {
		return errors.New("collector cannot be nil")
	}

	if pkg == nil {
		return errors.New("pkg(package) cannot be nil")
	}

	err := collector.Collect(pkg)

	if err != nil {
		return err
	}

	return nil
}
