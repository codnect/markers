package markers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsLower(t *testing.T) {
	assert.True(t, IsLower("any"))
	assert.False(t, IsLower("Any"))
}

func TestIsUpper(t *testing.T) {
	assert.True(t, IsUpper("ANY"))
	assert.False(t, IsUpper("Any"))
}

func TestLowerCamelCase(t *testing.T) {
	assert.Equal(t, "testAny", LowerCamelCase("TestAny"))
}

func TestUpperCamelCase(t *testing.T) {
	assert.Equal(t, "TestAny", UpperCamelCase("testAny"))
}
