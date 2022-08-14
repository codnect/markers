package package2

// Restaurant is an interface
// +marker:interface-type-level:Name=Restaurant
type Restaurant interface {
	// SpanishDishes is a method
	// +marker:interface-method-level:Name=SpanishDishes
	SpanishDishes() SpanishDishes
	// ItalianDishes is a method
	// +marker:interface-method-level:Name=ItalianDishes
	ItalianDishes() ItalianDishes
}

// ItalianDishes is a struct
// +marker:struct-type-level:Name=ItalianDishes
type ItalianDishes struct {
	// MinestroneSoup is a field
	// +marker:struct-field-level:Name=MinestroneSoup
	MinestroneSoup interface{}
	// Lasagna is a field
	// +marker:struct-field-level:Name=Lasagna
	Lasagna *string
	// ItalianMeatballs is a field
	// +marker:struct-field-level:Name=ItalianMeatballs
	ItalianMeatballs *int
	// FettuccineAlfredo is a field
	// +marker:struct-field-level:Name=FettuccineAlfredo
	FettuccineAlfredo [5]string
}
