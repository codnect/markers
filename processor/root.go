package processor

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "marker",
	Short: "CLI Tool for marker",
}

func Execute() {
	rootCmd.Execute()
}
