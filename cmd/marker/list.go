package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/procyon-projects/marker/processor"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"strings"
)

var listCmd = &cobra.Command{
	Use:   "list [pkg]",
	Short: "list marker package",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("pkg is required")
		}

		return listPackage(args[0])
	},
}

func init() {
	processor.AddCommand(listCmd)
}

func listPackage(pkg string) error {
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

	markerPackage := &MarkerPackage{
		Path:               pkgInfo.Path,
		Version:            pkgInfo.Version,
		DownloadedVersions: make([]string, 0),
		AvailableVersions:  pkgInfo.Versions,
		Dir:                "",
		GoVersion:          pkgInfo.GoVersion,
	}

	for _, version := range pkgInfo.Versions {
		fileInfo, err := os.Stat(getMarkerPackagePathFromString(pkgInfo.Path, version))
		if err == nil && fileInfo.IsDir() {
			markerPackage.DownloadedVersions = append(markerPackage.DownloadedVersions, version)
		}
	}

	if len(pkgInfo.Versions) != 0 {
		markerPackage.LatestVersion = pkgInfo.Versions[len(pkgInfo.Versions)-1]
	}

	if pkgVersion == "latest" && len(markerPackage.DownloadedVersions) != 0 {
		markerPackage.Version = markerPackage.DownloadedVersions[len(markerPackage.DownloadedVersions)-1]
		markerPackage.Dir = getMarkerPackagePathFromString(pkgInfo.Path, markerPackage.Version)
	}

	jsonText, _ := json.MarshalIndent(markerPackage, "", "\t")
	log.Println(string(jsonText))
	return nil
}
