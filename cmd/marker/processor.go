package main

import (
	"errors"
	"github.com/procyon-projects/markers/packages"
	"github.com/procyon-projects/markers/processor"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var commonFileContent = `package cmd

var (
	Package = "$MODULE_NAME"
	Version = "[your version goes here]"
)

`

var mainFileContent = `// Code generated by marker; DO NOT EDIT.
package main

import (
	"github.com/procyon-projects/markers/cmd"
	"github.com/procyon-projects/markers/processor"
	"log"
)

func init() {
	processor.Initialize(cmd.Package, "$PROCESSOR_NAME", cmd.Version)
	processor.SetGenerateCommandCallback(Generate)
	processor.SetValidateCommandCallback(Validate)
}

func main() {
	log.SetFlags(0)
	processor.Execute()
}
`

var generateFileContent = `package main

import (
	"github.com/procyon-projects/markers/processor"
)

func Generate(ctx *processor.Context) {

}
`

var validateFileContent = `package main

import (
	"github.com/procyon-projects/markers/processor"
)

func Validate(ctx *processor.Context) {

}
`

var moduleName string

var processorCmd = &cobra.Command{
	Use:   "processor",
	Short: "Create your own marker processor package",
}

var processorInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize marker processor package",
	RunE: func(cmd *cobra.Command, args []string) error {
		return initializeMarkerProcessorPackage()
	},
}

var processorAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add marker processor",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("name is required")
		}
		return addProcessorCommand(args[0])
	},
}

func init() {
	processorInitCmd.Flags().StringVarP(&moduleName, "module", "m", "", "module name")
	processorCmd.AddCommand(processorInitCmd)
	processorCmd.AddCommand(processorAddCmd)
	processor.AddCommand(processorCmd)
}

func initializeMarkerProcessorPackage() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if moduleName == "" {
		var loadResult *packages.LoadResult
		loadResult, err = packages.LoadPackages()
		if err != nil || len(loadResult.Packages()) == 0 {
			return errors.New("go.mod not found, module name is required")
		}
		pkg := loadResult.Packages()[0]
		moduleName = pkg.PkgPath
	} else {
		err = generateModFile(moduleName)
		if err != nil {
			return err
		}
	}

	generatedFolder := filepath.FromSlash(path.Join(wd, "generated"))
	err = createFolder(generatedFolder)
	if err != nil {
		return err
	}

	cmdFolder := filepath.FromSlash(path.Join(wd, "cmd"))
	err = createFolder(cmdFolder)
	if err != nil {
		return err
	}

	commonFilePath := filepath.FromSlash(path.Join(cmdFolder, "common.go"))
	_, err = os.Stat(commonFilePath)
	if err != nil {
		var commonGoFile *os.File
		commonGoFile, err = os.Create(commonFilePath)

		if err != nil {
			return errors.New("common.go is not created")
		}

		defer commonGoFile.Close()
		commonFileContent = strings.ReplaceAll(commonFileContent, "$MODULE_NAME", moduleName)
		_, err = commonGoFile.Write([]byte(commonFileContent))

		if err != nil {
			return errors.New("initializing processor package failed")
		}
	}

	yamlPath := filepath.FromSlash(path.Join(wd, "marker.processors.yaml"))
	_, err = os.Stat(yamlPath)
	if err == nil {
		log.Println("marker processor package is already initialized")
		return nil
	}

	var yamlFile *os.File
	yamlFile, err = os.Create(yamlPath)

	if err != nil {
		return errors.New("initializing processor package failed")
	}

	defer yamlFile.Close()
	yamlText, _ := yaml.Marshal([]byte{})
	_, err = yamlFile.Write(yamlText)

	if err != nil {
		return errors.New("initializing processor package failed")
	}

	log.Println("marker processor package is initialized")
	return nil
}

func addProcessorCommand(processorName string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	processorName = strings.ToLower(processorName)

	yamlPath := filepath.FromSlash(path.Join(wd, "marker.processors.yaml"))
	_, err = os.Stat(yamlPath)
	if err != nil {
		return errors.New("marker processor package is not initialized")
	}

	cmdFolder := filepath.FromSlash(path.Join(wd, "cmd"))
	_, err = os.Stat(cmdFolder)
	if err != nil {
		return errors.New("marker processor package is not initialized")
	}

	generatedFolder := filepath.FromSlash(path.Join(wd, "generated"))
	_, err = os.Stat(generatedFolder)
	if err != nil {
		return errors.New("marker processor package is not initialized")
	}

	processorFolder := filepath.FromSlash(path.Join(cmdFolder, processorName))
	_, err = os.Stat(processorFolder)
	if err != nil {
		return errors.New("processor already exists")
	}

	err = createFolder(processorFolder)
	if err != nil {
		return err
	}

	mainFilePath := filepath.FromSlash(path.Join(cmdFolder, "main.go"))
	_, err = os.Stat(mainFilePath)
	if err != nil {
		var mainGoFile *os.File
		mainGoFile, err = os.Create(mainFilePath)

		if err != nil {
			return errors.New("main.go is not created")
		}

		defer mainGoFile.Close()
		mainFileContent = strings.ReplaceAll(mainFileContent, "$PROCESSOR_NAME", processorName)
		_, err = mainGoFile.Write([]byte(mainFileContent))

		if err != nil {
			return errors.New("initializing processor package failed")
		}
	}

	validateFilePath := filepath.FromSlash(path.Join(processorFolder, "validate.go"))
	_, err = os.Stat(validateFilePath)
	if err != nil {
		var validateFile *os.File
		validateFile, err = os.Create(validateFilePath)

		if err != nil {
			return errors.New("validate.go is not created")
		}

		defer validateFile.Close()
		_, err = validateFile.Write([]byte(validateFileContent))

		if err != nil {
			return errors.New("creating processor failed")
		}
	}

	generateFilePath := filepath.FromSlash(path.Join(processorFolder, "generate.go"))
	_, err = os.Stat(generateFilePath)
	if err != nil {
		var generateFile *os.File
		generateFile, err = os.Create(generateFilePath)

		if err != nil {
			return errors.New("validate.go is not created")
		}

		defer generateFile.Close()
		_, err = generateFile.Write([]byte(generateFileContent))

		if err != nil {
			return errors.New("creating processor failed")
		}
	}

	return nil
}
