package main

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/marker/packages"
	"github.com/procyon-projects/marker/processor"
	"github.com/spf13/cobra"
	"os"
)

var downloadCmd = &cobra.Command{
	Use:   "download [pkg]",
	Short: "Download marker package",
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
	pkgInfo, err := packages.GetPackageInfo(pkg)
	if err != nil {
		return err
	}

	if err = createMarkerPackageFolder(pkgInfo); err != nil {
		return err
	}

	if err = packages.InstallPackage(pkgInfo); err != nil {
		return err
	}

	if !isMarkerProcessorPackage(pkgInfo) {
		os.RemoveAll(packages.MarkerPackagePathFromPackageInfo(pkgInfo))
		return fmt.Errorf("'%s' is not valid marker processor package", pkgInfo.Name())
	}

	err = copyFile(packages.MarkerProcessorYamlPath(pkgInfo), packages.MarkerPackageYamlPath(pkgInfo))
	if err != nil {
		return err
	}

	fmt.Printf("%s downloaded successfully", pkgInfo.Name())
	return nil
}
