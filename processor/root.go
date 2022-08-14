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
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "marker",
	Short: "CLI Tool for marker processor and code generation",
	Long:  `CLI Tool for marker processor and code generation`,
}

func Execute() {
	if processorInfo == nil || generateCallback == nil || len(registryFunctions) == 0 {
		log.Fatal("processor could not be initialized properly")
	}

	cobra.CheckErr(rootCmd.Execute())
}
