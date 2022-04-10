package visitor

import "strings"

type Map struct {
	key  Type
	elem Type
}

func (m *Map) Name() string {
	return ""
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
	var builder strings.Builder
	builder.WriteString("map[")
	builder.WriteString(m.key.String())
	builder.WriteString("]")
	builder.WriteString(m.elem.String())
	return builder.String()
}
