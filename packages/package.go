package packages

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/tools/go/packages"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type MarkerPackage struct {
	Path               string   `json:"Path"`
	Version            string   `json:"Version"`
	LatestVersion      string   `json:"LatestVersion"`
	DownloadedVersions []string `json:"DownloadedVersions"`
	AvailableVersions  []string `json:"AvailableVersions"`
	Dir                string   `json:"Dir"`
	GoVersion          string   `json:"GoVersion"`
}

type PackageInfo struct {
	Path      string    `json:"Path"`
	Version   string    `json:"Version"`
	Versions  []string  `json:"Versions"`
	Time      time.Time `json:"Time"`
	Dir       string    `json:"Dir"`
	GoMod     string    `json:"GoMod"`
	GoVersion string    `json:"GoVersion"`
}

func (p *PackageInfo) Name() string {
	return fmt.Sprintf("%s@%s", p.Path, p.Version)
}

func (p *PackageInfo) ModulePath() string {
	goPath := GoPath()
	filePath := filepath.FromSlash(path.Join(goPath, "pkg", "mod", p.Name()))
	return filePath
}

func GetMarkerPackage(path string) (*MarkerPackage, error) {
	pkgParts := strings.SplitN(path, "@", 2)

	pkgName := path
	pkgVersion := "latest"

	if len(pkgParts) == 2 {
		pkgName = pkgParts[0]
		pkgVersion = pkgParts[1]
	}

	pkgInfo, err := GetPackageInfo(fmt.Sprintf("%s@%s", pkgName, pkgVersion))
	if err != nil {
		switch typedErr := err.(type) {
		case *exec.ExitError:
			return nil, errors.New(string(typedErr.Stderr))
		}
		return nil, err
	}

	markerPackage := &MarkerPackage{
		Path:               pkgInfo.Path,
		Version:            pkgInfo.Version,
		DownloadedVersions: make([]string, 0),
		AvailableVersions:  pkgInfo.Versions,
		Dir:                "",
		GoVersion:          pkgInfo.GoVersion,
	}

	for _, version := range pkgInfo.Versions {
		fileInfo, err := os.Stat(MarkerPackagePath(pkgInfo.Path, version))
		if err == nil && fileInfo.IsDir() {
			markerPackage.DownloadedVersions = append(markerPackage.DownloadedVersions, version)
		}
	}

	if len(pkgInfo.Versions) != 0 {
		markerPackage.LatestVersion = pkgInfo.Versions[len(pkgInfo.Versions)-1]
	}

	if pkgVersion == "latest" && len(markerPackage.DownloadedVersions) != 0 {
		markerPackage.Version = markerPackage.DownloadedVersions[len(markerPackage.DownloadedVersions)-1]
		markerPackage.Dir = MarkerPackagePath(pkgInfo.Path, markerPackage.Version)
	}

	return markerPackage, nil
}

func GetPackageInfo(path string) (*PackageInfo, error) {
	pkgParts := strings.SplitN(path, "@", 2)

	pkgName := path
	pkgVersion := "latest"

	if len(pkgParts) == 2 {
		pkgName = pkgParts[0]
		pkgVersion = pkgParts[1]
	}

	goArgs := []string{"list", "-m", "-versions", "-json", fmt.Sprintf("%s@%s", pkgName, pkgVersion)}
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

func GoPath() string {
	goArgs := []string{"env", "GOPATH"}
	cmd := exec.Command("go", goArgs...)

	stdout, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.Trim(string(stdout), "\n ")
}

func GoModDir() (string, error) {
	var wd string
	var err error
	wd, err = os.Getwd()

	if err != nil {
		return "", fmt.Errorf("wtf - what a terrible failure! : %s", err.Error())
	}

	config := &packages.Config{}
	config.Mode |= packages.NeedModule

	var pkgs []*packages.Package
	pkgs, err = packages.Load(config, wd)

	if err != nil {
		return "", fmt.Errorf("an error occurred : %s", err.Error())
	}

	if pkgs == nil || len(pkgs) == 0 {
		return "", fmt.Errorf("package not found for the directory %s", wd)
	}

	pkg := pkgs[0]

	if pkg.Module == nil {
		return "", fmt.Errorf("go.mod does not exist for the directory %s", wd)
	}

	return pkg.Module.Dir, nil
}

func InstallPackage(info *PackageInfo) error {
	pkg := fmt.Sprintf("%s/...@%s", info.Path, info.Version)
	markerPath := MarkerPackagePathFromPackageInfo(info)

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

func MarkerPackagePath(pkg, version string) string {
	goPath := GoPath()
	return filepath.FromSlash(path.Join(goPath, "marker", "pkg", pkg, version))
}

func MarkerPackagePathFromPackageInfo(info *PackageInfo) string {
	return MarkerPackagePath(info.Path, info.Version)
}

func MarkerProcessorYamlPath(info *PackageInfo) string {
	goPath := GoPath()
	return filepath.FromSlash(path.Join(goPath, "pkg", "mod", info.Name(), "marker.processors.yaml"))
}

func MarkerPackageYamlPath(info *PackageInfo) string {
	goPath := GoPath()
	return filepath.FromSlash(path.Join(goPath, "marker", "pkg", info.Path, info.Version, "marker.procesors.yaml"))
}
