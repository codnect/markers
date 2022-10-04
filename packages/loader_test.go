package packages

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadResult_Packages(t *testing.T) {
	loadResult, err := LoadPackages("github.com/procyon-projects/marker/test/...")

	assert.Nil(t, err)
	assert.NotNil(t, loadResult)
	assert.Len(t, loadResult.Packages(), 2)
}

func TestLoadResult_Lookup(t *testing.T) {
	loadResult, err := LoadPackages("github.com/procyon-projects/marker/test/...")

	assert.Nil(t, err)
	assert.NotNil(t, loadResult)
	assert.Len(t, loadResult.Packages(), 2)

	pkg, err := loadResult.Lookup("github.com/procyon-projects/marker/test/menu")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "menu", pkg.Name)
	assert.False(t, pkg.IsStandardPackage())
	assert.Equal(t, "github.com/procyon-projects/marker/test/menu", pkg.ID)
	assert.Equal(t, "github.com/procyon-projects/marker/test/menu", pkg.PkgPath)

	pkg, err = loadResult.Lookup("github.com/procyon-projects/marker/test/any")
	assert.Nil(t, err)
	assert.NotNil(t, pkg)
	assert.Equal(t, "any", pkg.Name)
	assert.False(t, pkg.IsStandardPackage())
	assert.Equal(t, "github.com/procyon-projects/marker/test/any", pkg.ID)
	assert.Equal(t, "github.com/procyon-projects/marker/test/any", pkg.PkgPath)
}
