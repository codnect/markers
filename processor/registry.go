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
	"github.com/procyon-projects/marker"
)

type RegistryFunction func(marker *marker.Registry) error

func invokeRegistryFunctions(marker *marker.Registry) error {
	for _, registryFunction := range processorInfo.RegistryFunctions {
		if err := registryFunction(marker); err != nil {
			return err
		}
	}

	return nil
}
