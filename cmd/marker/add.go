package main

import (
	"github.com/procyon-projects/marker/processor"
	"github.com/spf13/cobra"
)

var processorName string

var addCmd = &cobra.Command{
	Use:   "add [processor]",
	Short: "Add a marker processor",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	processor.AddCommand(addCmd)
}
