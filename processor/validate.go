package processor

import (
	"github.com/spf13/cobra"
)

var validateArgs []string

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate markers' syntax and arguments",
	Long:  `The validate command helps you validate markers' syntax and arguments'`,
	Run: func(cmd *cobra.Command, args []string) {
		/*var err error
		var dirs []string

		dirs, err = getPackageDirectories()

		if err != nil {
			log.Println(err)
		}

		if dirs == nil || len(dirs) == 0 {
			return
		}

		var packages []*marker.Package
		packages, err = marker.LoadPackages(dirs...)

		if err != nil {
			log.Println(err)
			return
		}

		registry := marker.NewRegistry()
		err = RegisterDefinitions(registry)
		params := map[string]any{
			"args": validateArgs,
		}

		if err != nil {
			log.Println(err)
			return
		}

		collector := marker.NewCollector(registry)
		err = validateMarkers(collector, packages, dirs)

		if err != nil {
			log.Println(err)
			return
		}*/
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringSliceVarP(&validateArgs, "args", "a", validateArgs, "extra arguments for marker processors (key-value separated by comma)")
}
