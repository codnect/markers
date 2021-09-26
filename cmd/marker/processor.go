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
	"fmt"
	"github.com/procyon-projects/marker"
	"os/exec"
	"strings"
)

type MarkerProcessor struct {
	Module  string
	Version string
	Command string
}

// Register your marker definitions.
func RegisterDefinitions(registry *marker.Registry) error {
	return nil
}

var (
	processErr error
	processors = make(map[string]MarkerProcessor, 0)
)

// ProcessMarkers gets the import markers in the given directories.
// Then, it fetches marker processors and run them.
func ProcessMarkers(collector *marker.Collector, pkgs []*marker.Package, dirs []string) error {
	marker.EachFile(collector, pkgs, func(file *marker.File, markerErr error) {
		if processErr != nil {
			return
		}

		if markerErr != nil {
			processErr = markerErr
			return
		}

		if file.ImportMarkers == nil || len(file.ImportMarkers) == 0 {
			return
		}

		for _, markerValues := range file.ImportMarkers {
			importMarkers := markerValues[marker.ImportMarkerName]

			if importMarkers == nil || len(importMarkers) == 0 {
				continue
			}

			for _, value := range importMarkers {
				importMarker := value.(marker.ImportMarker)
				pkgId := importMarker.GetPkgId()

				processor, ok := processors[pkgId]

				if !ok {
					command := importMarker.GetCommand()

					if command == "" {
						command = importMarker.Value
					}

					processors[pkgId] = MarkerProcessor{
						Module:  pkgId,
						Version: importMarker.GetPkgVersion(),
						Command: command,
					}
				} else {
					version := importMarker.GetPkgVersion()

					if processor.Version != version {
						processErr = fmt.Errorf("conflict: PkgId with '%s' has got more than one version, versions: '%s' and '%s'",
							pkgId, processor.Version, version)
						break
					}

					command := importMarker.GetCommand()

					if command == "" {
						command = importMarker.Value
					}

					if processor.Command != command {
						processErr = fmt.Errorf("conflict: PkgId with '%s' has got more than one command, commands: '%s' and '%s'",
							pkgId, processor.Command, command)
						break
					}

				}
			}

		}
	})

	if processErr != nil {
		return processErr
	}

	fetchPackages()

	if processErr != nil {
		return processErr
	}

	runProcessors(dirs)

	return processErr
}

// runProcessors fetches the marker processors by making use of '+import' marker metadata.
func fetchPackages() {
	for _, processor := range processors {
		name := fmt.Sprintf("%s/...", processor.Module)

		if processor.Version != "" {
			name = fmt.Sprintf("%s@%s", name, processor.Version)
			fmt.Printf("Fetching %s@%s...\n", processor.Module, processor.Version)
		} else {
			fmt.Printf("Fetching %s...\n", processor.Module)
		}

		processErr = exec.Command("go", "get", "-u", name).Run()

		if processErr != nil {
			processErr = fmt.Errorf("an error occurred while fetching '%s'", name)
			break
		}
	}
}

// runProcessors runs the marker processors by making use of '+import' marker metadata.
func runProcessors(dirs []string) {
	args := make([]string, 0)

	args = append(args, "generate")
	args = append(args, "--output")
	args = append(args, outputPath)
	args = append(args, "--path")
	args = append(args, strings.Join(dirs, ","))

	if options != nil && len(options) != 0 {
		args = append(args, "--args")
		args = append(args, strings.Join(options, ","))
	}

	for _, processor := range processors {
		processErr = exec.Command(processor.Command, args...).Run()

		if processErr != nil {
			break
		}
	}
}
