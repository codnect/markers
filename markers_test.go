package markers

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

func TestImportMarker_PkgPath(t *testing.T) {
	importMarker := &Import{
		Pkg: "github.com/procyon-projects/marker@v1.2.3",
	}
	assert.Equal(t, "github.com/procyon-projects/marker", importMarker.PkgPath())
}

func TestImportMarker_PkgVersion(t *testing.T) {
	importMarker := &Import{
		Pkg: "github.com/procyon-projects/marker@v1.2.3",
	}
	assert.Equal(t, "v1.2.3", importMarker.PkgVersion())

}

func TestImportMarker_PkgVersionLatest(t *testing.T) {
	importMarker := &Import{
		Pkg: "github.com/procyon-projects/marker",
	}
	assert.Equal(t, "latest", importMarker.PkgVersion())

}

func TestImportMarker_Validate_IfValueIsMissing(t *testing.T) {
	importMarker := &Import{
		Pkg: "github.com/procyon-projects/marker",
	}
	assert.Error(t, importMarker.Validate())
}

func TestImportMarker_Validate_IfPkgIsMissing(t *testing.T) {
	importMarker := &Import{
		Value: "anyValue",
	}
	assert.Error(t, importMarker.Validate())
}

func TestImportMarker_ValidateShouldReturnNilIfValidationIsOkay(t *testing.T) {
	importMarker := &Import{
		Value: "anyValue",
		Pkg:   "anyPkg",
	}
	assert.Nil(t, importMarker.Validate())
}
