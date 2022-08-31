package processor

import (
	"errors"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "validate marker syntax and arguments",
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		var dirs []string

		dirs, err = getPackageDirectories()

		if err != nil {
			return errors.New("go.module not found")
		}

		if dirs == nil || len(dirs) == 0 {
			return nil
		}

		var loadResult *packages.LoadResult
		loadResult, err = packages.LoadPackages(dirs...)

		if err != nil {
			return errors.New("packages could not be loaded")
		}

		registry := marker.NewRegistry()
		ctx := &Context{
			dirs:       dirs,
			loadResult: loadResult,
			registry:   registry,
		}

		err = invokeRegistryFunctions(ctx)
		if err != nil {
			return err
		}

		collector := marker.NewCollector(registry)
		ctx.collector = collector

		validateCallback := getValidateCommandCallback()
		if validateCallback != nil {
			err = validateCallback(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
