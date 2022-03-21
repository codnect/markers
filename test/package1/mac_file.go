// +import=marker, Pkg="github.com/procyon-projects/marker@1.2.4:command"
// +import=marker, Alias=marker2, Pkg="github.com/procyon-project/marker@1.2.4:command"
// +import=chrono, Pkg="github.com/procyon-projects/chrono:chrono"

// This is a go document comment
// +marker:package-level1
// +marker2:package-level2
// This is a go document comment
package package1

import (
	_ "strings"
)

/*
type ABC int

const (
	Z     = 0
	C     = true || true
	X 		ABC = iota
	Y
	M = 2.4
)*/

type HttpStatus int

const (
	TEST HttpStatus = 2
	ACCESS_DENIED
)

const (
	OKAY = 23
	NOTFOUND
	X = "hello" + "hello2"
)

type Test interface {
	Dessert
}

func (a app) sort() {
}

func (f Fruit) Name() string {
	return ""
}

type app Fruit

type Base struct {
	Name int
}

// This is a go document comment
// +marker:method-level
// This is a go document comment
// +deprecated This method is deprecated

func (f *Fruit) String() interface{} {
	return ""
}

type y struct {
}

func (yx *y) Name(x, y int) error {
	return nil
}

func x() {
	d := &Fruit{
		Apple: &y{},
	}
	d.Apple.Name(3, 4)
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
func Coconut() *Fruit {
	return &Fruit{}
}

type Request struct {
}
