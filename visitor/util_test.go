package visitor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsStructType(t *testing.T) {
	assert.True(t, IsStructType(&Struct{}))
	assert.False(t, IsStructType(&Interface{}))
}

func TestIsInterface(t *testing.T) {
	assert.True(t, IsInterfaceType(&Interface{}))
	assert.False(t, IsInterfaceType(&Struct{}))
}
