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
package main

import (
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"github.com/procyon-projects/marker/processor"
	"log"
)

func init() {
	processor.Initialize(&processor.Processor{
		Name:    AppName,
		Version: AppVersion,
		GenerateCallback: func(collector *marker.Collector, loadResult *packages.LoadResult, params map[string]any) error {
			return nil
		},
		RegistryFunctions: []processor.RegistryFunction{},
	})
}

func main() {
	log.SetFlags(0)
	processor.Execute()
}
