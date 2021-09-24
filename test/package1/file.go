// +build windows
// This is a comment
// +marker:package-level
// This is a comment

// This is a go document comment
// +marker:package-level
// This is a go document comment
package package1

// +import=marker, Pkg="github.com/procyon-projects/marker@1.2.4:command"
// +import=chrono, Pkg="github.com/procyon-projects/chrono:chrono"
import ()

// This is a comment
// +marker:type-level
// +chrono:struct-level
// This is a comment

// This is a go document comment
// +marker:type-level
// +marker:struct-level
// This is a go document comment
// +deprecated This struct is deprecated
type Fruit struct {
	// This is a comment
	// +marker:field-level
	// This is a comment

	// This is a go document comment
	// +marker:field-level
	// This is a go document comment
	// +deprecated This field is deprecated
	Apple string
	// This is a comment
	// +marker:field-level
	// This is a comment

	// This is a go document comment
	// +marker:field-level
	// This is a go document comment
	Blackberry string
}

// This is a comment
// +marker:method-level
// This is a comment

// This is a go document comment
// +marker:method-level
// This is a go document comment
// +deprecated This method is deprecated
func (f *Fruit) Name() {

}

// This is a comment
// +marker:function-level
// This is a comment

// This is a go document comment
// +marker:function-level
// This is a go document comment
// +deprecated This function is deprecated
func Coconut() {

}

// This is a comment
// +marker:type-level
// +marker:interface-level
// This is a comment

// This is a go document comment
// +marker:type-level
// +marker:interface-level
// This is a go document comment
// +deprecated This interface is deprecated
type Dessert interface {
	// This is a comment
	// +marker:method-level
	// This is a comment

	// This is a go document comment
	// +marker:method-level
	// This is a go document comment
	IceCream(s string) string
	// This is a comment
	// +marker:method-level
	// This is a comment

	// This is a go document comment
	// +marker:method-level
	// This is a go document comment
	Cupcake() int
	// This is a comment
	// +marker:method-level
	// This is a comment

	// This is a go document comment
	// +marker:method-level
	// This is a go document comment
	Tart()
	// This is a comment
	// +marker:method-level
	// This is a comment

	// This is a go document comment
	// +marker:method-level
	// This is a go document comment
	Donut() interface{}
}
