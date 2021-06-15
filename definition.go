package marker

type Definition struct {
	Name  string
	Level TargetLevel
}

func MakeDefinition(name string, level TargetLevel) *Definition {
	return &Definition{
		Name:  name,
		Level: level,
	}
}
