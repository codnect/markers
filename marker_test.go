package markers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarkerValues_AllMarkers(t *testing.T) {
	markerValues := make(Values)
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest1")
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest2")
	markerValues["anyMarker2"] = append(markerValues["anyMarker2"], "anyTest3")

	markers, exists := markerValues.FindByName("anyMarker1")
	assert.True(t, exists)
	assert.Equal(t, []interface{}{"anyTest1", "anyTest2"}, markers)

	markers, exists = markerValues.FindByName("anyMarker2")
	assert.True(t, exists)
	assert.Equal(t, []interface{}{"anyTest3"}, markers)
}

func TestMarkerValues_Count(t *testing.T) {
	markerValues := make(Values)
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest1")
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest2")
	markerValues["anyMarker2"] = append(markerValues["anyMarker2"], "anyTest3")

	assert.Equal(t, 3, markerValues.Count())
}

func TestMarkerValues_CountByName(t *testing.T) {
	markerValues := make(Values)
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest1")
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest2")
	markerValues["anyMarker2"] = append(markerValues["anyMarker2"], "anyTest3")

	assert.Equal(t, 2, markerValues.CountByName("anyMarker1"))
	assert.Equal(t, 1, markerValues.CountByName("anyMarker2"))
}

func TestMarkerValues_First(t *testing.T) {
	markerValues := make(Values)
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest1")
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest2")
	markerValues["anyMarker2"] = append(markerValues["anyMarker2"], "anyTest3")

	marker, exists := markerValues.First("anyMarker1")
	assert.True(t, exists)
	assert.Equal(t, "anyTest1", marker)

	marker, exists = markerValues.First("anyMarker2")
	assert.True(t, exists)
	assert.Equal(t, "anyTest3", marker)
}
