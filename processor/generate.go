package processor

import (
	"errors"
	"github.com/procyon-projects/marker"
	"github.com/procyon-projects/marker/packages"
	"github.com/spf13/cobra"
	"log"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate Go files by processing markers",
	RunE: func(cmd *cobra.Command, args []string) error {
		dirs, err := getPackageDirectories()

		if err != nil {
			return errors.New("go.module not found")
		}

		if dirs == nil || len(dirs) == 0 {
			return nil
		}

		var loadResult *packages.LoadResult
		loadResult, err = packages.LoadPackages(dirs...)

		if err != nil {
			log.Println(err)
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

		generateCallback := getGenerateCommandCallback()
		if generateCallback != nil {
			err = generateCallback(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
