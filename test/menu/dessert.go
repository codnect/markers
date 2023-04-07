// +import=marker, Pkg=github.com/procyon-projects/markers
// +marker:package-level:Name=dessert.go

package menu

import (
	"fmt"
	_ "strings"
)

// BakeryShop is an interface
// +marker:interface-type-level:Name=BakeryShop
type BakeryShop interface {
	// Bread is a method
	// +marker:interface-method-level:Name=Bread
	Bread(i, k float64) struct{}
	// Dessert is an embedded interface
	// +marker:interface-method-level:Name=Dessert
	Dessert
}

// Eat is a method
// +marker:struct-method-level:Name=Eat
func (c *FriedCookie) Eat() bool {
	return true
}

// FriedCookie is a struct
// +marker:struct-type-level:Name=FriedCookie
type FriedCookie struct {

	// Cookie is an embedded struct
	// +marker:interface-method-level:Name=Cookie
	cookie
	// ChocolateChip is a field
	// +marker:struct-field-level:Name=CookieDough
	cookieDough any
}

// Buy is a method
// +marker:struct-method-level:Name=Buy
func (c *FriedCookie) Buy(i int) {

}

// NewYearsEveCookie is an interface
// +marker:interface-type-level:Name=newYearsEveCookie
type newYearsEveCookie interface {
	// Funfetti is a method
	// +marker:interface-method-level:Name=Funfetti
	Funfetti(v rune) byte
}

// Cookie is a struct
// +marker:struct-type-level:Name=cookie, Any={key:"value"}
type cookie struct {
	// ChocolateChip is a field
	// +marker:struct-field-level:Name=ChocolateChip
	ChocolateChip string
	// tripleChocolateCookie is a field
	// +marker:struct-field-level:Name=tripleChocolateCookie
	tripleChocolateCookie map[string]error
}

// FortuneCookie is a method
// +marker:struct-method-level:Name=FortuneCookie
func (c *cookie) FortuneCookie(v interface{}) []string {
	return nil
}

// Oreo is a method
// +marker:struct-method-level:Name=Oreo
func (c *cookie) Oreo(a []interface{}, v ...string) error {
	return nil
}

// Dessert is an interface
// +marker:interface-type-level:Name=Dessert
type Dessert interface {

	// IceCream is a method
	// +marker:interface-method-level:Name=IceCream
	// +marker:interface-type-level:Name=IceCream
	IceCream(s string, v ...bool) (r string)

	// CupCake is a method
	// +marker:interface-method-level:Name=CupCake
	CupCake(a []int, b bool) float32

	// Tart is a method
	// +marker:interface-method-level:Name=Tart
	Tart(s interface{})

	// Donut is a method
	// +marker:interface-method-level:Name=Donut
	Donut() error

	// Pudding is a method
	// +marker:interface-method-level:Name=Pudding
	Pudding() [5]string

	// Pie is a method
	// +marker:interface-method-level:Name=Pie
	Pie() interface{}

	// muffin is a method
	// +marker:interface-method-level:Name=muffin
	muffin() (*string, error)
}

// MakeACake is a function
// +marker:function-level:Name=MakeACake
func MakeACake(s interface{}) error {
	return nil
}

// BiscuitCake is a function
// +marker:function-level:Name=BiscuitCake
func BiscuitCake(s string, arr []int, v ...int16) (i int, b bool) {
	return
}

// SweetShop is an interface
// +marker:interface-type-level:Name=SweetShop
type SweetShop interface {

	// NewYearsEveCookie is an embedded interface
	// +marker:interface-method-level:Name=NewYearsEveCookie
	newYearsEveCookie

	// Macaron is a method
	// +marker:interface-method-level:Name=Macaron
	Macaron(c complex128) (chan string, fmt.Stringer)

	// Dessert is an embedded interface
	// +marker:interface-method-level:Name=Dessert
	Dessert
}
