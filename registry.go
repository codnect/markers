package marker

import (
	"fmt"
	"strings"
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

		registry.definitionMap[ImportMarkerName+"#"], _ = MakeDefinition(ImportMarkerName, "", ImportLevel, &ImportMarker{})
	})

}

// Register registers a new marker with the given name, target level, and output type.
func (registry *Registry) Register(name string, pkgId string, level TargetLevel, output interface{}, useValueSyntax ...bool) error {
	registry.initialize()

	def, err := MakeDefinition(name, pkgId, level, output)

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

	if definition.Level&ImportLevel == ImportLevel {
		return fmt.Errorf("level is not valid for %v, import level cannot be used", definition.Name)
	}

	nameParts := strings.Split(definition.Name, ":")
	name := nameParts[0]

	if _, ok := registry.definitionMap[name]; ok {
		return fmt.Errorf("reserved names such as 'import' cannot be used: %v", definition.Name)
	}

	if _, ok := registry.definitionMap[definition.Name+"#"+definition.PkgId]; ok {
		return fmt.Errorf("there is already registered definition : %v", definition.Name)
	}

	registry.definitionMap[definition.Name+"#"+definition.PkgId] = definition

	return nil
}

// Lookup fetches the definition corresponding to the given name and pkgId.
func (registry *Registry) Lookup(name string, pkgId string) *Definition {
	registry.initialize()

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	name, anonymousName, _ := splitMarker(name)
	// for syntax-free markers
	name = strings.Split(name, " ")[0]

	if def, exists := registry.definitionMap[anonymousName+"#"+pkgId]; exists {
		return def
	}

	return registry.definitionMap[name+"#"+pkgId]
}
