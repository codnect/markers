package visitor

import (
	"fmt"
)

type Slice struct {
	elem Type
}

func (s *Slice) Name() string {
	return s.String()
}

func (s *Slice) Elem() Type {
	return s.elem
}

func (s *Slice) Underlying() Type {
	return s
}

func (s *Slice) String() string {
	return fmt.Sprintf("[]%s", s.elem.Name())
}
