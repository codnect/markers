package visitor

type ChanDirection int

const (
	SEND ChanDirection = 1 << iota
	RECEIVE
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
	return ""
}

func (c *Chan) String() string {
	return ""
}
