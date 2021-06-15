package foo

// normal doc
// +test-marker:type-level
// +test-marker:doc
// normal doc
type Foo struct {
	// normal doc
	// +test-maker:field-level
	// +test-marker:doc
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
	// +test-marker:doc
	// normal doc
	TestFunction()
}

// normal doc
// +test-marker:function-level
// +test-marker:doc
// normal doc
func Baz() error {
	return nil
}
