package deneme

// normal doc
// +test-marker:type-level
// normal doc
type Foo struct {
	// normal doc
	// +test-maker:field-level
	// normal doc
	Milkshake string
	// normal doc
	// +test-maker:field-level
	// normal doc
	Donut int
}

// normal doc
// +test-marker:type-level
// normal doc
type Bar interface {
	// normal doc
	// +test-marker:function-level
	// normal doc
	TestFunction()
}

// normal doc
// +test-marker:function-level
// normal doc
func Baz() error {
	return nil
}
