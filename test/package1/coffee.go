package package1

type Lemonade uint

const (
	ClassicLemonade Lemonade = iota
	BlueberryLemonade
	WatermelonLemonade
	MangoLemonade
	StrawberryLemonade
)

type Coffee int

const (
	Cappuccino Coffee = -(iota + 1)
	Americano
	Latte
	TurkishCoffee
)
