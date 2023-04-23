package visitor

import "strings"

type GenericType struct {
	rawType   Type
	arguments []Type
}

func (g *GenericType) Name() string {
	return g.rawType.Name()
}

func (g *GenericType) ActualTypeArguments() TypeSets {
	return g.arguments
}

func (g *GenericType) RawType() Type {
	return g.rawType
}

func (g *GenericType) Underlying() Type {
	return g
}

func (g *GenericType) String() string {
	var builder strings.Builder
	builder.WriteString(g.rawType.Name())

	if g.ActualTypeArguments().Len() != 0 {
		builder.WriteString("[")

		for index := 0; index < g.ActualTypeArguments().Len(); index++ {
			typeParam := g.ActualTypeArguments().At(index)
			builder.WriteString(typeParam.String())

			if index != g.ActualTypeArguments().Len()-1 {
				builder.WriteString(",")
			}
		}

		builder.WriteString("]")
	}

	return builder.String()
}
