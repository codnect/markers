package marker

import "testing"

func TestCollector_Collect(t *testing.T) {
	pkgs, _ := LoadPackages("./test-packages/foo")
	registry := NewRegistry()
	collector := NewCollector(registry)
	collector.Collect(pkgs[0])
}
