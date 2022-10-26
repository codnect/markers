package processor

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"github.com/spf13/cobra"
	"path"
)

var configFilePath string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Go files by processing markers",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if configFilePath == "" {
			configFilePath, err = getConfigFilePath()
			if err != nil {
				return err
			}
		}

		var config *Config
		config, err = getConfig(configFilePath)

		if err != nil {
			return fmt.Errorf("%s not found", path.Join(configFilePath, "marker.json"))
		}

		// TODO check marker package details
		_, err = packages.GetMarkerPackage(fmt.Sprintf("%s@%s", packageName, processorVersion))

		if err != nil {
			return err
		}

		modDir, _ := packages.GoModDir()
		ctx := &Context{
			configFilePath: configFilePath,
			config:         *config,
			packageId:      packageName,
			version:        processorVersion,
			goModuleDir:    modDir,
			errors:         make([]error, 0),
			values:         map[string]any{},
			args:           args,
		}

		var dirs []string
		dirs, err = PackageDirectories()

		if err != nil {
			return errors.New("go.module not found")
		}

		if dirs == nil || len(dirs) == 0 {
			return nil
		}

		ctx.dirs = dirs

		var loadResult *packages.LoadResult
		loadResult, err = packages.LoadPackages(dirs...)

		if err != nil {
			return errors.New("packages could not be loaded")
		}

		registry := markers.NewRegistry()
		ctx.loadResult = loadResult
		ctx.registry = registry

		err = invokeRegistryFunctions(ctx)
		if err != nil {
			return err
		}

		collector := markers.NewCollector(registry)
		ctx.collector = collector

		generateCallback := getGenerateCommandCallback()
		if generateCallback != nil {
			generateCallback(ctx)
		}

		return nil
	},
}

func init() {
	generateCmd.Flags().StringVarP(&configFilePath, "file", "f", "", "config file path")
	rootCmd.AddCommand(generateCmd)
}
