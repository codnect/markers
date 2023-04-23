package visitor

import "fmt"

type TypeConstraint struct {
	typ           Type
	tildeOperator bool
}

func (c *TypeConstraint) Name() string {
	return c.typ.String()
}

func (c *TypeConstraint) Type() Type {
	return c.typ
}

func (c *TypeConstraint) Underlying() Type {
	return c
}

func (c *TypeConstraint) Satisfy(t Type) bool {
	//TODO: implement this method
	return false
}

func (c *TypeConstraint) String() string {
	if c.tildeOperator {
		return fmt.Sprintf("~%s", c.typ.String())
	}
	return c.typ.String()
}

type TypeConstraints struct {
	elements []*TypeConstraint
}

func (c *TypeConstraints) Len() int {
	return len(c.elements)
}

func (c *TypeConstraints) At(index int) *TypeConstraint {
	if index >= 0 && index < len(c.elements) {
		return c.elements[index]
	}

	return nil
}
