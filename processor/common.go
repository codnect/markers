package processor

import (
	"github.com/procyon-projects/marker"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// getPackageDirectories finds the go module directory and returns
// the package directories.
func getPackageDirectories() ([]string, error) {
	var err error
	var modDir string
	modDir, err = marker.GoModDir()

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

// printErrors prints error.
func PrintError(err error) {
	if err != nil {

		switch typedErr := err.(type) {
		case marker.ErrorList:
			PrintErrors(typedErr)
			return
		}

		log.Println(err)
		return
	}
}

// printErrors prints the error list.
func PrintErrors(errorList marker.ErrorList) {
	if errorList == nil || len(errorList) == 0 {
		return
	}

	for _, err := range errorList {
		switch typedErr := err.(type) {
		case marker.Error:
			pos := typedErr.Position
			log.Printf("%s (%d:%d) : %s\n", typedErr.FileName, pos.Line, pos.Column, typedErr.Error())
		case marker.ParserError:
			pos := typedErr.Position
			log.Printf("%s (%d:%d) : %s\n", typedErr.FileName, pos.Line, pos.Column, typedErr.Error())
		case marker.ErrorList:
			PrintErrors(typedErr)
		default:
			PrintError(err)
		}
	}
}
