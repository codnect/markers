package packages

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadResult_Packages(t *testing.T) {
	loadResult, err := LoadPackages("github.com/procyon-projects/marker/test/...")

	assert.Nil(t, err)
	assert.NotNil(t, loadResult)
	assert.Len(t, loadResult.Packages(), 1)

	pkg := loadResult.Packages()[0]
	assert.Equal(t, "package1", pkg.Name)
	assert.False(t, pkg.IsStandardPackage())
	assert.Equal(t, "github.com/procyon-projects/marker/test/package1", pkg.ID)
	assert.Equal(t, "github.com/procyon-projects/marker/test/package1", pkg.PkgPath)

	assert.NotNil(t, pkg.GoFiles)
	assert.NotNil(t, pkg.CompiledGoFiles)
	assert.NotNil(t, pkg.Syntax)

	assert.Len(t, pkg.GoFiles, 1)
	assert.Len(t, pkg.CompiledGoFiles, 1)
	assert.Len(t, pkg.Syntax, 1)

	assert.NotNil(t, pkg.Module)
}

func TestLoadResult_Lookup(t *testing.T) {
	loadResult, err := LoadPackages("github.com/procyon-projects/marker/test/...")

	assert.Nil(t, err)
	assert.NotNil(t, loadResult)
	assert.Len(t, loadResult.Packages(), 1)

	pkg, err := loadResult.Lookup("github.com/procyon-projects/marker/test/package1")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "package1", pkg.Name)
	assert.False(t, pkg.IsStandardPackage())
	assert.Equal(t, "github.com/procyon-projects/marker/test/package1", pkg.ID)
	assert.Equal(t, "github.com/procyon-projects/marker/test/package1", pkg.PkgPath)
}
