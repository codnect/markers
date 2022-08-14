package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type MarkerYamlResource struct {
	Id         string                             `json:"id,omitempty"`
	Processors map[string]MarkerProcessorResource `json:"processors,omitempty"`
}

type MarkerProcessorResource struct {
	Markers []Marker `json:"markers,omitempty"`
}

type Marker struct {
	Name       string            `json:"name,omitempty"`
	Targets    []string          `json:"targets,omitempty"`
	Parameters []MarkerParameter `json:"parameters,omitempty"`
}

type MarkerParameter struct {
	Name        string            `json:"name,omitempty"`
	Type        string            `json:"type,omitempty"`
	Description string            `json:"description,omitempty"`
	Required    bool              `json:"required,omitempty"`
	Default     string            `json:"default,omitempty"`
	Enum        []MarkerEnumValue `json:"enum,omitempty"`
}

type MarkerEnumValue struct {
	Value       string `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
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

func (p *PackageInfo) Id() string {
	return fmt.Sprintf("%s@%s", p.Path, p.Version)
}

func (p *PackageInfo) ModulePath() string {
	goPath := getGoPath()
	return filepath.FromSlash(path.Join(goPath, "pkg", "mod", p.Id()))
}

func (p *PackageInfo) CmdPath() string {
	return filepath.FromSlash(path.Join(p.ModulePath(), "cmd"))
}

func (p *PackageInfo) MarkerYamlPath() string {
	return filepath.FromSlash(path.Join(p.CmdPath(), "marker.yaml"))
}

func (p *PackageInfo) MarkerProcessorPath() string {
	goPath := getGoPath()
	return filepath.FromSlash(path.Join(goPath, "marker", p.Path, p.Version))
}

func (p *PackageInfo) MarkerProcessorTempPath() string {
	return filepath.FromSlash(path.Join(p.MarkerProcessorPath(), "temp"))
}

func getGoPath() string {
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

func packageExists(info *PackageInfo) bool {
	if _, err := os.Stat(info.ModulePath()); err != nil {
		return true
	}

	return false
}

func installPackage(info *PackageInfo) error {
	if packageExists(info) {
		return fmt.Errorf("package '%s' already installed", info.Id())
	}

	goArgs := []string{"install", info.Id()}
	exec.Command("go", goArgs...).Run()

	if !packageExists(info) {
		return fmt.Errorf("could not install package '%s'", info.Id())
	}

	return nil
}

func createMarkerProcessorFolder(info *PackageInfo) {
	markerProcessorPath := info.MarkerProcessorPath()

	if _, err := os.Stat(markerProcessorPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(markerProcessorPath, os.ModePerm)
		if err != nil {
			log.Println(err)
		}

		CopyDir(info.ModulePath(), info.MarkerProcessorTempPath())
	}
}

func CopyFile(src, dst string) error {
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

func CopyDir(src, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}

	for _, fd := range fds {
		srcfb := path.Join(src, fd.Name())
		dstfb := path.Join(dst, fd.Name())
		if fd.IsDir() {
			if err = CopyDir(srcfb, dstfb); err != nil {
				fmt.Println(err)
				return err
			}
		} else {
			if err = CopyFile(srcfb, dstfb); err != nil {
				fmt.Println(err)
				return err
			}
		}
	}

	return nil
}

func buildMarkerProcessors(info *PackageInfo) error {
	goArgs := []string{"mod", "download"}
	cmd := exec.Command("go", goArgs...)
	cmd.Dir = info.MarkerProcessorTempPath()
	err := cmd.Run()
	if err != nil {
		return err
	}

	buildPath := filepath.FromSlash(path.Join(info.MarkerProcessorTempPath(), "cmd", "procyon"))
	cmd = exec.Command("go", "build", "-o", filepath.FromSlash(path.Join(info.MarkerProcessorPath(), "procyon")), buildPath)
	cmd.Dir = buildPath
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()

	return nil
}

func containMarkerProcessor(info *PackageInfo) bool {
	if _, err := os.Stat(info.CmdPath()); errors.Is(err, os.ErrNotExist) {
		return false
	}

	if _, err := os.Stat(info.MarkerYamlPath()); errors.Is(err, os.ErrNotExist) {
		return false
	}

	validateMarkerYaml(info.Path, info.MarkerYamlPath())

	return true
}

func validateMarkerYaml(pkgId, path string) bool {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	resource := &MarkerYamlResource{}
	err = yaml.Unmarshal(content, resource)
	if err != nil {
		return false
	}

	if resource.Id != pkgId {
		return false
	}

	if len(resource.Processors) == 0 {
		return false
	}

	return validateMarkerProcessors(resource.Processors)
}

func validateMarkerProcessors(processors map[string]MarkerProcessorResource) bool {
	for name, processor := range processors {
		validateMarkersV(name, processor.Markers)
	}

	return false
}

func validateMarkersV(processorName string, markers []Marker) bool {
	for _, marker := range markers {
		if marker.Name != processorName && !strings.HasPrefix(marker.Name, fmt.Sprintf("%s:", processorName)) {
			return false
		}

		validateMarker(marker)
	}

	return true
}

func validateMarker(marker Marker) {

}

func copyMarkerYamlFile(info *PackageInfo) {
	err := os.Rename(info.MarkerYamlPath(), filepath.FromSlash(path.Join(info.MarkerProcessorPath(), "marker.yaml")))
	if err != nil {

	}
}
