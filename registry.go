package marker

type Registry struct {
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (registry *Registry) Register(name string, level TargetLevel) {

}

func (registry *Registry) RegisterWithDefinition(definition *Definition) {

}
