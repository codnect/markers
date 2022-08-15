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
	"github.com/procyon-projects/marker/processor"
	"github.com/spf13/cobra"
)

var module string

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a marker processor project",
	Long:  `The init command lets you create a marker processor project.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("marker name is required")
		}

		return nil
	},
}

func init() {
	processor.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&module, "module", "m", "", "Module Name (required)")

	err := initCmd.MarkFlagRequired("module")

	if err != nil {
		panic(err)
	}
}
