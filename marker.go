package marker

type TargetLevel int

const (
	TypeLevel TargetLevel = 1 << iota
	ImportLevel
	FieldLevel
	FunctionLevel
	MethodLevel
	PackageLevel
)
