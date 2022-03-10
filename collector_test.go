package marker

import (
	"log"
	"testing"
)

type X interface {
	Print()
}

type TestOutput struct {
}

type Z struct {
}

func (z Z) Print() {
	log.Println("")
}

func n(x X) {

}

func TestCollector_Collect(t *testing.T) {

	result, _ := LoadPackages("./test/package1")
	/*pkg, _ := result.Lookup("fmt")
	stringInterface := pkg.Types.Scope().Lookup("Stringer")

	pkgs := result.GetPackages()
	pkg1 := pkgs[0]
	str := types.NewPointer(pkg1.Types.Scope().Lookup("Fruit").Type())

	s := stringInterface.Type().Underlying()
	x := types.Implements(str, s.(*types.Interface))

	if x {

	}

	if str == nil {

	}

	if stringInterface == nil {

	}

	if pkg.IsStandardPackage() {

	}*/

	registry := NewRegistry()

	registry.Register("marker:package-level1", "github.com/procyon-projects/marker", PackageLevel, &TestOutput{})
	registry.Register("marker:package-level2", "github.com/procyon-project/marker", PackageLevel, &TestOutput{})

	collector := NewCollector(registry)

	eachFile(collector, result.GetPackages(), func(file *SourceFile, err error) {

	})

	/*EachFile(collector, result.GetPackages(), func(file *File, err error) {
		if file == nil {
		}
	})*/
}
