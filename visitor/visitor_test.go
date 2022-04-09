package visitor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"testing"
)

type PackageLevel struct {
	Name string `marker:"Name"`
}

type StructTypeLevel struct {
	Name string `marker:"Name"`
}

type StructMethodLevel struct {
	Name string `marker:"Name"`
}

type StructFieldLevel struct {
	Name string `marker:"Name"`
}

type InterfaceTypeLevel struct {
	Name string `marker:"Name"`
}

type InterfaceMethodLevel struct {
	Name string `marker:"Name"`
}

type FunctionLevel struct {
	Name string `marker:"Name"`
}

func TestEachFile(t *testing.T) {
	result, _ := packages.LoadPackages("../test/package1")
	registry := marker.NewRegistry()

	registry.Register("marker:package-level", "github.com/procyon-projects/marker", marker.PackageLevel, &PackageLevel{})
	registry.Register("marker:interface-type-level", "github.com/procyon-projects/marker", marker.InterfaceTypeLevel, &InterfaceTypeLevel{})
	registry.Register("marker:interface-method-level", "github.com/procyon-projects/marker", marker.InterfaceMethodLevel, &InterfaceMethodLevel{})
	registry.Register("marker:function-level", "github.com/procyon-projects/marker", marker.FunctionLevel, &FunctionLevel{})
	registry.Register("marker:struct-type-level", "github.com/procyon-projects/marker", marker.StructTypeLevel, &StructTypeLevel{})
	registry.Register("marker:struct-method-level", "github.com/procyon-projects/marker", marker.StructMethodLevel, &StructMethodLevel{})
	registry.Register("marker:struct-field-level", "github.com/procyon-projects/marker", marker.FieldLevel, &StructFieldLevel{})

	collector := marker.NewCollector(registry)

	err := EachFile(collector, result.GetPackages(), func(file *File, err error) error {
		if file.pkg.ID == "builtin" {
			return nil
		}

		iface := file.Interfaces().At(0)
		methods := iface.getInterfaceMethods()
		methods[5].Results()

		if methods != nil {

		}

		function := file.Functions().At(0)
		fResults := function.Results()
		fParams := function.Params()

		if fResults != nil {

		}

		if fParams != nil {

		}

		s := file.Structs().At(0)
		s.Fields()
		s.AllFields()
		s.NumEmbeddedFields()
		s.AllMethods()
		s.EmbeddedFields()
		sMethods := s.Methods()
		sFields := s.Fields()

		if sMethods != nil {

		}

		if sFields != nil {

		}

		return nil
	})

	if err != nil {

	}
}
