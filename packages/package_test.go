package packages

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/procyon-projects/markers/internal/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type mockExecutor struct {
	mock.Mock
}

func (m *mockExecutor) Execute(cmd cmd.Command) (output []byte, err error) {
	args := m.Called(cmd)

	if args.Get(0) != nil {
		output = args.Get(0).([]byte)
	}

	if args.Get(1) != nil {
		err = args.Get(1).(error)
	}

	return
}

func TestGetPackageInfo(t *testing.T) {
	mockExecutor := &mockExecutor{}
	cmd.SetCommandExecutor(mockExecutor)

	execLookupPath, _ := exec.LookPath("go")
	goPathCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "env", "GOPATH"},
	}

	goListCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "list", "-m", "-versions", "-json", fmt.Sprintf("%s@%s", "github.com/procyon-projects/marker", "latest")},
	}

	packageInfo := &PackageInfo{
		Path:      "github.com/procyon-projects/marker",
		Version:   "v2.0.5",
		Versions:  []string{"v1.2.6", "v2.0.5"},
		Time:      time.Time{},
		Dir:       "anyDir",
		GoMod:     "",
		GoVersion: "1.18",
	}
	byteData, _ := json.Marshal(packageInfo)

	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)
	mockExecutor.On("Execute", goListCmd).Return(byteData, nil)

	pkg, err := GetPackageInfo("github.com/procyon-projects/marker")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, packageInfo, pkg)
	assert.Equal(t, filepath.FromSlash("anyGoPath/pkg/mod/github.com/procyon-projects/marker@v2.0.5"), packageInfo.ModulePath())
}

func TestGetMarkerPackageShouldReturnErrorIfAnyErrorOccurs(t *testing.T) {
	mockExecutor := &mockExecutor{}
	cmd.SetCommandExecutor(mockExecutor)

	execLookupPath, _ := exec.LookPath("go")
	goPathCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "env", "GOPATH"},
	}

	goListCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "list", "-m", "-versions", "-json", fmt.Sprintf("%s@%s", "github.com/procyon-projects/marker", "latest")},
	}

	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)
	anyError := errors.New("anyError")
	mockExecutor.On("Execute", goListCmd).Return(nil, anyError)

	pkg, err := GetMarkerPackage("github.com/procyon-projects/marker")
	assert.Nil(t, pkg)
	assert.Equal(t, anyError, err)
}

func TestGetMarkerPackage(t *testing.T) {
	mockExecutor := &mockExecutor{}
	cmd.SetCommandExecutor(mockExecutor)

	execLookupPath, _ := exec.LookPath("go")
	goPathCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "env", "GOPATH"},
	}

	goListCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "list", "-m", "-versions", "-json", fmt.Sprintf("%s@%s", "github.com/procyon-projects/marker", "latest")},
	}

	packageInfo := &PackageInfo{
		Path:      "github.com/procyon-projects/marker",
		Version:   "v2.0.5",
		Versions:  []string{"v1.2.6", "v2.0.5"},
		Time:      time.Time{},
		Dir:       "anyDir",
		GoMod:     "",
		GoVersion: "1.18",
	}
	byteData, _ := json.Marshal(packageInfo)

	expectedMarkerPackage := &MarkerPackage{
		Path:               "github.com/procyon-projects/marker",
		Version:            "v2.0.5",
		LatestVersion:      "v2.0.5",
		DownloadedVersions: []string{},
		AvailableVersions:  []string{"v1.2.6", "v2.0.5"},
		Dir:                "",
		GoVersion:          "1.18",
	}

	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)
	mockExecutor.On("Execute", goListCmd).Return(byteData, nil)

	pkg, err := GetMarkerPackage("github.com/procyon-projects/marker")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, expectedMarkerPackage, pkg)
}

func TestGetMarkerPackageWithVersion(t *testing.T) {
	mockExecutor := &mockExecutor{}
	cmd.SetCommandExecutor(mockExecutor)

	execLookupPath, _ := exec.LookPath("go")
	goPathCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "env", "GOPATH"},
	}

	goListCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "list", "-m", "-versions", "-json", fmt.Sprintf("%s@%s", "github.com/procyon-projects/marker", "v2.0.5")},
	}

	packageInfo := &PackageInfo{
		Path:      "github.com/procyon-projects/marker",
		Version:   "v2.0.5",
		Versions:  []string{"v1.2.6", "v2.0.5"},
		Time:      time.Time{},
		Dir:       "anyDir",
		GoMod:     "",
		GoVersion: "1.18",
	}
	byteData, _ := json.Marshal(packageInfo)

	expectedMarkerPackage := &MarkerPackage{
		Path:               "github.com/procyon-projects/marker",
		Version:            "v2.0.5",
		LatestVersion:      "v2.0.5",
		DownloadedVersions: []string{},
		AvailableVersions:  []string{"v1.2.6", "v2.0.5"},
		Dir:                "",
		GoVersion:          "1.18",
	}

	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)
	mockExecutor.On("Execute", goListCmd).Return(byteData, nil)

	pkg, err := GetMarkerPackage("github.com/procyon-projects/marker@v2.0.5")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, expectedMarkerPackage, pkg)
}

func TestGoPath(t *testing.T) {
	mockExecutor := &mockExecutor{}
	cmd.SetCommandExecutor(mockExecutor)

	execLookupPath, _ := exec.LookPath("go")
	goPathCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "env", "GOPATH"},
	}

	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)
	assert.Equal(t, "anyGoPath", GoPath())
}

func TestMarkerPackagePath(t *testing.T) {
	mockExecutor := &mockExecutor{}
	cmd.SetCommandExecutor(mockExecutor)

	execLookupPath, _ := exec.LookPath("go")
	goPathCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "env", "GOPATH"},
	}

	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)

	assert.Equal(t, filepath.FromSlash("anyGoPath/marker/pkg/github.com/procyon-projects/marker/anyVersion"),
		MarkerPackagePath("github.com/procyon-projects/marker", "anyVersion"))
}

func TestMarkerPackagePathFromPackageInfo(t *testing.T) {
	mockExecutor := &mockExecutor{}
	cmd.SetCommandExecutor(mockExecutor)

	execLookupPath, _ := exec.LookPath("go")
	goPathCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "env", "GOPATH"},
	}

	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)

	assert.Equal(t, filepath.FromSlash("anyGoPath/marker/pkg/github.com/procyon-projects/marker/anyVersion"),
		MarkerPackagePathFromPackageInfo(&PackageInfo{
			Path:    "github.com/procyon-projects/marker",
			Version: "anyVersion",
		}))
}

func TestMarkerProcessorYamlPath(t *testing.T) {
	mockExecutor := &mockExecutor{}
	cmd.SetCommandExecutor(mockExecutor)

	execLookupPath, _ := exec.LookPath("go")
	goPathCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "env", "GOPATH"},
	}

	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)

	assert.Equal(t, filepath.FromSlash("anyGoPath/pkg/mod/github.com/procyon-projects/marker@anyVersion/marker.processors.yaml"),
		MarkerProcessorYamlPath(&PackageInfo{
			Path:    "github.com/procyon-projects/marker",
			Version: "anyVersion",
		}),
	)
}

func TestMarkerPackageYamlPath(t *testing.T) {
	mockExecutor := &mockExecutor{}
	cmd.SetCommandExecutor(mockExecutor)

	execLookupPath, _ := exec.LookPath("go")
	goPathCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "env", "GOPATH"},
	}

	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)

	assert.True(t, strings.HasSuffix(MarkerPackageYamlPath(&PackageInfo{
		Path:    "github.com/procyon-projects/marker",
		Version: "anyVersion",
	}), filepath.FromSlash("marker/pkg/github.com/procyon-projects/marker/anyVersion/marker.procesors.yaml")))
}

func TestGoModDir(t *testing.T) {
	_, err := GoModDir()
	assert.Nil(t, err)
}

func TestInstallPackageShouldInstallPackage(t *testing.T) {
	mockExecutor := &mockExecutor{}
	cmd.SetCommandExecutor(mockExecutor)

	execLookupPath, _ := exec.LookPath("go")
	goPathCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "env", "GOPATH"},
	}

	goInstallCmd := &exec.Cmd{
		Path:   execLookupPath,
		Args:   []string{"go", "install", "github.com/procyon-projects/chrono/...@latest"},
		Env:    []string{},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	environmentVariables := []string{fmt.Sprintf("GOBIN=%s", filepath.FromSlash("anyGoPath/marker/pkg/github.com/procyon-projects/chrono/latest"))}
	goInstallCmd.Env = append(goInstallCmd.Env, os.Environ()...)
	goInstallCmd.Env = append(goInstallCmd.Env, environmentVariables...)

	mockExecutor.On("Execute", goInstallCmd).Return(nil, nil)
	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)

	err := InstallPackage(&PackageInfo{
		Path:    "github.com/procyon-projects/chrono",
		Version: "latest",
	})
	assert.Nil(t, err)
}

func TestInstallPackageReturnsErrorIfInstallationIsFailed(t *testing.T) {
	mockExecutor := &mockExecutor{}
	cmd.SetCommandExecutor(mockExecutor)

	execLookupPath, _ := exec.LookPath("go")
	goPathCmd := &exec.Cmd{
		Path: execLookupPath,
		Args: []string{"go", "env", "GOPATH"},
	}

	goInstallCmd := &exec.Cmd{
		Path:   execLookupPath,
		Args:   []string{"go", "install", "github.com/procyon-projects/chrono/...@latest"},
		Env:    []string{},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	environmentVariables := []string{fmt.Sprintf("GOBIN=%s", filepath.FromSlash("anyGoPath/marker/pkg/github.com/procyon-projects/chrono/latest"))}
	goInstallCmd.Env = append(goInstallCmd.Env, os.Environ()...)
	goInstallCmd.Env = append(goInstallCmd.Env, environmentVariables...)

	mockExecutor.On("Execute", goInstallCmd).Return(nil, errors.New("anyInstallationError"))
	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)

	err := InstallPackage(&PackageInfo{
		Path:    "github.com/procyon-projects/chrono",
		Version: "latest",
	})
	assert.NotNil(t, err)
	assert.Equal(t, "could not install package github.com/procyon-projects/chrono@latest", err.Error())
}
