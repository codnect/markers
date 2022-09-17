package marker

import (
	"github.com/procyon-projects/marker/packages"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	result, _ := packages.LoadPackages("github.com/procyon-projects/marker/test/package1")
	pkg, _ := result.Lookup("github.com/procyon-projects/marker/test/package1")

	registry := NewRegistry()
	collector := NewCollector(registry)
	nodes, err := collector.Collect(pkg)
	assert.NotNil(t, nodes)
	assert.NoError(t, err)
}
