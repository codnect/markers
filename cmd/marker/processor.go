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
	"log"
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
	err        error
	processors = make(map[string]MarkerProcessor, 0)
)

// ProcessMarkers gets the import markers in the given directories.
// Then, it fetches marker processors and run them for validation.
func ValidateMarkers(collector *marker.Collector, pkgs []*marker.Package, dirs []string) error {
	collectMarkerProcessors(collector, pkgs, dirs)

	if err != nil {
		switch typedErr := err.(type) {
		case marker.ErrorList:
			printErrors(typedErr)
			return nil
		}
		return err
	}

	fetchPackages()

	if err != nil {
		return err
	}

	validate(dirs)

	return err
}

// ProcessMarkers gets the import markers in the given directories.
// Then, it fetches marker processors and run them for code generation.
func ProcessMarkers(collector *marker.Collector, pkgs []*marker.Package, dirs []string) error {
	collectMarkerProcessors(collector, pkgs, dirs)

	if err != nil {
		switch typedErr := err.(type) {
		case marker.ErrorList:
			printErrors(typedErr)
			return nil
		}
		return err
	}

	fetchPackages()

	if err != nil {
		return err
	}

	generateCode(dirs)

	return err
}

// collectImportMarkers collects marker processor metadata by parsing '+import' marker
func collectMarkerProcessors(collector *marker.Collector, pkgs []*marker.Package, dirs []string) {
	marker.EachFile(collector, pkgs, func(file *marker.File, fileErr error) {
		if err != nil {
			return
		}

		if fileErr != nil {
			err = fileErr
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
						fileErr = fmt.Errorf("conflict: PkgId with '%s' has got more than one version, versions: '%s' and '%s'",
							pkgId, processor.Version, version)
						break
					}

					command := importMarker.GetCommand()

					if command == "" {
						command = importMarker.Value
					}

					if processor.Command != command {
						fileErr = fmt.Errorf("conflict: PkgId with '%s' has got more than one command, commands: '%s' and '%s'",
							pkgId, processor.Command, command)
						break
					}

				}
			}

		}
	})
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

		err = exec.Command("go", "get", "-u", name).Run()

		if err != nil {
			log.Errorf("an error occurred while fetching '%s'", name)
			break
		}
	}
}

// generateCode runs the marker processors to generate code
func generateCode(dirs []string) {
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

	runProcessors(args)
}

// validate runs the marker processors to validate markers
func validate(dirs []string) {
	args := make([]string, 0)

	args = append(args, "validate")
	args = append(args, "--path")
	args = append(args, strings.Join(dirs, ","))

	if validateArgs != nil && len(validateArgs) != 0 {
		args = append(args, "--args")
		args = append(args, strings.Join(validateArgs, ","))
	}

	runProcessors(args)
}

// runProcessor runs processors by passing given args
func runProcessors(args []string) {
	for _, processor := range processors {
		cmd := exec.Command(processor.Command, args...)
		output, err := cmd.CombinedOutput()

		if err != nil {
			log.Errorf("An error occurred while running command '%s %s' : ", processor.Command, strings.Join(args, " "))
			log.Errorf(err.Error())
		}

		if output != nil {
			log.Errorf(string(output))
		}

		if err != nil || output != nil {
			log.Println()
		}

	}
}
