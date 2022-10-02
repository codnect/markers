package visitor

import (
	"fmt"
)

type Map struct {
	key  Type
	elem Type
}

func (m *Map) Name() string {
	return m.String()
}

func (m *Map) Key() Type {
	return m.key
}

func (m *Map) Elem() Type {
	return m.elem
}

func (m *Map) Underlying() Type {
	return m
}

func (m *Map) String() string {
	return fmt.Sprintf("map[%s]%s", m.key.Name(), m.elem.Name())
}
