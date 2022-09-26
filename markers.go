package marker

import (
	"errors"
	"strings"
)

// Reserved markers
const (
	ImportMarkerName     = "import"
	DeprecatedMarkerName = "deprecated"
	OverrideMarkerName   = "override"
)

func IsReservedMarker(marker string) bool {
	if marker == ImportMarkerName || marker == DeprecatedMarkerName || marker == OverrideMarkerName {
		return true
	}

	return false
}

type ImportMarker struct {
	Value string `parameter:"Value" required:"true"`
	Alias string `parameter:"Alias" required:"false"`
	Pkg   string `parameter:"Pkg" required:"true"`
}

func (m ImportMarker) Validate() error {
	var errs []error

	if strings.Trim(m.Value, " \t") == "" {
		errs = append(errs, errors.New("'Value' argument cannot be nil or empty"))
	}

	if strings.Trim(m.Pkg, " \t") == "" {
		errs = append(errs, errors.New("'Pkg' argument cannot be nil or empty"))
	}

	if len(errs) != 0 {
		return NewErrorList(errs)
	}

	return nil
}

func (m ImportMarker) PkgPath() string {
	pkgParts := strings.Split(m.Pkg, "@")
	return pkgParts[0]
}

func (m ImportMarker) PkgVersion() string {
	pkgParts := strings.Split(m.Pkg, "@")

	if len(pkgParts) > 1 {
		return pkgParts[1]

	}

	return "latest"
}

type DeprecatedMarker struct {
	Value string `parameter:"Value"`
}

type OverrideMarker struct {
	Value string `parameter:"Value"`
}

type Marker struct {
	Value       string `parameter:"Value" required:"true"`
	Description string `parameter:"Description" required:"true"`
	Repeatable  bool   `parameter:"Repeatable" required:"false"`
	SyntaxFree  bool   `parameter:"SyntaxFree" required:"false"`
}

type Parameter struct {
	Value       string `parameter:"Value" required:"true"`
	Description string `parameter:"Description" required:"true"`
	Required    bool   `parameter:"Required" required:"false"`
	Deprecated  bool   `parameter:"Deprecated" required:"false"`
	Default     any    `parameter:"Default" required:"false"`
}

type ParameterEnum struct {
	Value string `parameter:"Value" required:"true"`
	Name  string `parameter:"Name" required:"true"`
}
