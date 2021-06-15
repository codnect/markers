package marker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadPackages(t *testing.T) {
	pkgs, err := LoadPackages("./test-packages/foo")

	assert.Nil(t, err)
	assert.NotNil(t, pkgs)
	assert.Len(t, pkgs, 1)
}
