/*
Copyright © 2021 Marker Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Contents
const (
	constantsFileContent = `/*
Copyright © 2021 [Company/Organization] Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

const (
	AppName    = "$APP_NAME"
	AppVersion = "1.0.0"
)
`
)

var module string

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a marker processor",
	Long:  `The init command lets you create a new marker processor.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("marker name is required")
		}

		name := args[0]

		wd, err := os.Getwd()

		if err != nil {
			return fmt.Errorf("wtf - what a terrible failure! : %s", err.Error())
		}

		err = downloadTemplateZip(wd)

		if err != nil {
			return err
		}

		err = extractTemplate(wd)
		deleteTemplateZip(wd)

		if err != nil {
			return err
		}

		writeConstantFile(wd, name)
		err = renameTemplateFiles(wd, name)

		if err != nil {
			return err
		}

		err = overwriteGoModFile(wd, name, module)

		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&module, "module", "m", "", "Module Name (required)")

	err := initCmd.MarkFlagRequired("module")

	if err != nil {
		panic(err)
	}
}

func downloadTemplateZip(wd string) error {
	templateUrl := fmt.Sprintf("%s/%s.zip", templateZipBaseUrl, templateVersionTag)

	var err error
	var resp *http.Response
	resp, err = http.Get(templateUrl)

	if err != nil {
		return fmt.Errorf("an error occurred while downloading template file : %s", err.Error())
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("downloading template file could not be downloaded successfully")
	}

	// Create the file
	var out *os.File
	out, err = os.Create(filepath.Join(wd, templateZipFileName))
	if err != nil {
		return fmt.Errorf("template.zip could not be created successfully")
	}

	defer out.Close()
	_, err = io.Copy(out, resp.Body)

	if err != nil {
		return fmt.Errorf("response body could not be written to template.zip file successfully")
	}

	return nil
}

func extractTemplate(wd string) error {
	archive, err := zip.OpenReader(templateZipFileName)

	if err != nil {
		return fmt.Errorf("template.zip could not be opened : %s", err.Error())
	}

	defer archive.Close()

	for _, file := range archive.File {
		filePath := filepath.Join(wd, file.Name)

		if !strings.HasPrefix(filePath, filepath.Clean(wd)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid file path %s", filePath)
		}
		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := file.Open()

		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}

func renameTemplateFiles(wd string, processorName string) error {
	// rename cmd folder
	err := os.Rename(filepath.Join(wd, templateRootFolderName, templateCmdFolderName, templateProcessorFolderName),
		filepath.Join(wd, templateRootFolderName, templateCmdFolderName, processorName))

	if err != nil {
		os.RemoveAll(filepath.Join(wd, templateRootFolderName))
		return err
	}

	// rename root folder
	err = os.Rename(filepath.Join(wd, templateRootFolderName), filepath.Join(wd, processorName))

	if err != nil {
		os.RemoveAll(filepath.Join(wd, templateRootFolderName))
		return err
	}

	return nil
}

func deleteTemplateZip(wd string) {
	os.Remove(filepath.Join(wd, templateZipFileName))
}

func writeConstantFile(wd, processorName string) {
	filePath := filepath.Join(wd, templateRootFolderName, templateCmdFolderName, templateProcessorFolderName, constantsFileName)
	os.Remove(filePath)

	content := strings.ReplaceAll(constantsFileContent, "$APP_NAME", processorName)
	writeFile(filePath, content)
}

func overwriteGoModFile(wd, processorName, module string) error {
	err := os.Remove(filepath.Join(wd, processorName, goModFileName))

	if err != nil {
		return err
	}

	err = exec.Command("go", "mod", "init", module).Run()

	if err != nil {
		return err
	}

	err = exec.Command("go", "get", "-t", "-v", "./...").Run()

	return err
}

func writeFile(filePath, content string) error {
	return ioutil.WriteFile(filePath, []byte(content), os.ModePerm)
}
