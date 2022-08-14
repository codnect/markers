package processor

import (
	"github.com/procyon-projects/marker"
	"sync"
)

type RegistryFunction func(marker *marker.Registry) error

var (
	registryFunctionsMu sync.RWMutex
	registryFunctions   = make([]RegistryFunction, 0)
)

func register(function RegistryFunction) {
	registryFunctionsMu.Lock()
	defer registryFunctionsMu.Unlock()
	registryFunctions = append(registryFunctions, function)
}

func invokeRegistryFunctions(marker *marker.Registry) error {
	for _, registryFunction := range registryFunctions {
		if err := registryFunction(marker); err != nil {
			return err
		}
	}

	return nil
}
