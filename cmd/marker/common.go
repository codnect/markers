package main

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/marker/internal/cmd"
	"github.com/procyon-projects/marker/packages"
	"io"
	"os"
	"os/exec"
)

func generateModFile(moduleName string) error {
	command := exec.Command("go", "mod", "init", moduleName)
	command.Stdout, command.Stderr = os.Stdout, os.Stderr
	executor := cmd.GetCommandExecutor()
	_, err := executor.Execute(command)

	if err != nil {
		return fmt.Errorf("could not create go.mod file %s", moduleName)
	}

	return nil
}

func createMarkerPackageFolder(info *packages.PackageInfo) error {
	markerPackagePath := packages.MarkerPackagePathFromPackageInfo(info)

	if _, err := os.Stat(markerPackagePath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(markerPackagePath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("marker path is not created for pkg %s", info.Name())
		}
	}

	return nil
}

func isMarkerProcessorPackage(info *packages.PackageInfo) bool {
	if _, err := os.Stat(packages.MarkerProcessorYamlPath(info)); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
func createFolder(folderPath string) error {
	if _, err := os.Stat(folderPath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("folder %s is not created", folderPath)
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	var err error
	var srcfd *os.File
	var dstfd *os.File
	var srcinfo os.FileInfo

	if srcfd, err = os.Open(src); err != nil {
		return err
	}
	defer srcfd.Close()

	if dstfd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstfd.Close()

	if _, err = io.Copy(dstfd, srcfd); err != nil {
		return err
	}

	if srcinfo, err = os.Stat(src); err != nil {
		return err
	}

	return os.Chmod(dst, srcinfo.Mode())
}
