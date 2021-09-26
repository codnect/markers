package marker

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"os"
	"strings"
	"unicode"
)

func IsDecimal(character rune) bool {
	return '0' <= character && character <= '9'
}

func IsIdentifier(character rune, index int) bool {
	return character == '_' || unicode.IsLetter(character) || unicode.IsDigit(character) && index > 0
}

func LowerCamelCase(str string) string {
	isFirst := true

	return strings.Map(func(r rune) rune {
		if isFirst {
			isFirst = false
			return unicode.ToLower(r)
		}

		return r
	}, str)

}

// GoModDir returns the directory of go.mod file.
func GoModDir() (string, error) {
	var wd string
	var err error
	wd, err = os.Getwd()

	if err != nil {
		return "", fmt.Errorf("wtf - what a terrible failure! : %s", err.Error())
	}

	config := &packages.Config{}
	config.Mode |= packages.NeedModule

	var pkgs []*packages.Package
	pkgs, err = packages.Load(config, wd)

	if err != nil {
		return "", fmt.Errorf("an error occurred : %s", err.Error())
	}

	if pkgs == nil || len(pkgs) == 0 {
		return "", fmt.Errorf("package not found for the directory %s", wd)
	}

	pkg := pkgs[0]

	if pkg.Module == nil {
		return "", fmt.Errorf("go.mod does not exist for the directory %s", wd)
	}

	return pkg.Module.Dir, nil
}
