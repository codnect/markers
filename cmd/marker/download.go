package main

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/marker/processor"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

var downloadCmd = &cobra.Command{
	Use:   "download [pkg]",
	Short: "download marker package",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("pkg is required")
		}

		return downloadPackage(args[0])
	},
}

func init() {
	processor.AddCommand(downloadCmd)
}

func downloadPackage(pkg string) error {
	pkgParts := strings.SplitN(pkg, "@", 2)

	pkgName := pkg
	pkgVersion := "latest"

	if len(pkgParts) == 2 {
		pkgName = pkgParts[0]
		pkgVersion = pkgParts[1]
	}

	pkgInfo, err := getPackageInfo(fmt.Sprintf("%s@%s", pkgName, pkgVersion))
	if err != nil {
		switch typedErr := err.(type) {
		case *exec.ExitError:
			return errors.New(string(typedErr.Stderr))
		}
		return err
	}

	if err = createMarkerPackageFolder(pkgInfo); err != nil {
		return err
	}

	if err = installPackage(pkgInfo); err != nil {
		return err
	}

	if !isMarkerProcessorPackage(pkgInfo) {
		os.RemoveAll(getMarkerPackagePath(pkgInfo))
		return fmt.Errorf("'%s' is not valid marker processor package", pkgInfo.Name())
	}

	err = copyFile(getModuleMarkerProcessorYamlPath(pkgInfo), getMarkerPackageYamlPath(pkgInfo))
	if err != nil {
		return err
	}

	fmt.Printf("%s downloaded successfully", pkgInfo.Name())
	return nil
}
