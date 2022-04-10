// +import=marker, Pkg=github.com/procyon-projects/marker
// +marker:package-level:Name=dessert.go

package package1

// BakeryShop is an interface
// +marker:interface-type-level:Name=BakeryShop
type BakeryShop interface {
	// Dessert is an embedded interface
	// +marker:interface-method-level:Name=Dessert
	Dessert

	// Bread is a method
	// +marker:interface-method-level:Name=Bread
	Bread(i, k float64) struct{}
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
	Cookie
}

// Buy is a method
// +marker:struct-method-level:Name=Buy
func (c *FriedCookie) Buy(i int) {

}

// NewYearsEveCookie is an interface
// +marker:interface-type-level:Name=NewYearsEveCookie
type NewYearsEveCookie interface {
	Funfetti(v rune) byte
}

// Cookie is a struct
// +marker:struct-type-level:Name=Cookie
type Cookie struct {
	// ChocolateChip is a field
	// +marker:struct-field-level:Name=ChocolateChip
	ChocolateChip string
}

// FortuneCookie is a method
// +marker:struct-method-level:Name=FortuneCookie
func (c *Cookie) FortuneCookie(v interface{}) []string {
	return nil
}

// Oreo is a method
// +marker:struct-method-level:Name=Oreo
func (c *Cookie) Oreo(a []interface{}, v ...string) error {
	return nil
}

// Dessert is an interface
// +marker:interface-type-level:Name=Dessert
type Dessert interface {

	// IceCream is a method
	// +marker:interface-method-level:Name=IceCream
	// +marker:interface-type-level:Name=IceCream
	IceCream(s string, v ...string) string

	// Cupcake is a method
	// +marker:interface-method-level:Name=Cupcake
	Cupcake(a []int, b bool) float32

	// Tart is a method
	// +marker:interface-method-level:Name=Tart
	Tart(s any)

	// Donut is a method
	// +marker:interface-method-level:Name=Donut
	Donut() interface{}

	// Pudding is a method
	// +marker:interface-method-level:Name=Pudding
	Pudding() []string

	// Pie is a method
	// +marker:interface-method-level:Name=Pie
	Pie() any

	// muffin is a method
	// +marker:interface-method-level:Name=muffin
	muffin() (string, error)
}

// MakeACake is a function
// +marker:function-level:Name=MakeACake
func MakeACake(s any) error {
	return nil
}

// SweetShop is an interface
// +marker:interface-type-level:Name=SweetShop
type SweetShop interface {

	// NewYearsEveCookie is an embedded interface
	// +marker:interface-method-level:Name=NewYearsEveCookie
	NewYearsEveCookie

	// Dessert is an embedded interface
	// +marker:interface-method-level:Name=Dessert
	Dessert

	// Macaron is a method
	// +marker:interface-method-level:Name=Macaron
	Macaron(i complex128) bool
}
