package packages

import (
	"errors"
	"fmt"
	"github.com/procyon-projects/marker/internal/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"os/exec"
	"strings"
	"testing"
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
	pkg, err := GetPackageInfo("github.com/procyon-projects/marker")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "github.com/procyon-projects/marker", pkg.Path)
	assert.NotEmpty(t, pkg.Name())
	assert.NotEmpty(t, pkg.ModulePath())
	assert.True(t, strings.HasSuffix(pkg.ModulePath(), "/pkg/mod/"+pkg.Name()))
}

func TestGetMarkerPackage(t *testing.T) {
	pkg, err := GetMarkerPackage("github.com/procyon-projects/marker")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "github.com/procyon-projects/marker", pkg.Path)
}

func TestGoPath(t *testing.T) {
	assert.NotEmpty(t, GoPath())
}

func TestMarkerPackagePath(t *testing.T) {
	assert.True(t, strings.HasSuffix(MarkerPackagePath("github.com/procyon-projects/marker", "anyVersion"),
		"/marker/pkg/github.com/procyon-projects/marker/anyVersion"))
}

func TestMarkerPackagePathFromPackageInfo(t *testing.T) {
	assert.True(t, strings.HasSuffix(MarkerPackagePathFromPackageInfo(&PackageInfo{
		Path:    "github.com/procyon-projects/marker",
		Version: "anyVersion",
	}), "/marker/pkg/github.com/procyon-projects/marker/anyVersion"))
}

func TestMarkerProcessorYamlPath(t *testing.T) {
	assert.True(t, strings.HasSuffix(MarkerProcessorYamlPath(&PackageInfo{
		Path:    "github.com/procyon-projects/marker",
		Version: "anyVersion",
	}), "/pkg/mod/github.com/procyon-projects/marker@anyVersion/marker.processors.yaml"))
}

func TestMarkerPackageYamlPath(t *testing.T) {
	assert.True(t, strings.HasSuffix(MarkerPackageYamlPath(&PackageInfo{
		Path:    "github.com/procyon-projects/marker",
		Version: "anyVersion",
	}), "marker/pkg/github.com/procyon-projects/marker/anyVersion/marker.procesors.yaml"))
}

func TestGoModDir(t *testing.T) {
	//_, err := GoModDir()
	//assert.Nil(t, err)
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
		Path:   "/usr/local/go/bin/go",
		Args:   []string{"go", "install", "github.com/procyon-projects/chrono/...@latest"},
		Env:    []string{},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	environmentVariables := []string{fmt.Sprintf("GOBIN=%s", "anyGoPath/marker/pkg/github.com/procyon-projects/chrono/latest")}
	goInstallCmd.Env = append(goInstallCmd.Env, os.Environ()...)
	goInstallCmd.Env = append(goInstallCmd.Env, environmentVariables...)

	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)
	mockExecutor.On("Execute", goInstallCmd).Return(nil, nil)

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
		Path:   "/usr/local/go/bin/go",
		Args:   []string{"go", "install", "github.com/procyon-projects/chrono/...@latest"},
		Env:    []string{},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	environmentVariables := []string{fmt.Sprintf("GOBIN=%s", "anyGoPath/marker/pkg/github.com/procyon-projects/chrono/latest")}
	goInstallCmd.Env = append(goInstallCmd.Env, os.Environ()...)
	goInstallCmd.Env = append(goInstallCmd.Env, environmentVariables...)

	mockExecutor.On("Execute", goPathCmd).Return([]byte("anyGoPath\n "), nil)
	mockExecutor.On("Execute", goInstallCmd).Return(nil, errors.New("anyInstallationError"))

	err := InstallPackage(&PackageInfo{
		Path:    "github.com/procyon-projects/chrono",
		Version: "latest",
	})
	assert.NotNil(t, err)
	assert.Equal(t, "could not install package github.com/procyon-projects/chrono@latest", err.Error())
}
