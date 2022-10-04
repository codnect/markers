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
