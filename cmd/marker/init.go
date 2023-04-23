package main

import (
	"encoding/json"
	"errors"
	"github.com/procyon-projects/markers/processor"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"path/filepath"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize new marker project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return initializeMarkerProject()
	},
}

func init() {
	processor.AddCommand(initCmd)
}

func initializeMarkerProject() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	markerJsonPath := filepath.FromSlash(path.Join(wd, "marker.json"))

	_, err = os.Stat(markerJsonPath)
	if err == nil {
		log.Println("marker project is already initialized")
		return nil
	}

	var markerJsonFile *os.File
	markerJsonFile, err = os.Create(markerJsonPath)

	if err != nil {
		return errors.New("marker project is not initialized")
	}

	defer markerJsonFile.Close()

	config := &processor.Config{
		Version: Version,
		Parameters: []processor.Parameter{
			{
				Name:  "OUTPUT_PATH",
				Value: "${MODULE_ROOT}/generated",
			},
		},
		Overrides: make([]processor.Override, 0),
	}

	jsonText, _ := json.MarshalIndent(config, "", "\t")
	_, err = markerJsonFile.Write(jsonText)

	if err != nil {
		return errors.New("marker project is not initialized")
	}

	log.Println("marker project is initialized")
	return nil
}
