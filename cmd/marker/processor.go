package main

/*
var module string

var initCmd = &cobra.Command{
	Use:   "processor [name]",
	Short: "Create marker processor package",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("marker name is required")
		}

		return nil
	},
}

func init() {
	processor.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&module, "module", "m", "", "Module Name (required)")

	err := initCmd.MarkFlagRequired("module")

	if err != nil {
		panic(err)
	}
}

type MarkerProcessor struct {
	Module  string
	Version string
	Command string
}

var (
	processors       = make(map[string]MarkerProcessor, 0)
	validationErrors []error
)

// ProcessMarkers gets the import markers in the given directories.
// Then, it fetches marker processors and run them for code generation.
func ProcessMarkers(collector *marker.Collector, loadResult *packages.LoadResult, dirs []string) error {
	err := collectMarkers(collector, loadResult.GetPackages())

	if validationErrors != nil {
		switch typedErr := err.(type) {
		case marker.ErrorList:
			printErrors(typedErr)
			return nil
		}
		return err
	}

	err = fetchPackages()

	if err != nil {
		return err
	}

	generateCode(dirs)

	return err
}

// CollectMarkers collects markers by scanning metadata
func collectMarkers(collector *marker.Collector, pkgs []*packages.Package) error {
	visitor.EachFile(collector, pkgs, func(file *visitor.File, fileErr error) {
		if fileErr != nil {
			validationErrors = append(validationErrors, fileErr)
			return
		}

		if file.NumImportMarkers() == 0 {
			return
		}

		for _, importMarker := range file.ImportMarkers() {
			pkgId := importMarker.GetPkgId()

			_, ok := processors[pkgId]

			if !ok {
				command := importMarker.GetCommand()

				if command == "" {
					command = importMarker.Value
				}

				processors[pkgId] = MarkerProcessor{
					Module:  pkgId,
					Version: importMarker.GetPkgVersion(),
					Command: command,
				}
			}
		}
	})
	return marker.NewErrorList(validationErrors)
}

// ProcessMarkers gets the import markers in the given directories.
// Then, it fetches marker processors and run them for validation.
func validateMarkers(collector *marker.Collector, pkgs []*packages.Package, dirs []string) error {
	err := collectMarkers(collector, pkgs)

	if err != nil {
		switch typedErr := err.(type) {
		case marker.ErrorList:
			printErrors(typedErr)
			return nil
		}
		return err
	}

	err = fetchPackages()

	if err != nil {
		return err
	}

	validate(dirs)

	return err
}

// runProcessors fetches the marker processors by making use of '+import' marker metadata.
func fetchPackages() error {
	for _, processor := range processors {
		name := fmt.Sprintf("%s/...", processor.Module)

		if processor.Version != "" {
			name = fmt.Sprintf("%s@%s", name, processor.Version)
			fmt.Printf("Fetching %s@%s...\n", processor.Module, processor.Version)
		} else {
			fmt.Printf("Fetching %s...\n", processor.Module)
		}

		err := exec.Command("go", "get", "-u", name).Run()

		if err != nil {
			return fmt.Errorf("an error occurred while fetching '%s'", name)
		}
	}

	return nil
}

// generateCode runs the marker processors to generate code
func generateCode(dirs []string) {
	args := make([]string, 0)

	args = append(args, "generate")
	args = append(args, "--output")
	args = append(args, outputPath)
	args = append(args, "--path")
	args = append(args, strings.Join(dirs, ","))

	if options != nil && len(options) != 0 {
		args = append(args, "--args")
		args = append(args, strings.Join(options, ","))
	}

	runProcessors(args)
}

// validate runs the marker processors to validate markers
func validate(dirs []string) {
	args := make([]string, 0)

	args = append(args, "validate")
	args = append(args, "--path")
	args = append(args, strings.Join(dirs, ","))

	if validateArgs != nil && len(validateArgs) != 0 {
		args = append(args, "--args")
		args = append(args, strings.Join(validateArgs, ","))
	}

	runProcessors(args)
}

// runProcessor runs processors by passing given args
func runProcessors(args []string) {
	for _, processor := range processors {
		cmd := exec.Command(processor.Command, args...)
		output, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("An error occurred while running command '%s %s' : ", processor.Command, strings.Join(args, " "))
			log.Fatalf(err.Error())
		}

		if output != nil {
			log.Printf(string(output))
		}

		if err != nil || output != nil {
			log.Println()
		}

	}
}
*/
