package processor

import (
	"encoding/json"
	"github.com/procyon-projects/marker/packages"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func getConfigFilePath() (string, error) {
	var err error
	var modDir string
	modDir, err = packages.GoModDir()

	if err != nil {
		return "", err
	}

	markerJsonFilePath := filepath.FromSlash(path.Join(modDir, "marker.json"))
	_, err = os.Stat(markerJsonFilePath)
	if err != nil {
		return "", err
	}

	return markerJsonFilePath, nil
}

func getConfig(configFilePath string) (*Config, error) {
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// getPackageDirectories finds the go module directory and returns
// the package directories.
func getPackageDirectories() ([]string, error) {
	var err error
	var modDir string
	modDir, err = packages.GoModDir()

	if err != nil {
		return nil, err
	}

	var dirs []string
	dirs, err = findDirectoriesWithGoFiles(modDir)

	if err != nil {
		return nil, err
	}

	return dirs, nil
}

// findDirectoriesWithGoFiles returns the go directories with go files.
// if not, an error might occur while loading packages.
func findDirectoriesWithGoFiles(root string) ([]string, error) {
	dirMap := make(map[string]bool, 0)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// skip hidden directories
		if strings.HasPrefix(path, ".") && !strings.HasPrefix(path, "./") {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if matched, err := filepath.Match("*.go", filepath.Base(path)); err != nil {
			return err
		} else if matched {
			dirMap[filepath.Dir(path)] = true
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	var dirs []string

	for dir, _ := range dirMap {
		dirs = append(dirs, dir)
	}

	return dirs, nil
}
