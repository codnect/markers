package main

import (
	"encoding/json"
	"errors"
	"github.com/procyon-projects/marker/packages"
	"github.com/procyon-projects/marker/processor"
	"github.com/spf13/cobra"
	"log"
)

var listCmd = &cobra.Command{
	Use:   "list [pkg]",
	Short: "List marker package",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("pkg is required")
		}

		return listPackage(args[0])
	},
}

func init() {
	processor.AddCommand(listCmd)
}

func listPackage(pkg string) error {
	markerPackage, err := packages.GetMarkerPackage(pkg)
	if err != nil {
		return err
	}

	jsonText, _ := json.MarshalIndent(markerPackage, "", "\t")
	log.Println(string(jsonText))
	return nil
}
