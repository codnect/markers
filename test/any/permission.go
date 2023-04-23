// +import=marker, Pkg=github.com/procyon-projects/markers, Alias=test
// +import=test-marker, Pkg=github.com/procyon-projects/test-markers
// +test-marker:package-level:Name=permission.go

package any

type Permission int

const (
	Read Permission = 1 << iota
	Write
	ReadWrite = Read | Write
)

type RequestMethod string

const (
	RequestGet    RequestMethod = "GET"
	RequestPost   RequestMethod = "POST"
	RequestPatch  RequestMethod = "PATCH"
	RequestDelete RequestMethod = "DELETE"
)

type Chan int

const (
	SendDir Chan = 2 >> iota
	ReceiveDir
	BothDir = SendDir | ReceiveDir
)
