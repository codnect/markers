package marker

import (
	"errors"
	"fmt"
	"sync"
)

type DefinitionMap map[string]*Definition

// Registry keeps the registered marker definitions.
type Registry struct {
	packageMap map[string]DefinitionMap
	mu         sync.RWMutex
}

// NewRegistry returns a new registry to register the markers.
func NewRegistry() *Registry {
	registry := &Registry{}
	registry.initialize()
	return registry
}

// initialize initializes the registry when needed.
func (registry *Registry) initialize() {
	if registry.packageMap == nil {
		registry.packageMap = make(map[string]DefinitionMap)
		registry.packageMap[""] = make(DefinitionMap)
	}
	registry.packageMap[""][ImportMarkerName], _ = MakeDefinition(ImportMarkerName, "", PackageLevel, &ImportMarker{})
	registry.packageMap[""][OverrideMarkerName], _ = MakeDefinition(OverrideMarkerName, "", StructMethodLevel, &OverrideMarker{})
	registry.packageMap[""][DeprecatedMarkerName], _ = MakeDefinition(DeprecatedMarkerName, "", TypeLevel|MethodLevel|FieldLevel|FunctionLevel, &DeprecatedMarker{})
}

// Register registers a new marker with the given name, target level, and output type.
func (registry *Registry) Register(name, pkg string, level TargetLevel, output any) error {
	def, err := MakeDefinition(name, pkg, level, output)

	if err != nil {
		return err
	}

	return registry.RegisterWithDefinition(def)
}

// RegisterWithDefinition registers a new marker with the given definition.
func (registry *Registry) RegisterWithDefinition(definition *Definition) error {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	if definition == nil {
		return errors.New("definition cannot be nil")
	}

	err := definition.validate()

	if err != nil {
		return err
	}

	definitionMap, exists := registry.packageMap[definition.Package]
	if !exists {
		definitionMap = make(DefinitionMap)
		registry.packageMap[definition.Package] = definitionMap
	}

	if _, ok := definitionMap[definition.Name]; ok {
		return fmt.Errorf("there is already registered definition : %v", definition.Name)
	}

	registry.packageMap[definition.Package][definition.Name] = definition
	return nil
}

// Lookup fetches the definition corresponding to the given name and pkgId.
func (registry *Registry) Lookup(name, pkg string, targetLevel TargetLevel) (*Definition, bool) {
	registry.initialize()

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	definitionMap, exists := registry.packageMap[pkg]
	if !exists {
		return nil, false
	}

	if definition, ok := definitionMap[name]; ok && definition.TargetLevel&targetLevel == targetLevel {
		return definition, true
	}

	return nil, false
}
