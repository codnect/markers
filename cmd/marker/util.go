package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func goPath() string {
	goArgs := []string{"env", "GOPATH"}
	cmd := exec.Command("go", goArgs...)

	stdout, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.Trim(string(stdout), "\n ")
}

func getPackageInfo(path string) (*PackageInfo, error) {
	goArgs := []string{"list", "-m", "-versions", "-json", path}
	cmd := exec.Command("go", goArgs...)

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	packageInfo := &PackageInfo{}
	if stdout != nil {
		err = json.Unmarshal(stdout, packageInfo)
		if err != nil {
			return nil, err
		}
	}

	return packageInfo, nil
}

func installPackage(info *PackageInfo) error {
	pkg := fmt.Sprintf("%s/...@%s", info.Path, info.Version)
	markerPath := getMarkerPackagePath(info)

	cmd := exec.Command("go", "install", pkg)
	environmentVariables := []string{fmt.Sprintf("GOBIN=%s", markerPath)}
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, environmentVariables...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("could not install package %s", info.Name())
	}

	return nil
}

func getMarkerPackagePath(info *PackageInfo) string {
	goPath := goPath()
	return filepath.FromSlash(path.Join(goPath, "marker", "pkg", info.Path, info.Version))
}

func getMarkerPackagePathFromString(pkg, version string) string {
	goPath := goPath()
	return filepath.FromSlash(path.Join(goPath, "marker", "pkg", pkg, version))
}

func markerPackageExists(info *PackageInfo) bool {
	fileInfo, err := os.Stat(getMarkerPackagePath(info))
	if err == nil && fileInfo.IsDir() {
		return true
	}

	return false
}

func getModuleMarkerProcessorYamlPath(info *PackageInfo) string {
	goPath := goPath()
	return filepath.FromSlash(path.Join(goPath, "pkg", "mod", info.Name(), "marker.processors.yaml"))
}

func getMarkerPackageYamlPath(info *PackageInfo) string {
	goPath := goPath()
	return filepath.FromSlash(path.Join(goPath, "marker", "pkg", info.Path, info.Version, "marker.procesors.yaml"))
}

func isMarkerProcessorPackage(info *PackageInfo) bool {
	if _, err := os.Stat(getModuleMarkerProcessorYamlPath(info)); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}

func createMarkerPackageFolder(info *PackageInfo) error {
	markerPackagePath := getMarkerPackagePath(info)

	if _, err := os.Stat(markerPackagePath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(markerPackagePath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("marker path is not created for pkg %s", info.Name())
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
