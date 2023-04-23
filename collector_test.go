package markers

import (
	"github.com/procyon-projects/markers/packages"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	result, _ := packages.LoadPackages("github.com/procyon-projects/markers/test/...")
	pkg, _ := result.Lookup("github.com/procyon-projects/markers/test/menu")

	registry := NewRegistry()
	collector := NewCollector(registry)

	nodes, err := collector.Collect(pkg)
	assert.NotNil(t, nodes)
	assert.NoError(t, err)

	pkg, _ = result.Lookup("github.com/procyon-projects/markers/test/any")

	nodes, err = collector.Collect(pkg)
	assert.NotNil(t, nodes)
	assert.NoError(t, err)
}
