package visitor

type testFile struct {
	interfaces map[string]interfaceInfo
	structs    map[string]structInfo
	functions  map[string]functionInfo
}
