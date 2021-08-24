package marker

import "testing"

type TestOutput struct {
	MinValue interface{} `marker:"Min"`
}

func TestCollector_Collect(t *testing.T) {
	pkgs, _ := LoadPackages("./test-packages/foo")
	registry := NewRegistry()
	registry.Register("test-marker:method-level", MethodLevel, TestOutput{})
	registry.Register("test-marker:package-level", PackageLevel, TestOutput{})
	registry.Register("test-marker:type-level", TypeLevel, TestOutput{})
	registry.Register("test-marker:function-level", FunctionLevel, TestOutput{})
	collector := NewCollector(registry)

	EachFile(collector, pkgs[0], func(file *File) {
		if file == nil {

		}
	})
}
