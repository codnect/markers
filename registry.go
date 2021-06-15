package marker

type Registry struct {
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (registry *Registry) Register(definition *Definition) {

}

func (registry *Registry) Define(name string, level TargetLevel) {

}
