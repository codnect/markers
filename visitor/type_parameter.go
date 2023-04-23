package visitor

import "strings"

type TypeParameter struct {
	name        string
	constraints *TypeConstraints
}

func (t *TypeParameter) Name() string {
	return t.name
}

func (t *TypeParameter) TypeConstraints() *TypeConstraints {
	return t.constraints
}

func (t *TypeParameter) Underlying() Type {
	return t
}

func (t *TypeParameter) String() string {
	var builder strings.Builder
	builder.WriteString(t.name + " ")

	if t.TypeConstraints().Len() != 0 {
		for i := 0; i < t.TypeConstraints().Len(); i++ {
			constraint := t.TypeConstraints().At(i)
			builder.WriteString(constraint.String())

			if i != t.TypeConstraints().Len()-1 {
				builder.WriteString("|")
			}
		}
	}

	return builder.String()
}

type TypeParameters struct {
	elements []*TypeParameter
}

func (t *TypeParameters) Len() int {
	return len(t.elements)
}

func (t *TypeParameters) At(index int) *TypeParameter {
	if index >= 0 && index < len(t.elements) {
		return t.elements[index]
	}

	return nil
}

func (t *TypeParameters) FindByName(name string) (*TypeParameter, bool) {
	for _, typeParameter := range t.elements {
		if typeParameter.name == name {
			return typeParameter, true
		}
	}

	return nil, false
}
