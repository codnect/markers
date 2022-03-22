package visitor

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"testing"
)

type TestOutput struct {
}

func TestEachFile(t *testing.T) {

	result, _ := packages.LoadPackages("std")
	registry := marker.NewRegistry()

	registry.Register("marker:package-level1", "github.com/procyon-projects/marker", marker.PackageLevel, &TestOutput{})
	registry.Register("marker:package-level2", "github.com/procyon-project/marker", marker.PackageLevel, &TestOutput{})

	collector := marker.NewCollector(registry)

	err := EachFile(collector, result.GetPackages(), func(file *File, err error) error {
		return nil
	})

	if err != nil {

	}
}
