package marker

import "testing"

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	registry.Register("test-marker:field-level", TargetsField)
	registry.Register("test-marker:doc", TargetsType|TargetsField|TargetsFunction|TargetsMethod)
}
