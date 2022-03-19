package marker

import (
	"testing"
)

type TestOutput struct {
}

func TestCollector_Collect(t *testing.T) {

	result, _ := LoadPackages("./test/package1")
	registry := NewRegistry()

	registry.Register("marker:package-level1", "github.com/procyon-projects/marker", PackageLevel, &TestOutput{})
	registry.Register("marker:package-level2", "github.com/procyon-project/marker", PackageLevel, &TestOutput{})

	collector := NewCollector(registry)

	eachFile(collector, result.GetPackages(), func(file *SourceFile, err error) {

	})

	/*EachFile(collector, result.GetPackages(), func(file *File, err error) {
		if file == nil {
		}
	})*/
}
