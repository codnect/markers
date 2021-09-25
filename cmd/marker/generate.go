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
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var outputPath string
var options []string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go files by processing markers",
	Long:  `The generate command helps your code generation process by running marker processors`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var dirs []string

		dirs, err = getPackageDirectories()

		if err != nil {
			return err
		}

		if dirs == nil || len(dirs) == 0 {
			return nil
		}

		var packages []*marker.Package
		packages, err = marker.LoadPackages(dirs...)

		if err != nil {
			return err
		}

		registry := marker.NewRegistry()
		err = RegisterDefinitions(registry)

		if err != nil {
			return err
		}

		collector := marker.NewCollector(registry)
		return ProcessMarkers(collector, packages, dirs)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output path")
	err := generateCmd.MarkFlagRequired("output")

	if err != nil {
		panic(err)
	}

	generateCmd.Flags().StringSliceVarP(&options, "args", "a", options, "extra arguments for marker processors (key-value separated by comma)")
}

// getPackageDirectories finds the go module directory and returns
// the package directories.
func getPackageDirectories() ([]string, error) {
	var err error
	var modDir string
	modDir, err = marker.GoModDir()

	if err != nil {
		return nil, err
	}

	var dirs []string
	dirs, err = findDirectoriesWithGoFiles(modDir)

	if err != nil {
		return nil, err
	}

	return dirs, nil
}

// findDirectoriesWithGoFiles returns the go directories with go files.
// if not, an error might occur while loading packages.
func findDirectoriesWithGoFiles(root string) ([]string, error) {
	dirMap := make(map[string]bool, 0)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip hidden directories
		if strings.HasPrefix(path, ".") && !strings.HasPrefix(path, "./") {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match("*.go", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			dirMap[filepath.Dir(path)] = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var dirs []string

	for dir, _ := range dirMap {
		dirs = append(dirs, dir)
	}

	return dirs, nil
}
