package marker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsReservedMarker(t *testing.T) {
	assert.True(t, IsReservedMarker("import"))
	assert.True(t, IsReservedMarker("deprecated"))
	assert.True(t, IsReservedMarker("override"))
	assert.False(t, IsReservedMarker("nonReservedMarker"))
}
