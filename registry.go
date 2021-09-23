package marker

import (
	"fmt"
	"sync"
)

// Registry keeps the registered marker definitions.
type Registry struct {
	definitionMap map[string]*Definition

	initOnce sync.Once
	mu       sync.RWMutex
}

// NewRegistry returns a new registry to register the markers.
func NewRegistry() *Registry {
	return &Registry{}
}

// initialize initializes the registry when needed.
func (registry *Registry) initialize() {
	registry.initOnce.Do(func() {

		if registry.definitionMap == nil {
			registry.definitionMap = make(map[string]*Definition)
		}

		registry.definitionMap["import"], _ = MakeDefinition("import", ImportLevel, &ImportMarker{}, true)
	})

}

// Register registers a new marker with the given name, target level, and output type.
func (registry *Registry) Register(name string, level TargetLevel, output interface{}, useValueSyntax ...bool) error {
	registry.initialize()

	def, err := MakeDefinition(name, level, output, useValueSyntax...)

	if err != nil {
		return err
	}

	return registry.RegisterWithDefinition(def)
}

// RegisterWithDefinition registers a new marker with the given definition.
func (registry *Registry) RegisterWithDefinition(definition *Definition) error {
	registry.initialize()

	registry.mu.Lock()
	defer registry.mu.Unlock()

	if definition.Level == 0 {
		return fmt.Errorf("specify target levels for the definition : %v", definition.Name)
	}

	if _, ok := registry.definitionMap[definition.Name]; ok {
		return fmt.Errorf("there is already registered definition : %v", definition.Name)
	}

	registry.definitionMap[definition.Name] = definition

	return nil
}

// Lookup fetches the definition corresponding to the given name.
func (registry *Registry) Lookup(name string) *Definition {
	registry.initialize()

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	name, anonymousName, _ := splitMarker(name)

	if def, exists := registry.definitionMap[anonymousName]; exists {
		return def
	}

	return registry.definitionMap[name]
}
