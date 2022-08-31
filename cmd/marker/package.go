package main

import (
	"fmt"
	"path"
	"path/filepath"
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
	goPath := goPath()
	filePath := filepath.FromSlash(path.Join(goPath, "pkg", "mod", p.Name()))
	return filePath
}
