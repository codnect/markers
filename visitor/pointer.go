package visitor

import (
	"fmt"
)

type Pointer struct {
	base Type
}

func (p *Pointer) Name() string {
	return fmt.Sprintf("*%s", p.base.Name())
}

func (p *Pointer) Elem() Type {
	return p.base
}

func (p *Pointer) Underlying() Type {
	return p
}

func (p *Pointer) String() string {
	return fmt.Sprintf("*%s", p.base.Name())
}
