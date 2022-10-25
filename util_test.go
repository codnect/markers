package markers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsLower(t *testing.T) {
	assert.True(t, isLower("any"))
	assert.False(t, isLower("Any"))
}

func TestIsUpper(t *testing.T) {
	assert.True(t, isUpper("ANY"))
	assert.False(t, isUpper("Any"))
}

func TestLowerCamelCase(t *testing.T) {
	assert.Equal(t, "testAny", lowerCamelCase("TestAny"))
}

func TestUpperCamelCase(t *testing.T) {
	assert.Equal(t, "TestAny", upperCamelCase("testAny"))
}
