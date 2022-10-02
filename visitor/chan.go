package visitor

import "fmt"

type ChanDirection int

const (
	SendDir ChanDirection = 1 << iota
	ReceiveDir
	BothDir = SendDir | ReceiveDir
)

type Chan struct {
	direction ChanDirection
	elem      Type
}

func (c *Chan) Direction() ChanDirection {
	return c.direction
}

func (c *Chan) Elem() Type {
	return c.elem
}

func (c *Chan) Underlying() Type {
	return c
}

func (c *Chan) Name() string {
	return c.String()
}

func (c *Chan) String() string {
	if c.direction == BothDir {
		return fmt.Sprintf("chan %s", c.elem.Name())
	} else if c.direction == SendDir {
		return fmt.Sprintf("chan<- %s", c.elem.Name())
	}

	return fmt.Sprintf("<-chan %s", c.elem.Name())
}
