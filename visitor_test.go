package marker

import (
	"github.com/procyon-projects/marker/packages"
	"github.com/stretchr/testify/assert"
	"go/ast"
	"path/filepath"
	"testing"
)

func TestCommentVisitor_Visit(t *testing.T) {
	result, _ := packages.LoadPackages("./test/...")

	testCases := map[string]map[string]struct {
		commentGroup   int
		nodeMarkers    int
		packageMarkers int
	}{
		"github.com/procyon-projects/marker/test/menu": {
			"coffee.go": {
				commentGroup:   1,
				nodeMarkers:    1,
				packageMarkers: 2,
			},
			"dessert.go": {
				commentGroup:   30,
				nodeMarkers:    34,
				packageMarkers: 2,
			},
			"fresh.go": {
				commentGroup:   1,
				nodeMarkers:    1,
				packageMarkers: 2,
			},
		},
		"github.com/procyon-projects/marker/test/any": {
			"permission.go": {
				commentGroup:   1,
				nodeMarkers:    2,
				packageMarkers: 2,
			},
		},
	}

	for _, pkg := range result.Packages() {
		for _, file := range pkg.Syntax {
			fileName := filepath.Base(pkg.Fset.File(file.Pos()).Name())
			testCase, exists := testCases[pkg.ID][fileName]

			if !exists {
				t.Errorf("file %s not found in test cases", fileName)
				return
			}

			visitor := newCommentVisitor(file.Comments)
			ast.Walk(visitor, file)

			assert.Len(t, visitor.allComments, testCase.commentGroup)
			assert.Len(t, visitor.nodeMarkers, testCase.nodeMarkers)
			assert.Len(t, visitor.packageMarkers, testCase.packageMarkers)
		}
	}
}
