// +import=marker, Pkg=github.com/procyon-projects/marker
// +marker:package-level

package package1

// Cookie is a struct
// +marker:struct-type-level
type Cookie struct {
	// ChocolateChip is a field
	// +marker:struct-field-level
	ChocolateChip string
}

// FortuneCookie is a method
// +marker:struct-method-level
func (c *Cookie) FortuneCookie(v interface{}) []string {
	return nil
}

// Oreo is a method
// +marker:struct-method-level
func (c *Cookie) Oreo(a []interface{}, v ...string) error {
	return nil
}

// FriedCookie is a struct
type FriedCookie struct {
	Cookie
}

// Buy is a method
// +marker:struct-method-level
func (c *FriedCookie) Buy(i int) {

}

// Dessert is an interface
// +marker:interface-type-level
type Dessert interface {

	// IceCream is a method
	// +marker:interface-method-level
	// +marker:interface-type-level
	IceCream(s string, v ...string) string

	// Cupcake is a method
	// +marker:interface-method-level
	Cupcake(a []int, b bool) float32

	// Tart is a method
	// +marker:interface-method-level
	Tart(s any)

	// Donut is a method
	// +marker:interface-method-level
	Donut() interface{}

	// Pudding is a method
	// +marker:interface-method-level
	Pudding() []string

	// Pie is a method
	// +marker:interface-method-level
	Pie() any

	// muffin is a method
	// +marker:interface-method-level
	muffin() (string, error)
}

// MakeACake is a function
// +marker:function-level
func MakeACake(s any) error {
	return nil
}
