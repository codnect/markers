package processor

import "github.com/spf13/cobra"

func AddCommand(cmd *cobra.Command) {
	rootCmd.AddCommand(cmd)
}
