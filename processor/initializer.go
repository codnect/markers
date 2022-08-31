package processor

var (
	processorName    = "marker"
	processorVersion = "1.0.0"
)

func Initialize(name, version string) {
	processorName = name
	processorVersion = version
}
