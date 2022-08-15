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
	"github.com/spf13/cobra"
)

var validateArgs []string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate markers' syntax and arguments",
	Long:  `The validate command helps you validate markers' syntax and arguments'`,
	Run: func(cmd *cobra.Command, args []string) {
		/*var err error
		var dirs []string

		dirs, err = getPackageDirectories()

		if err != nil {
			log.Println(err)
		}

		if dirs == nil || len(dirs) == 0 {
			return
		}

		var packages []*marker.Package
		packages, err = marker.LoadPackages(dirs...)

		if err != nil {
			log.Println(err)
			return
		}

		registry := marker.NewRegistry()
		err = RegisterDefinitions(registry)
		params := map[string]any{
			"args": validateArgs,
		}

		if err != nil {
			log.Println(err)
			return
		}

		collector := marker.NewCollector(registry)
		err = validateMarkers(collector, packages, dirs)

		if err != nil {
			log.Println(err)
			return
		}*/
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringSliceVarP(&validateArgs, "args", "a", validateArgs, "extra arguments for marker processors (key-value separated by comma)")
}
