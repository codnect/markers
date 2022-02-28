// +import=marker, Pkg="github.com/procyon-projects/marker@1.2.4:command"
// +import=marker, Alias=marker2, Pkg="github.com/procyon-project/marker@1.2.4:command"
// +import=chrono, Pkg="github.com/procyon-projects/chrono:chrono"

// This is a go document comment
// +marker:package-level1
// +marker2:package-level2
// This is a go document comment
package package1

import (
	"github.com/procyon-projects/marker/test/package2"
	_ "strings"
)

type Base struct {
	Name package2.AnotherType
}

// This is a go document comment
// +marker:method-level
// This is a go document comment
// +deprecated This method is deprecated
func (f Fruit) Name() string {
	return ""
}

func (f *Fruit) String() interface{} {
	return ""
}

// This is a go document comment
// +marker:type-level
// +marker:struct-level
// This is a go document comment
// +deprecated This struct is deprecated
type Fruit struct {
	//package2.IFace
	//Base
	x *string `tag1=val1,tag2=val2`
	// This is a comment
	// +marker:field-level
	// This is a comment

	// This is a go document comment
	// +marker:field-level
	// This is a go document comment
	// +deprecated This field is deprecated
	Apple interface {
		Name(x, y int) error
	}
	// This is a comment
	// +marker:field-level
	// This is a comment

	// This is a go document comment
	// +marker:field-level
	// This is a go document comment
	Blackberry <-chan string
}

// This is a comment
// +marker:method-level
// This is a comment

func (f *Fruit) V(x int) {

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
