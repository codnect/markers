package packages

import (
	"testing"
)

func TestLoadPackages(t *testing.T) {
	/*pkgs, err := LoadPackages("./test/package1")

	assert.Nil(t, err)
	assert.NotNil(t, pkgs)
	assert.Len(t, pkgs, 1)

	assert.Equal(t, "package1", pkgs[0].Name)
	assert.Equal(t, "github.com/procyon-projects/marker/test/package1", pkgs[0].ID)
	assert.Equal(t, "github.com/procyon-projects/marker/test/package1", pkgs[0].PkgPath)

	assert.NotNil(t, pkgs[0].GoFiles)
	assert.NotNil(t, pkgs[0].CompiledGoFiles)
	assert.NotNil(t, pkgs[0].Syntax)

	assert.Len(t, pkgs[0].GoFiles, 1)
	assert.Len(t, pkgs[0].CompiledGoFiles, 1)
	assert.Len(t, pkgs[0].Syntax, 1)

	assert.NotNil(t, pkgs[0].Module)*/
}

func TestLoadMultiPackages(t *testing.T) {
	/*pkgs, err := LoadPackages("./test/...")

	assert.Nil(t, err)
	assert.NotNil(t, pkgs)
	assert.Len(t, pkgs, 2)

	assert.Equal(t, "package1", pkgs[0].Name)
	assert.Equal(t, "github.com/procyon-projects/marker/test/package1", pkgs[0].ID)
	assert.Equal(t, "github.com/procyon-projects/marker/test/package1", pkgs[0].PkgPath)

	assert.NotNil(t, pkgs[0].GoFiles)
	assert.NotNil(t, pkgs[0].CompiledGoFiles)
	assert.NotNil(t, pkgs[0].Syntax)

	assert.Len(t, pkgs[0].GoFiles, 1)
	assert.Len(t, pkgs[0].CompiledGoFiles, 1)
	assert.Len(t, pkgs[0].Syntax, 1)

	assert.NotNil(t, pkgs[0].Module)

	assert.Equal(t, "package2", pkgs[1].Name)
	assert.Equal(t, "github.com/procyon-projects/marker/test/package2", pkgs[1].ID)
	assert.Equal(t, "github.com/procyon-projects/marker/test/package2", pkgs[1].PkgPath)

	assert.NotNil(t, pkgs[1].GoFiles)
	assert.NotNil(t, pkgs[1].CompiledGoFiles)
	assert.NotNil(t, pkgs[1].Syntax)

	assert.Len(t, pkgs[1].GoFiles, 1)
	assert.Len(t, pkgs[1].CompiledGoFiles, 1)
	assert.Len(t, pkgs[1].Syntax, 1)

	assert.NotNil(t, pkgs[1].Module)*/
}
