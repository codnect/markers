// +import=marker, Pkg=github.com/procyon-projects/markers
// +import=test-marker, Pkg=github.com/procyon-projects/test-markers
// +test-marker:package-level:Name=coffee.go

package menu

type Coffee int

const (
	Cappuccino Coffee = -(iota + 1)
	Americano
	Latte
	TurkishCoffee
)

func (c *cookie) PrintCookie(v interface{}) []string {
	return nil
}

type CustomBakeryShop BakeryShop
