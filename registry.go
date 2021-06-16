package marker

import (
	"fmt"
	"strings"
	"sync"
)

type Definition struct {
	Name  string
	Level TargetLevel
}

func MakeDefinition(name string, level TargetLevel) (*Definition, error) {
	return &Definition{
		Name:  name,
		Level: level,
	}, nil
}

type Registry struct {
	definitionMap map[string]*Definition

	initOnce sync.Once
	mu       sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (registry *Registry) initialize() {
	registry.initOnce.Do(func() {

		if registry.definitionMap == nil {
			registry.definitionMap = make(map[string]*Definition)
		}

	})
}

func (registry *Registry) Register(name string, level TargetLevel) error {
	registry.initialize()

	def, err := MakeDefinition(name, level)

	if err != nil {
		return err
	}

	return registry.RegisterWithDefinition(def)
}

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

func (registry *Registry) Lookup(name string) *Definition {
	registry.initialize()

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	return nil
}

func splitMarker(marker string) (name string, candidateName string, options string) {
	marker = marker[1:]

	nameFieldParts := strings.SplitN(marker, "=", 2)

	if len(nameFieldParts) == 1 {
		return nameFieldParts[0], nameFieldParts[0], ""
	}

	candidateName = nameFieldParts[0]
	name = candidateName

	nameParts := strings.Split(name, ":")

	if len(nameParts) > 1 {
		name = strings.Join(nameParts[:len(nameParts)-1], ":")
	}

	return name, candidateName, nameFieldParts[1]
}
