package marker

import "testing"

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	//registry.Register("test-marker:field-level", FieldLevel)
	//registry.Register("test-marker:doc", TypeLevel|FieldLevel|FunctionLevel|MethodLevel)

	registry.Lookup("+test-marker:doc=true")
}
