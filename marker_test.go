package marker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarkerValues_AllMarkers(t *testing.T) {
	markerValues := make(MarkerValues)
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest1")
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest2")
	markerValues["anyMarker2"] = append(markerValues["anyMarker2"], "anyTest3")

	assert.Equal(t, []interface{}{"anyTest1", "anyTest2"}, markerValues.AllMarkers("anyMarker1"))
	assert.Equal(t, []interface{}{"anyTest3"}, markerValues.AllMarkers("anyMarker2"))
}

func TestMarkerValues_Count(t *testing.T) {
	markerValues := make(MarkerValues)
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest1")
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest2")
	markerValues["anyMarker2"] = append(markerValues["anyMarker2"], "anyTest3")

	assert.Equal(t, 3, markerValues.Count())
}

func TestMarkerValues_CountByName(t *testing.T) {
	markerValues := make(MarkerValues)
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest1")
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest2")
	markerValues["anyMarker2"] = append(markerValues["anyMarker2"], "anyTest3")

	assert.Equal(t, 2, markerValues.CountByName("anyMarker1"))
	assert.Equal(t, 1, markerValues.CountByName("anyMarker2"))
}

func TestMarkerValues_First(t *testing.T) {
	markerValues := make(MarkerValues)
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest1")
	markerValues["anyMarker1"] = append(markerValues["anyMarker1"], "anyTest2")
	markerValues["anyMarker2"] = append(markerValues["anyMarker2"], "anyTest3")

	assert.Equal(t, "anyTest1", markerValues.First("anyMarker1"))
	assert.Equal(t, "anyTest3", markerValues.First("anyMarker2"))
}
