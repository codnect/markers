package marker

import (
	"testing"
)

type TestOutput struct {
	MinValue interface{} `marker:"Min"`
}

func TestCollector_Collect(t *testing.T) {
	pkgs, _ := LoadPackages("./test/package1", "./test/package2")
	registry := NewRegistry()

	collector := NewCollector(registry)

	EachFile(collector, pkgs, func(file *File, err error) {
		if file == nil {

		}
	})
}
