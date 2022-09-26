package packages

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

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

func TestGoModDir(t *testing.T) {
	modDir, err := GoModDir()
	assert.Nil(t, err)
	assert.NotEmpty(t, modDir)
	assert.True(t, strings.HasSuffix(modDir, "/marker"))
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
