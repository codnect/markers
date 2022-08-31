package processor

/*
var outputPath string
var options []string
var packageName string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate Go files by processing sources",
	Run: func(cmd *cobra.Command, args []string) {
		dirs, err := getPackageDirectories()

		if err != nil {
			log.Println(err)
			return
		}

		if dirs == nil || len(dirs) == 0 {
			return
		}

		var loadResult *packages.LoadResult
		loadResult, err = packages.LoadPackages(dirs...)

		if err != nil {
			log.Println(err)
		}

		registry := marker.NewRegistry()
		err = invokeRegistryFunctions(registry)

		if err != nil {
			log.Println(err)
			return
		}

		collector := marker.NewCollector(registry)
		params := map[string]any{
			"directories": dirs,
			"output":      outputPath,
			"package":     packageName,
			"args":        options,
		}

		err = processorInfo.GenerateCallback(collector, loadResult, params)

		if err != nil {
			log.Println(err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "", "output path")
	err := generateCmd.MarkFlagRequired("output")

	if err != nil {
		panic(err)
	}

	generateCmd.Flags().StringVarP(&packageName, "package", "p", "generated", "package name")
	generateCmd.Flags().StringSliceVarP(&options, "args", "a", options, "extra arguments for marker processors (key-value separated by comma)")
}
*/
