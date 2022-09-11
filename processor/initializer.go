package processor

var (
	packageName      = ""
	processorName    = "marker"
	processorVersion = "1.0.0"
)

func Initialize(pkg, name, version string) {
	packageName = pkg
	processorName = name
	processorVersion = version
}
