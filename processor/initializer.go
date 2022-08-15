/*
Copyright © 2021 Marker Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package processor

import (
	"log"
	"sync"
)

type Processor struct {
	Name              string
	Version           string
	GenerateCallback  GenerateCallback
	ValidateCallback  ValidateCallback
	RegistryFunctions []RegistryFunction
}

var (
	processorMu   sync.RWMutex
	processorInfo *Processor
)

func Initialize(processor *Processor) {
	if processor == nil {
		log.Fatal("processor cannot be nil")
	}

	if processor.GenerateCallback == nil {
		log.Fatal("GenerateCallback cannot be nil")
	}

	if processor.ValidateCallback == nil {
		log.Fatal("ValidateCallback cannot be nil")
	}

	if len(processor.RegistryFunctions) == 0 {
		log.Fatal("there should be at least one registry function")
	}

	processorMu.Lock()
	defer processorMu.Unlock()
	processorInfo = processor
}
