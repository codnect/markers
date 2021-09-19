package marker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadPackages(t *testing.T) {
	pkgs, err := LoadPackages("./test/package1")

	assert.Nil(t, err)
	assert.NotNil(t, pkgs)
	assert.Len(t, pkgs, 1)

	assert.Equal(t, "package1", pkgs[0].Name)
	assert.Equal(t, "github.com/procyon-projects/marker/test/package1", pkgs[0].ID)
}
