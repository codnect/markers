package visitor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"testing"
)

type TestOutput struct {
}

func TestEachFile(t *testing.T) {

	result, _ := packages.LoadPackages("../test/package2")
	registry := marker.NewRegistry()

	registry.Register("marker:package-level1", "github.com/procyon-projects/marker", marker.PackageLevel, &TestOutput{})
	registry.Register("marker:package-level2", "github.com/procyon-project/marker", marker.PackageLevel, &TestOutput{})

	collector := marker.NewCollector(registry)

	err := EachFile(collector, result.GetPackages(), func(file *File, err error) error {
		function := file.Functions().At(0)
		params := function.Params()
		results := function.Results()
		isVariadic := function.IsVariadic()
		if params != nil {

		}

		if results != nil {

		}

		if isVariadic {

		}

		structType := file.Structs().At(0)
		fieldList := structType.AllFields()
		methods := structType.Methods()

		if methods != nil {

		}

		if fieldList != nil {

		}
		if structType != nil {

		}

		return nil
	})

	if err != nil {

	}
}
