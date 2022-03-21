package marker

import (
	"testing"
)

type HttpStatus int

const (
	OKAY HttpStatus = -(ACCESS_DENIED + 1)
	NOTFOUND
)

const ACCESS_DENIED HttpStatus = 4

type TestOutput struct {
}

func TestCollector_Collect(t *testing.T) {
	c := NOTFOUND
	if c == OKAY {

	}
	result, _ := LoadPackages("std")
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
