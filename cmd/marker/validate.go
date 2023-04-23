package main

import (
	"github.com/procyon-projects/markers/processor"
	"os"
	"os/exec"
)

func Validate(ctx *processor.Context) {

}

func runValidate(processorName, processorPath, configFilePath string) error {
	args := make([]string, 0)
	args = append(args, "validate")
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
