package package2

import "github.com/procyon-projects/marker/test/package1"

// SpanishDishes is a custom type
// +marker:struct-type-level:Name=SpanishDishes
type SpanishDishes int

// SouthKoreanDishes is a struct
// +marker:struct-type-level:Name=SouthKoreanDishes
type SouthKoreanDishes struct {
	// Bulgogi is a field
	// +marker:struct-field-level:Name=Bulgogi
	Bulgogi string
}

// Kimchi is a method
// +marker:struct-method-level:Name=Kimchi
func (s SouthKoreanDishes) Kimchi() interface{} {
	return nil
}

// Menu is a struct
// +marker:struct-type-level:Name=Menu
type Menu struct {
	// Desserts is a field
	// +marker:struct-field-level:Name=Desserts
	Desserts package1.Dessert
	// ItalianDishes is a field
	// +marker:struct-field-level:Name=ItalianDishes
	ItalianDishes ItalianDishes
}

const (
	Tortilla SpanishDishes = iota
	Gazpacho
	Paella
)

// Bibimbap is a method
// +marker:struct-method-level:Name=Bibimbap
func (s SouthKoreanDishes) Bibimbap(v []string) {

}
