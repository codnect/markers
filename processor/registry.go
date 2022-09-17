package processor

import "sync"

type RegistryFunction func(ctx *Context) error

var (
	registryFunctionMu sync.Mutex
	registryFunctions  = make([]RegistryFunction, 0)
)

func AddRegistryFunction(function RegistryFunction) {
	defer registryFunctionMu.Unlock()
	registryFunctionMu.Lock()

	if function == nil {
		return
	}

	registryFunctions = append(registryFunctions, function)
}

func invokeRegistryFunctions(ctx *Context) error {
	defer registryFunctionMu.Unlock()
	registryFunctionMu.Lock()

	for _, function := range registryFunctions {
		err := function(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
