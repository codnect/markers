package main

import (
	"github.com/procyon-projects/marker/packages"
	"github.com/procyon-projects/marker/processor"
	"github.com/procyon-projects/marker/visitor"
	"os"
	"os/exec"
)

func Generate(ctx *processor.Context) {
	pkg, _ := ctx.LoadResult().Lookup("github.com/procyon-projects/marker/test/package1")
	err := visitor.EachFile(ctx.Collector(), []*packages.Package{pkg}, func(file *visitor.File, err error) error {
		if file.NumImportMarkers() == 0 {
			return nil
		}

		file.Structs().At(0).Markers()

		return err
	})

	if err != nil {
		return
	}
}

func runGenerate(processorName, processorPath, configFilePath string) error {
	args := make([]string, 0)
	args = append(args, "generate")
	args = append(args, "-f")
	args = append(args, configFilePath)
	cmd := exec.Command(processorName, args...)
	cmd.Dir = processorPath
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
