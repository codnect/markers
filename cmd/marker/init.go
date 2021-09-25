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
	"errors"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [name] [version]",
	Short: "Initialize a marker processor",
	Long:  `The init command lets you create a new marker processor.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return errors.New("processor name and version are required")
		}

		if len(args) == 1 {
			return errors.New("processor version is missing")
		}

		if len(args) > 2 {
			return errors.New("too many arguments")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
