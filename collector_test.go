package marker

import "testing"

type TestOutput struct {
	MinValue []string `marker:"Min"`
}

func TestCollector_Collect(t *testing.T) {
	pkgs, _ := LoadPackages("./test-packages/foo")
	registry := NewRegistry()
	registry.Register("test-maker:method-level", MethodLevel, TestOutput{})
	collector := NewCollector(registry)
	collector.Collect(pkgs[0])
}
