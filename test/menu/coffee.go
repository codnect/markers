// +import=marker, Pkg=github.com/procyon-projects/marker
// +marker:package-level:Name=coffee.go

package menu

type Coffee int

const (
	Cappuccino Coffee = -(iota + 1)
	Americano
	Latte
	TurkishCoffee
)
