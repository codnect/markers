package processor

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: fmt.Sprintf("Print version"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s version %s\n", processorName, processorVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
