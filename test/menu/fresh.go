// +import=marker, Pkg=github.com/procyon-projects/markers
// +import=test-marker, Pkg=github.com/procyon-projects/test-markers
// +test-marker:package-level:Name=fresh.go

package menu

type Lemonade uint

const (
	ClassicLemonade Lemonade = iota
	BlueberryLemonade
	WatermelonLemonade
	MangoLemonade
	StrawberryLemonade
)
