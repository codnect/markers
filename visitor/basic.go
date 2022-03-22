package visitor

type BasicKind int

const (
	Invalid BasicKind = iota

	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String
	UnsafePointer

	UntypedBool
	UntypedInt
	UntypedRune
	UntypedFloat
	UntypedComplex
	UntypedString
	UntypedNil

	Byte = Uint8
	Rune = Int32
)

// BasicInfo is a set of flags describing properties of a basic type.
type BasicInfo int

// Properties of basic types.
const (
	IsBoolean BasicInfo = 1 << iota
	IsInteger
	IsUnsigned
	IsFloat
	IsComplex
	IsString
	IsUntyped

	IsNumeric = IsInteger | IsFloat | IsComplex
)

var basicTypesMap = map[string]*Basic{
	"bool":       basicTypes[Bool],
	"int":        basicTypes[Int],
	"int8":       basicTypes[Int8],
	"int16":      basicTypes[Int16],
	"int32":      basicTypes[Int32],
	"int64":      basicTypes[Int64],
	"uint":       basicTypes[Uint],
	"uint8":      basicTypes[Uint8],
	"uint16":     basicTypes[Uint16],
	"uint32":     basicTypes[Uint32],
	"uint64":     basicTypes[Uint64],
	"uintptr":    basicTypes[Uintptr],
	"float32":    basicTypes[Float32],
	"float64":    basicTypes[Float64],
	"complex64":  basicTypes[Complex64],
	"complex128": basicTypes[Complex128],
	"string":     basicTypes[String],
	"byte":       basicTypes[Byte],
	"rune":       basicTypes[Rune],
}

var basicTypes = []*Basic{
	Invalid: {Invalid, 0, "invalid type"},

	Bool:          {Bool, IsBoolean, "bool"},
	Int:           {Int, IsInteger, "int"},
	Int8:          {Int8, IsInteger, "int8"},
	Int16:         {Int16, IsInteger, "int16"},
	Int32:         {Int32, IsInteger, "int32"},
	Int64:         {Int64, IsInteger, "int64"},
	Uint:          {Uint, IsInteger | IsUnsigned, "uint"},
	Uint8:         {Uint8, IsInteger | IsUnsigned, "uint8"},
	Uint16:        {Uint16, IsInteger | IsUnsigned, "uint16"},
	Uint32:        {Uint32, IsInteger | IsUnsigned, "uint32"},
	Uint64:        {Uint64, IsInteger | IsUnsigned, "uint64"},
	Uintptr:       {Uintptr, IsInteger | IsUnsigned, "uintptr"},
	Float32:       {Float32, IsFloat, "float32"},
	Float64:       {Float64, IsFloat, "float64"},
	Complex64:     {Complex64, IsComplex, "complex64"},
	Complex128:    {Complex128, IsComplex, "complex128"},
	String:        {String, IsString, "string"},
	UnsafePointer: {UnsafePointer, 0, "Pointer"},

	UntypedBool:    {UntypedBool, IsBoolean | IsUntyped, "untyped bool"},
	UntypedInt:     {UntypedInt, IsInteger | IsUntyped, "untyped int"},
	UntypedRune:    {UntypedRune, IsInteger | IsUntyped, "untyped rune"},
	UntypedFloat:   {UntypedFloat, IsFloat | IsUntyped, "untyped float"},
	UntypedComplex: {UntypedComplex, IsComplex | IsUntyped, "untyped complex"},
	UntypedString:  {UntypedString, IsString | IsUntyped, "untyped string"},
	UntypedNil:     {UntypedNil, IsUntyped, "untyped nil"},

	{Byte, IsInteger | IsUnsigned, "byte"},
	{Rune, IsInteger, "rune"},
}

type Basic struct {
	kind BasicKind
	info BasicInfo
	name string
}

func (b *Basic) Kind() BasicKind {
	return b.kind
}

func (b *Basic) Name() string {
	return b.name
}

func (b *Basic) Underlying() Type {
	return b
}

func (b *Basic) String() string {
	return b.name
}
