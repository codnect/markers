package package2

type BaseRequest struct {
	Id int
}

type Request struct {
	BaseRequest
	Name string
	Test string
}
type ComplexRequest struct {
	Request
	Unique int64
}

type WriteI interface {
	Write()
}

type ReadI interface {
	Read()
}

type IFace interface {
	WriteI
	ReadI
}
