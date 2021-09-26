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
)

var validateArgs []string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate markers' syntax and arguments",
	Long:  `The validate command helps you validate markers' syntax and arguments'`,
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
		return ValidateMarkers(collector, packages, dirs)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringSliceVarP(&validateArgs, "args", "a", validateArgs, "extra arguments for marker processors (key-value separated by comma)")
}
