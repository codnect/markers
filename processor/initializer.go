package processor

import (
	"log"
	"sync"
)

type Processor struct {
	Name    string
	Version string
}

var (
	processorMu   sync.RWMutex
	processorInfo *Processor
)

func Initialize(processor *Processor, generateCallback GenerateCallback, registryFunctions ...RegistryFunction) {
	if processor == nil {
		log.Fatal("processor parameter cannot be nil")
	}

	if generateCallback == nil {
		log.Fatal("generateCallback cannot be nil")
	}

	if len(registryFunctions) == 0 {
		log.Fatal("there should be at least one registry function")
	}

	processorMu.Lock()
	defer processorMu.Unlock()
	processorInfo = processor

	setGenerateCallback(generateCallback)

	for _, registryFunction := range registryFunctions {
		register(registryFunction)
	}
}
